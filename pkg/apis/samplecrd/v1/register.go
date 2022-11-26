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

var SchemeGroupVersion = schema.GroupVersion{
	Group:   samplecrd.GroupName,
	Version: samplecrd.Version,
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(
		SchemeGroupVersion,
		&Network{},
		&NetworkList{})

	api.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil

}
