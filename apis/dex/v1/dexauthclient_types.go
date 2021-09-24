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

// DexSecretRef defines the k8s secret ref of Dex
type DexSecretRef struct {
	// Namespace of the k8s secret
	NameSpace string `json:"namespace"`
	// Name of the the k8s secret
	Name string `json:"name"`
	// The config file key in the k8s secret
	Key string `json:"key"`
}

// DexAuthClientSpec defines the desired state of DexAuthClient
type DexAuthClientSpec struct {
	Id           string   `json:"id"           yaml:"id"`
	Name         string   `json:"name"         yaml:"name"`
	RedirectURIs []string `json:"redirectURIs" yaml:"redirectURIs,flow"`
	Secret       string   `json:"secret"       yaml:"secret"`
	// The K8s Secret Ref would be modified
	SecretRef DexSecretRef `json:"secretRef"   yaml:"-"`
}

// DexAuthClientStatus defines the observed state of DexAuthClient
type DexAuthClientStatus struct {
	Success bool   `json:"success"`
	Finish  bool   `json:"finish"`
	Message string `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DexAuthClient is the Schema for the dexauthclients API
type DexAuthClient struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DexAuthClientSpec   `json:"spec,omitempty"`
	Status DexAuthClientStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DexAuthClientList contains a list of DexAuthClient
type DexAuthClientList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DexAuthClient `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DexAuthClient{}, &DexAuthClientList{})
}

func (status *DexAuthClientStatus) Reset() {
	status.Finish = false
	status.Success = false
	status.Message = ""
}

func (status *DexAuthClientStatus) Error(message string) {
	status.Finish = true
	status.Success = false
	status.Message = message
}

func (status *DexAuthClientStatus) DoSuccess() {
	status.Finish = true
	status.Success = true
	status.Message = ""
}
