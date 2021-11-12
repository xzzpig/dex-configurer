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
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type UrlScheme string

const (
	SchemeHTTP  UrlScheme = "http"
	SchemeHTTPS UrlScheme = "https"
)

type DexRedirectUrl struct {
	// +kubebuilder:default:https
	// +kubebuilder:validation:Enum=http;https
	// +kubebuilder:validation:Optional
	Scheme UrlScheme `json:"scheme,omitempty"`

	// If Empty will Find Host from Target Ingress
	// +kubebuilder:validation:Optional
	Host string `json:"host,omitempty"`

	// +kubebuilder:default:443
	// +kubebuilder:validation:Optional
	Port int `json:"port,omitempty"`

	// +kubebuilder:default:/oauth2
	// +kubebuilder:validation:Optional
	Path string `json:"path,omitempty"`
}

type NamespacedName struct {
	// +kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name"`
}

// DexClientOrderSpec defines the desired state of DexClientOrder
type DexClientOrderSpec struct {
	// Config Ref of DexClientOrder
	Config NamespacedName `json:"config"`

	// Target Ingress will Add Nginx Auth Annotations
	TargetIngress NamespacedName `json:"target-ingress"`

	// +kubebuilder:validation:Optional
	RedirectUrl DexRedirectUrl `json:"redirect-url"`

	// +kubebuilder:validation:Optional
	ClientId string `json:"client-id,omitempty"`
	// +kubebuilder:validation:Optional
	ClientName string `json:"client-name,omitempty"`
	// +kubebuilder:validation:Optional
	ClientSecret string `json:"client-secret,omitempty"`

	// Groups Allow user to login, All groups if null
	// +kubebuilder:validation:Optional
	AllowedGroups []string `json:"allowed-groups"`
}

type DexClientOrderRefObjects struct {
	ClientRef     NamespacedName `json:"client-ref,omitempty"`
	DeploymentRef NamespacedName `json:"deployment-ref,omitempty"`
	ServiceRef    NamespacedName `json:"service-ref,omitempty"`
	IngressRef    NamespacedName `json:"ingress-ref,omitempty"`
}

// DexClientOrderStatus defines the observed state of DexClientOrder
type DexClientOrderStatus struct {
	Created    bool                     `json:"created,omitempty"`
	Message    string                   `json:"message,omitempty"`
	RefObjects DexClientOrderRefObjects `json:"ref-objects,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DexClientOrder is the Order to create a dex auth client
type DexClientOrder struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DexClientOrderSpec   `json:"spec,omitempty"`
	Status DexClientOrderStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DexClientOrderList contains a list of DexClientOrder
type DexClientOrderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DexClientOrder `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DexClientOrder{}, &DexClientOrderList{})
}

func (url *DexRedirectUrl) GoString() string {
	portStr := ""
	path := url.GetPath()
	if (url.Scheme == SchemeHTTPS && url.Port != 443) || (url.Scheme == SchemeHTTP && url.Port != 80) {
		portStr = ":" + fmt.Sprint(url.Port)
	}
	return fmt.Sprintf("%v://%v%v%v/callback", url.Scheme, url.Host, portStr, path)
}

func (url *DexRedirectUrl) GetPath() string {
	return strings.TrimRight(url.Path, "/")
}

func (url *DexRedirectUrl) AuthSignin() string {
	portStr := ""
	path := url.GetPath()
	if (url.Scheme == SchemeHTTPS && url.Port != 443) || (url.Scheme == SchemeHTTP && url.Port != 80) {
		portStr = ":" + fmt.Sprint(url.Port)
	}
	return fmt.Sprintf("%v://%v%v%v/start?rd=$escaped_request_uri", url.Scheme, "$host", portStr, path)
}
