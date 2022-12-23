package informer

import (
	"fmt"

	"github.com/kubefirst/kubefirst-watcher/pkg/k1/crd"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

func WatchPods(conditions []crd.PodCondition, matchConditions chan Condition, stopper chan struct{}, informer cache.SharedIndexInformer) {
	logger.Debug(fmt.Sprintf("Started Wacher for %#v", conditions))

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mObj := obj.(*corev1.Pod)
			labels := obj.(*corev1.Pod).Labels

			logger.Debug(fmt.Sprintf("New Pod updated: %s, %s, %s", mObj.GetName(), mObj.GetNamespace(), mObj.Status.Phase))
			CheckMatchConditionPod(mObj, labels, conditions, matchConditions)

		},
		UpdateFunc: func(old, new interface{}) {
			newObj := new.(*corev1.Pod)
			labels := new.(*corev1.Pod).Labels
			logger.Debug(fmt.Sprintf("Pod updated: %s, %s, %s", newObj.GetName(), newObj.GetNamespace(), newObj.Status.Phase))
			CheckMatchConditionPod(newObj, labels, conditions, matchConditions)
		},
		DeleteFunc: func(obj interface{}) {
			mObj := obj.(*corev1.Pod)
			logger.Debug(fmt.Sprintf("New Pod deleted from Store: %s", mObj.GetName()))
		},
	})
	informer.Run(stopper)
}

func CheckMatchConditionPod(obj *corev1.Pod, labelsFound map[string]string, conditions []crd.PodCondition, matchCondition chan Condition) {
	//check on conditions list if there is a match
	for k, _ := range conditions {
		propertyExpected := ExtractPodConditionMap(&conditions[k])
		propertyFound := ExtractPodMap(obj)
		if MatchesGeneric(&propertyFound, &labelsFound, &propertyExpected, &conditions[k].Labels) {
			logger.Debug(fmt.Sprintf("Interest Pod event found -  status:  %s, %s, %s", obj.GetName(), obj.GetNamespace(), obj.Status.Phase))
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

//ExtractPodMap - Convert Pod to Map
func ExtractPodMap(obj *corev1.Pod) map[string]string {
	result := map[string]string{}
	if len(obj.Name) > 0 {
		result["name"] = obj.Name
	}
	if len(obj.Namespace) > 0 {
		result["namespace"] = obj.Namespace
	}
	if len(string(obj.Status.Phase)) > 0 {
		result["phase"] = string(obj.Status.Phase)
	}

	return result
}

//ExtractPodConditionMap - Convert PodCondition to Map
func ExtractPodConditionMap(obj *crd.PodCondition) map[string]string {
	result := map[string]string{}
	if len(obj.Name) > 0 {
		result["name"] = obj.Name
	}
	if len(obj.Namespace) > 0 {
		result["namespace"] = obj.Namespace
	}
	if len(obj.Phase) > 0 {
		result["phase"] = obj.Phase
	}

	return result
}
