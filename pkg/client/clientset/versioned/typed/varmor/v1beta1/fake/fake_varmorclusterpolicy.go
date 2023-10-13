/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1beta1 "github.com/bytedance/vArmor/apis/varmor/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeVarmorClusterPolicies implements VarmorClusterPolicyInterface
type FakeVarmorClusterPolicies struct {
	Fake *FakeCrdV1beta1
}

var varmorclusterpoliciesResource = schema.GroupVersionResource{Group: "crd.varmor.org", Version: "v1beta1", Resource: "varmorclusterpolicies"}

var varmorclusterpoliciesKind = schema.GroupVersionKind{Group: "crd.varmor.org", Version: "v1beta1", Kind: "VarmorClusterPolicy"}

// Get takes name of the varmorClusterPolicy, and returns the corresponding varmorClusterPolicy object, and an error if there is any.
func (c *FakeVarmorClusterPolicies) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.VarmorClusterPolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(varmorclusterpoliciesResource, name), &v1beta1.VarmorClusterPolicy{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VarmorClusterPolicy), err
}

// List takes label and field selectors, and returns the list of VarmorClusterPolicies that match those selectors.
func (c *FakeVarmorClusterPolicies) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.VarmorClusterPolicyList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(varmorclusterpoliciesResource, varmorclusterpoliciesKind, opts), &v1beta1.VarmorClusterPolicyList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.VarmorClusterPolicyList{ListMeta: obj.(*v1beta1.VarmorClusterPolicyList).ListMeta}
	for _, item := range obj.(*v1beta1.VarmorClusterPolicyList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested varmorClusterPolicies.
func (c *FakeVarmorClusterPolicies) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(varmorclusterpoliciesResource, opts))
}

// Create takes the representation of a varmorClusterPolicy and creates it.  Returns the server's representation of the varmorClusterPolicy, and an error, if there is any.
func (c *FakeVarmorClusterPolicies) Create(ctx context.Context, varmorClusterPolicy *v1beta1.VarmorClusterPolicy, opts v1.CreateOptions) (result *v1beta1.VarmorClusterPolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(varmorclusterpoliciesResource, varmorClusterPolicy), &v1beta1.VarmorClusterPolicy{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VarmorClusterPolicy), err
}

// Update takes the representation of a varmorClusterPolicy and updates it. Returns the server's representation of the varmorClusterPolicy, and an error, if there is any.
func (c *FakeVarmorClusterPolicies) Update(ctx context.Context, varmorClusterPolicy *v1beta1.VarmorClusterPolicy, opts v1.UpdateOptions) (result *v1beta1.VarmorClusterPolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(varmorclusterpoliciesResource, varmorClusterPolicy), &v1beta1.VarmorClusterPolicy{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VarmorClusterPolicy), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeVarmorClusterPolicies) UpdateStatus(ctx context.Context, varmorClusterPolicy *v1beta1.VarmorClusterPolicy, opts v1.UpdateOptions) (*v1beta1.VarmorClusterPolicy, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(varmorclusterpoliciesResource, "status", varmorClusterPolicy), &v1beta1.VarmorClusterPolicy{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VarmorClusterPolicy), err
}

// Delete takes name of the varmorClusterPolicy and deletes it. Returns an error if one occurs.
func (c *FakeVarmorClusterPolicies) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(varmorclusterpoliciesResource, name, opts), &v1beta1.VarmorClusterPolicy{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeVarmorClusterPolicies) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(varmorclusterpoliciesResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.VarmorClusterPolicyList{})
	return err
}

// Patch applies the patch and returns the patched varmorClusterPolicy.
func (c *FakeVarmorClusterPolicies) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.VarmorClusterPolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(varmorclusterpoliciesResource, name, pt, data, subresources...), &v1beta1.VarmorClusterPolicy{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.VarmorClusterPolicy), err
}
