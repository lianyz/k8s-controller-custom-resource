/*
@Time : 2022/11/22 18:16
@Author : lianyz
@Description :
*/

package v1

import (
	api "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Network describes a Network resource
type Network struct {
	// TypeMeta is the metadata for the resource, like kind and apiversion
	api.TypeMeta `json:",inline"`

	// ObjectMeta contains the metadata for the particular object, including
	// things like...
	// - name
	// - namespace
	// - labels
	// - annotations
	// - self link
	api.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the custom resource spec
	Spec NetworkSpec `json:"spec"`
}

// NetworkSpec is the spec for a Network resource
type NetworkSpec struct {
	Cidr    string `json:"cidr"`
	Gateway string `json:"gateway"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkList is a list of Network resources
type NetworkList struct {
	api.TypeMeta `json:",inline"`
	api.ListMeta `json:"metadata"`

	Items []Network `json:"items"`
}
