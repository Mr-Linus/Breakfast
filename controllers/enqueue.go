package controllers

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type EnqueueRequest struct {
	handler.EnqueueRequestForObject
}

func (e *EnqueueRequest) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	if evt.Meta == nil {
		return
	}
}

func (e *EnqueueRequest) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	if evt.MetaOld != nil && evt.MetaNew != nil {
		if _, ok := evt.MetaOld.GetLabels()["bread"]; ok {
			if !reflect.DeepEqual(evt.ObjectOld.(*v1.Pod).Status, evt.ObjectNew.(*v1.Pod).Status) {
				q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
					Name:      evt.MetaNew.GetName(),
					Namespace: evt.MetaNew.GetNamespace(),
				}})
			}
		}
	}
}

func (e *EnqueueRequest) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	if evt.Meta == nil {
		return
	}
	if _, ok := evt.Meta.GetLabels()["bread"]; ok {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      evt.Meta.GetName(),
			Namespace: evt.Meta.GetNamespace(),
		}})
	}
}
