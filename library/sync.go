package library

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func SetControllerRef(owner metav1.Object, owned metav1.Object) error {
	return controllerutil.SetControllerReference(owner, owned, scheme)
}

// Ported from controllerutil
func OwnerRefsMatch(a, b *metav1.OwnerReference) bool {
	if a == nil || b == nil {
		return false
	}
	aGV, err := schema.ParseGroupVersion(a.APIVersion)
	if err != nil {
		return false
	}

	bGV, err := schema.ParseGroupVersion(b.APIVersion)
	if err != nil {
		return false
	}

	return aGV.Group == bGV.Group && a.Kind == b.Kind && a.Name == b.Name
}