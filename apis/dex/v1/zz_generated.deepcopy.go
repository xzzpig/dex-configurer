//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2021 xzzpig.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DexAuthClient) DeepCopyInto(out *DexAuthClient) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DexAuthClient.
func (in *DexAuthClient) DeepCopy() *DexAuthClient {
	if in == nil {
		return nil
	}
	out := new(DexAuthClient)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DexAuthClient) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DexAuthClientList) DeepCopyInto(out *DexAuthClientList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DexAuthClient, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DexAuthClientList.
func (in *DexAuthClientList) DeepCopy() *DexAuthClientList {
	if in == nil {
		return nil
	}
	out := new(DexAuthClientList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DexAuthClientList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DexAuthClientSpec) DeepCopyInto(out *DexAuthClientSpec) {
	*out = *in
	if in.RedirectURIs != nil {
		in, out := &in.RedirectURIs, &out.RedirectURIs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.SecretRef = in.SecretRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DexAuthClientSpec.
func (in *DexAuthClientSpec) DeepCopy() *DexAuthClientSpec {
	if in == nil {
		return nil
	}
	out := new(DexAuthClientSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DexAuthClientStatus) DeepCopyInto(out *DexAuthClientStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DexAuthClientStatus.
func (in *DexAuthClientStatus) DeepCopy() *DexAuthClientStatus {
	if in == nil {
		return nil
	}
	out := new(DexAuthClientStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DexClientOrder) DeepCopyInto(out *DexClientOrder) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DexClientOrder.
func (in *DexClientOrder) DeepCopy() *DexClientOrder {
	if in == nil {
		return nil
	}
	out := new(DexClientOrder)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DexClientOrder) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DexClientOrderList) DeepCopyInto(out *DexClientOrderList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DexClientOrder, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DexClientOrderList.
func (in *DexClientOrderList) DeepCopy() *DexClientOrderList {
	if in == nil {
		return nil
	}
	out := new(DexClientOrderList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DexClientOrderList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DexClientOrderRefObjects) DeepCopyInto(out *DexClientOrderRefObjects) {
	*out = *in
	out.ClientRef = in.ClientRef
	out.DeploymentRef = in.DeploymentRef
	out.ServiceRef = in.ServiceRef
	out.IngressRef = in.IngressRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DexClientOrderRefObjects.
func (in *DexClientOrderRefObjects) DeepCopy() *DexClientOrderRefObjects {
	if in == nil {
		return nil
	}
	out := new(DexClientOrderRefObjects)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DexClientOrderSpec) DeepCopyInto(out *DexClientOrderSpec) {
	*out = *in
	out.Config = in.Config
	out.TargetIngress = in.TargetIngress
	out.RedirectUrl = in.RedirectUrl
	if in.AllowedGroups != nil {
		in, out := &in.AllowedGroups, &out.AllowedGroups
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ExtraArguments != nil {
		in, out := &in.ExtraArguments, &out.ExtraArguments
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ExtraRedirectUrls != nil {
		in, out := &in.ExtraRedirectUrls, &out.ExtraRedirectUrls
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DexClientOrderSpec.
func (in *DexClientOrderSpec) DeepCopy() *DexClientOrderSpec {
	if in == nil {
		return nil
	}
	out := new(DexClientOrderSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DexClientOrderStatus) DeepCopyInto(out *DexClientOrderStatus) {
	*out = *in
	out.RefObjects = in.RefObjects
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DexClientOrderStatus.
func (in *DexClientOrderStatus) DeepCopy() *DexClientOrderStatus {
	if in == nil {
		return nil
	}
	out := new(DexClientOrderStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DexProxyConfig) DeepCopyInto(out *DexProxyConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DexProxyConfig.
func (in *DexProxyConfig) DeepCopy() *DexProxyConfig {
	if in == nil {
		return nil
	}
	out := new(DexProxyConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DexProxyConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DexProxyConfigList) DeepCopyInto(out *DexProxyConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DexProxyConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DexProxyConfigList.
func (in *DexProxyConfigList) DeepCopy() *DexProxyConfigList {
	if in == nil {
		return nil
	}
	out := new(DexProxyConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DexProxyConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DexProxyConfigSpec) DeepCopyInto(out *DexProxyConfigSpec) {
	*out = *in
	out.SecretRef = in.SecretRef
	out.DefaultUrl = in.DefaultUrl
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DexProxyConfigSpec.
func (in *DexProxyConfigSpec) DeepCopy() *DexProxyConfigSpec {
	if in == nil {
		return nil
	}
	out := new(DexProxyConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DexProxyConfigStatus) DeepCopyInto(out *DexProxyConfigStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DexProxyConfigStatus.
func (in *DexProxyConfigStatus) DeepCopy() *DexProxyConfigStatus {
	if in == nil {
		return nil
	}
	out := new(DexProxyConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DexRedirectUrl) DeepCopyInto(out *DexRedirectUrl) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DexRedirectUrl.
func (in *DexRedirectUrl) DeepCopy() *DexRedirectUrl {
	if in == nil {
		return nil
	}
	out := new(DexRedirectUrl)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DexSecretRef) DeepCopyInto(out *DexSecretRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DexSecretRef.
func (in *DexSecretRef) DeepCopy() *DexSecretRef {
	if in == nil {
		return nil
	}
	out := new(DexSecretRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespacedName) DeepCopyInto(out *NamespacedName) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespacedName.
func (in *NamespacedName) DeepCopy() *NamespacedName {
	if in == nil {
		return nil
	}
	out := new(NamespacedName)
	in.DeepCopyInto(out)
	return out
}
