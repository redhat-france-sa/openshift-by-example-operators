/*
MIT License

Copyright (c) 2020 RH France Solution Architects

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FruitsCatalogGSpec defines the desired state of FruitsCatalogG
type FruitsCatalogGSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:default:=fruits-catalog
	AppName string `json:"appName,omitempty"`
	// +kubebuilder:validation:Optional
	WebApp  WebAppSpec  `json:"webapp,omitempty"`
	MongoDB MongoDBSpec `json:"mongodb,omitempty"`
}

// WebAppSpec defines the desired state of WebApp
// +k8s:openapi-gen=true
type WebAppSpec struct {
	ReplicaCount int32 `json:"replicaCount,omitempty"`
	// +kubebuilder:default:="quay.io/lbroudoux/fruits-catalog:latest"
	Image   string      `json:"image,omitempty"`
	Ingress IngressSpec `json:"ingress,omitempty"`
}

// IngressSpec defines the desired state of WebApp Ingress
// +k8s:openapi-gen=true
type IngressSpec struct {
	// +kubebuilder:default:=true
	Enabled bool `json:"enabled,omitempty" default:"true"`
}

// MongoDBSpec defines the desired state of MongoDB
// +k8s:openapi-gen=true
type MongoDBSpec struct {
	// +kubebuilder:default:=true
	Install bool `json:"install,omitempty"`
	// +kubebuilder:default:="centos/mongodb-34-centos7:latest"
	Image    string `json:"image,omitempty"`
	URI      string `json:"uri,omitempty"`
	Database string `json:"database,omitempty"`
	// +kubebuilder:default:=true
	Persistent bool `json:"persistent,omitempty"`
	// +kubebuilder:default:="2Gi"
	VolumeSize string `json:"volumeSize,omitempty"`
	// +kubebuilder:default:="myusername"
	Username string `json:"username,omitempty"`
	// +kubebuilder:default:="mypassword"
	Password  string        `json:"password,omitempty"`
	SecretRef SecretRefSpec `json:"secretRef,omitempty"`
}

// SecretRefSpec defines a reference to a Secret
// +k8s:openapi-gen=true
type SecretRefSpec struct {
	Secret      string `json:"secret"`
	UsernameKey string `json:"usernameKey"`
	PasswordKey string `json:"passwordKey"`
}

// FruitsCatalogGStatus defines the observed state of FruitsCatalogG
type FruitsCatalogGStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	WebApp  string `json:"webapp"`
	MongoDB string `json:"mongodb"`
	Secret  string `json:"secret"`
	Route   string `json:"route"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// FruitsCatalogG is the Schema for the fruitscataloggs API
type FruitsCatalogG struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FruitsCatalogGSpec   `json:"spec,omitempty"`
	Status FruitsCatalogGStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FruitsCatalogGList contains a list of FruitsCatalogG
type FruitsCatalogGList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FruitsCatalogG `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FruitsCatalogG{}, &FruitsCatalogGList{})
}
