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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DexProxyConfigSpec defines the desired state of DexProxyConfig
type DexProxyConfigSpec struct {
	OcidIssuerUrl string `json:"ocid-issuer-url"`

	// +kubebuilder:default:Dex Login
	// +kubebuilder:validation:Optional
	ProviderDisplayName string `json:"provider-display-name"`

	// The K8s Secret Ref would be modified
	SecretRef DexSecretRef `json:"secretRef"`

	// +kubebuilder:validation:Optional
	CookieSecret string `json:"cookie-secret"`

	// +kubebuilder:validation:Optional
	ProxyImage string `json:"proxy-image"`

	// Default RedirectUrl will be Used when RedirectUrl's Filed is Empty in the DexClientOrder
	// +kubebuilder:validation:Optional
	DefaultUrl DexRedirectUrl `json:"default-url"`
}

// DexProxyConfigStatus defines the observed state of DexProxyConfig
type DexProxyConfigStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// DexProxyConfig is the Schema for the dexproxyconfigs API
type DexProxyConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DexProxyConfigSpec   `json:"spec,omitempty"`
	Status DexProxyConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DexProxyConfigList contains a list of DexProxyConfig
type DexProxyConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DexProxyConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DexProxyConfig{}, &DexProxyConfigList{})
}
