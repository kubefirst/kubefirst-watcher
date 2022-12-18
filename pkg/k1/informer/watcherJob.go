package informer

import (
	"fmt"
	"strconv"

	"github.com/kubefirst/kubefirst-watcher/pkg/k1/crd"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/tools/cache"
)

func WatchJobs(conditions []crd.JobCondition, matchConditions chan Condition, stopper chan struct{}, informer cache.SharedIndexInformer) {
	logger.Debug(fmt.Sprintf("Started Wacher for %#v", conditions))
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mObj := obj.(*batchv1.Job)
			labels := obj.(*batchv1.Job).Labels

			logger.Debug(fmt.Sprintf("New Job updated: %s, %s, %d", mObj.GetName(), mObj.GetNamespace(), mObj.Status.Succeeded))
			CheckMatchConditionJob(mObj, labels, conditions, matchConditions)

		},
		UpdateFunc: func(old, new interface{}) {
			newObj := new.(*batchv1.Job)
			labels := new.(*batchv1.Job).Labels
			logger.Debug(fmt.Sprintf("Job updated: %s, %s, %d", newObj.GetName(), newObj.GetNamespace(), newObj.Status.Succeeded))
			CheckMatchConditionJob(newObj, labels, conditions, matchConditions)
		},
		DeleteFunc: func(obj interface{}) {
			mObj := obj.(*batchv1.Job)
			logger.Debug(fmt.Sprintf("New Job deleted from Store: %s", mObj.GetName()))
		},
	})
	informer.Run(stopper)
}

func CheckMatchConditionJob(obj *batchv1.Job, labelsFound map[string]string, conditions []crd.JobCondition, matchCondition chan Condition) {
	//check on conditions list if there is a match
	for k, _ := range conditions {
		propertyExpected := ExtractJobConditionMap(&conditions[k])
		propertyFound := ExtractJobMap(obj)
		if MatchesGeneric(&propertyFound, &labelsFound, &propertyExpected, &conditions[k].Labels) {
			logger.Debug(fmt.Sprintf("Interest Job event found -  status:  %s, %s, %d", obj.GetName(), obj.GetNamespace(), obj.Status.Succeeded))
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

//ExtractJobMap - Convert Job to Map
func ExtractJobMap(obj *batchv1.Job) map[string]string {
	result := map[string]string{}
	if len(obj.Name) > 0 {
		result["name"] = obj.Name
	}
	if len(obj.Namespace) > 0 {
		result["namespace"] = obj.Namespace
	}
	if obj.Status.Succeeded > 0 {
		result["succeeded"] = strconv.FormatInt(int64(obj.Status.Succeeded), 10)
	}

	return result
}

//ExtractJobConditionMap - Converts JobCondition to Map
func ExtractJobConditionMap(obj *crd.JobCondition) map[string]string {
	result := map[string]string{}
	if len(obj.Name) > 0 {
		result["name"] = obj.Name
	}
	if len(obj.Namespace) > 0 {
		result["namespace"] = obj.Namespace
	}
	if obj.Succeeded > 0 {
		result["succeeded"] = strconv.FormatInt(int64(obj.Succeeded), 10)
	}

	return result
}
