// Copyright 2021-2023 vArmor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package profile

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilrand "k8s.io/apimachinery/pkg/util/rand"

	varmor "github.com/bytedance/vArmor/apis/varmor/v1beta1"
	varmorconfig "github.com/bytedance/vArmor/internal/config"
	apparmorprofile "github.com/bytedance/vArmor/internal/profile/apparmor"
	bpfprofile "github.com/bytedance/vArmor/internal/profile/bpf"
	varmortypes "github.com/bytedance/vArmor/internal/types"
	varmorinterface "github.com/bytedance/vArmor/pkg/client/clientset/versioned/typed/varmor/v1beta1"
)

// profileNameTemplate is the name of ArmorProfile object in k8s and AppArmor profile in host machine.
//
//	For namespace-scope profile, its format is "varmor-{VarmorProfile Namespace}-{VarmorProfile Name}"
//	For cluster-scope profile, its format is "varmor-cluster-{vArmor Namespace}-{VarmorClusterProfile Name}"
const (
	ClusterProfileNameTemplate = "varmor-cluster-%s-%s"
	ProfileNameTemplate        = "varmor-%s-%s"
)

func GenerateArmorProfileName(ns string, name string, clusterScope bool) string {
	profileName := ""

	if clusterScope {
		profileName = fmt.Sprintf(ClusterProfileNameTemplate, varmorconfig.Namespace, name)
	} else {
		profileName = fmt.Sprintf(ProfileNameTemplate, ns, name)
	}

	return strings.ToLower(profileName)
}

func GenerateProfile(policy varmor.Policy, name string, namespace string, varmorInterface varmorinterface.CrdV1beta1Interface, complete bool) (*varmor.Profile, error) {
	var err error

	profile := varmor.Profile{
		Name:     name,
		Enforcer: policy.Enforcer,
		Mode:     "enforce",
	}

	switch policy.Mode {
	case varmortypes.AlwaysAllowMode:
		switch policy.Enforcer {
		case "AppArmor":
			profile.Content = apparmorprofile.GenerateAlwaysAllowProfile(name)
		case "BPF":
			var bpfContent varmor.BpfContent
			profile.BpfContent = bpfContent
		default:
			return nil, fmt.Errorf("unknown enforcer")
		}

	case varmortypes.RuntimeDefaultMode:
		switch policy.Enforcer {
		case "AppArmor":
			profile.Content = apparmorprofile.GenerateRuntimeDefaultProfile(name)
		case "BPF":
			var bpfContent varmor.BpfContent
			err = bpfprofile.GenerateRuntimeDefaultProfile(&bpfContent)
			if err != nil {
				return nil, err
			}
			profile.BpfContent = bpfContent
		default:
			return nil, fmt.Errorf("unknown enforcer")
		}

	case varmortypes.EnhanceProtectMode:
		switch policy.Enforcer {
		case "AppArmor":
			profile.Content = apparmorprofile.GenerateEnhanceProtectProfile(&policy.EnhanceProtect, name, policy.Privileged)
		case "BPF":
			var bpfContent varmor.BpfContent
			err = bpfprofile.GenerateEnhanceProtectProfile(&policy.EnhanceProtect, &bpfContent, policy.Privileged)
			if err != nil {
				return nil, err
			}
			profile.BpfContent = bpfContent
		default:
			return nil, fmt.Errorf("unknown enforcer")
		}

	// [Experimental feature] Compatible with KubeArmor's SecurityPolicy
	case varmortypes.CustomPolicyMode:
		switch policy.Enforcer {
		case "AppArmor":
			profile.Content = apparmorprofile.GenerateCustomPolicyProfile(policy, name)
		case "BPF":
			return nil, fmt.Errorf("not supported by the BPF enforcer")
		default:
			return nil, fmt.Errorf("unknown enforcer")
		}

	case varmortypes.DefenseInDepthMode:
		switch policy.Enforcer {
		case "AppArmor":
			if policy.ModelOptions.UseExistingModel {
				profile.Mode = "enforce"
				apm, err := varmorInterface.ArmorProfileModels(namespace).Get(context.Background(), name, metav1.GetOptions{})
				if err == nil {
					profile.Content = apm.Spec.Profile.Content
				} else {
					return nil, fmt.Errorf("no models found")
				}
			} else {
				if complete {
					// Create profile based on the AlwaysAllow template after the behvior modeling was completed.
					profile.Content = apparmorprofile.GenerateAlwaysAllowProfile(name)
				} else {
					profile.Mode = "complain"
					profile.Content = apparmorprofile.GenerateBehaviorModelingProfile(name)
				}
			}
		case "BPF":
			return nil, fmt.Errorf("not supported by the BPF enforcer")
		default:
			return nil, fmt.Errorf("unknown enforcer")
		}

	default:
		return nil, fmt.Errorf("unknown mode")
	}

	return &profile, nil
}

func NewArmorProfile(obj interface{}, varmorInterface varmorinterface.CrdV1beta1Interface, clusterScope bool) (*varmor.ArmorProfile, error) {
	ap := varmor.ArmorProfile{}

	if clusterScope {
		vcp := obj.(*varmor.VarmorClusterPolicy)
		profileName := GenerateArmorProfileName("", vcp.Name, clusterScope)

		ap.Name = profileName
		ap.Namespace = varmorconfig.Namespace
		ap.Labels = vcp.ObjectMeta.DeepCopy().Labels

		profile, err := GenerateProfile(vcp.Spec.Policy, ap.Name, ap.Namespace, varmorInterface, false)
		if err != nil {
			return nil, err
		}
		ap.Spec.Profile = *profile
		ap.Spec.Target = *vcp.Spec.Target.DeepCopy()

	} else {
		vp := obj.(*varmor.VarmorPolicy)
		profileName := GenerateArmorProfileName(vp.Namespace, vp.Name, clusterScope)

		ap.Name = profileName
		ap.Namespace = vp.Namespace
		ap.Labels = vp.ObjectMeta.DeepCopy().Labels

		profile, err := GenerateProfile(vp.Spec.Policy, ap.Name, ap.Namespace, varmorInterface, false)
		if err != nil {
			return nil, err
		}
		ap.Spec.Profile = *profile
		ap.Spec.Target = *vp.Spec.Target.DeepCopy()

		if vp.Spec.Policy.Mode == varmortypes.DefenseInDepthMode {
			ap.Spec.BehaviorModeling.Enable = true
			ap.Spec.BehaviorModeling.ModelingDuration = vp.Spec.Policy.ModelOptions.ModelingDuration
			ap.Spec.BehaviorModeling.UniqueID = utilrand.String(8)
		}
	}

	return &ap, nil
}
