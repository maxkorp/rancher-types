package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type PrincipalLifecycle interface {
	Create(obj *Principal) (runtime.Object, error)
	Remove(obj *Principal) (runtime.Object, error)
	Updated(obj *Principal) (runtime.Object, error)
}

type principalLifecycleAdapter struct {
	lifecycle PrincipalLifecycle
}

func (w *principalLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*Principal))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *principalLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*Principal))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *principalLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*Principal))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewPrincipalLifecycleAdapter(name string, clusterScoped bool, client PrincipalInterface, l PrincipalLifecycle) PrincipalHandlerFunc {
	adapter := &principalLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *Principal) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
