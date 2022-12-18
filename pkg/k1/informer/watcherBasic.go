package informer

import (
	"fmt"

	"github.com/kubefirst/kubefirst-watcher/pkg/k1/crd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

// TODO: Make this more generic

func WatchBasic(conditions []crd.BasicConfigurationCondition, matchConditions chan Condition, stopper chan struct{}, informer cache.SharedIndexInformer) {
	logger.Debug(fmt.Sprintf("Started Wacher for %#v", conditions))

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mObj := obj.(metav1.Object)
			labels := mObj.GetLabels()
			logger.Debug(fmt.Sprintf("New Basic updated: %s, %s", mObj.GetName(), mObj.GetNamespace()))
			CheckMatchBasicConfigurationCondition(&crd.BasicConfigurationCondition{Namespace: mObj.GetNamespace(), Name: mObj.GetName()}, labels, conditions, matchConditions)

		},
		UpdateFunc: func(old, new interface{}) {
			newObj := new.(metav1.Object)
			labels := newObj.GetLabels()
			logger.Debug(fmt.Sprintf("Basic updated: %s, %s", newObj.GetName(), newObj.GetNamespace()))
			CheckMatchBasicConfigurationCondition(&crd.BasicConfigurationCondition{Namespace: newObj.GetNamespace(), Name: newObj.GetName()}, labels, conditions, matchConditions)
		},
		DeleteFunc: func(obj interface{}) {
			mObj := obj.(metav1.Object)
			logger.Debug(fmt.Sprintf("New Basic deleted from Store: %s", mObj.GetName()))
		},
	})
	informer.Run(stopper)
}

func CheckMatchBasicConfigurationCondition(obj *crd.BasicConfigurationCondition, labelsFound map[string]string, conditions []crd.BasicConfigurationCondition, matchCondition chan Condition) {
	//check on conditions list if there is a match
	for k, _ := range conditions {
		propertyExpected := ExtractBasicConfigurationMap(&conditions[k])
		propertyFound := ExtractBasicConfigurationMap(obj)
		if MatchesGeneric(&propertyFound, &labelsFound, &propertyExpected, &conditions[k].Labels) {
			logger.Debug(fmt.Sprintf("Interest BasicConfigurationCondition event found -  status: %#v", obj))
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

//ExtractBasicConfigurationMap - converts BasicConfigurationCondition to Map
func ExtractBasicConfigurationMap(obj *crd.BasicConfigurationCondition) map[string]string {
	result := map[string]string{}
	if len(obj.Name) > 0 {
		result["name"] = obj.Name
	}
	if len(obj.Namespace) > 0 {
		result["namespace"] = obj.Namespace
	}

	return result
}
