package informer

import (
	"fmt"

	"github.com/kubefirst/kubefirst-watcher/pkg/k1/crd"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

func WatchJobs(conditions []crd.JobCondition, matchConditions chan Condition, stopper chan struct{}, informer cache.SharedIndexInformer) {
	logger.Debug(fmt.Sprintf("Started Wacher for %#v", conditions))
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(*batchv1.Job)
			labels := obj.(*corev1.Pod).Labels

			logger.Debug(fmt.Sprintf("New Pod updated: %s, %s, %s", mObj.GetName(), mObj.GetNamespace(), mObj.Status.Succeeded))
			checkMatchConditionJob(mObj, labels, conditions, matchConditions)

		},
		UpdateFunc: func(old, new interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			newObj := new.(*batchv1.Job)
			labels := new.(*corev1.Pod).Labels
			logger.Debug(fmt.Sprintf("Pod updated: %s, %s, %s", newObj.GetName(), newObj.GetNamespace(), newObj.Status.Succeeded))
			checkMatchConditionJob(newObj, labels, conditions, matchConditions)
		},
		DeleteFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(*batchv1.Job)
			logger.Debug(fmt.Sprintf("New Pod deleted from Store: %s", mObj.GetName()))
		},
	})
	informer.Run(stopper)
}

func checkMatchConditionJob(obj *batchv1.Job, labels map[string]string, conditions []crd.JobCondition, matchCondition chan Condition) {
	//check on conditions list if there is a match
	for k, _ := range conditions {
		if obj.Namespace == conditions[k].Namespace &&
			obj.Name == conditions[k].Name &&
			obj.Status.Succeeded == conditions[k].Succeeded {
			matchMap, _ := IsMapPresent(labels, conditions[k].Labels)
			if matchMap {
				logger.Debug(fmt.Sprintf("Interest Job event found -  status:  %s, %s, %s", obj.GetName(), obj.GetNamespace(), obj.Status.Succeeded))
				foundCondition := Condition{
					ID:  conditions[k].ID,
					Met: true,
				}
				logger.Debug(fmt.Sprintf("Sending Condition -  status:  %#v ", foundCondition))
				matchCondition <- foundCondition
				//Remove Condition found
				//https://github.com/golang/go/wiki/SliceTricks
				// conditions = append(conditions[:k], conditions[k+1:]...)
				// it may fail on nil scenarios - extra checks needed
				//This need to be global, as this checks may run in parallel.
				//TODO: need to find an list that is thread safe
				logger.Debug(fmt.Sprintf("Remaning Condition -  status:  %#v ", foundCondition))
			}

		}
	}
}
