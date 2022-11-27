/*
@Time : 2022/11/22 18:16
@Author : lianyz
@Description :
*/

package v1

import (
	"github.com/lianyz/k8s-controller-custom-resource/pkg/apis/samplecrd"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	api "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SchemeGroupVersion is the identifier for the API which includes
// the name of the group and the version of the API
var SchemeGroupVersion = schema.GroupVersion{
	Group:   samplecrd.GroupName,
	Version: samplecrd.Version,
}

// create a SchemeBuilder which use functions to add types to the scheme
var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// Kind takes an unqualified kind and returns a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(
		SchemeGroupVersion,
		&Network{},
		&NetworkList{})

	api.AddToGroupVersion(scheme, SchemeGroupVersion)

	return nil
}
