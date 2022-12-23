package informer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kubefirst/kubefirst-watcher/pkg/k1/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/tools/cache"
)

func WatchDeployments(conditions []v1beta1.DeploymentCondition, matchConditions chan Condition, stopper chan struct{}, informer cache.SharedIndexInformer) {
	logger.Debug(fmt.Sprintf("Started Wacher for %#v", conditions))

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mObj := obj.(*appsv1.Deployment)
			labels := obj.(*appsv1.Deployment).Labels
			CheckMatchConditionDeployment(mObj, labels, conditions, matchConditions)

		},
		UpdateFunc: func(old, new interface{}) {
			newObj := new.(*appsv1.Deployment)
			labels := new.(*appsv1.Deployment).Labels
			CheckMatchConditionDeployment(newObj, labels, conditions, matchConditions)
		},
		DeleteFunc: func(obj interface{}) {

		},
	})
	informer.Run(stopper)
}

func CheckMatchConditionDeployment(obj *appsv1.Deployment, labelsFound map[string]string, conditions []v1beta1.DeploymentCondition, matchCondition chan Condition) {
	//check on conditions list if there is a match
	for k, _ := range conditions {
		propertyExpected := ExtractDeploymentConditionMap(&conditions[k])
		propertyFound := ExtractDeploymentMap(obj)
		if MatchesGeneric(&propertyFound, &labelsFound, &propertyExpected, &conditions[k].Labels) {
			foundCondition := Condition{
				ID:  conditions[k].ID,
				Met: true,
			}
			matchCondition <- foundCondition
		}

	}
}

//ExtractPodMap - Convert Pod to Map
func ExtractDeploymentMap(obj *appsv1.Deployment) map[string]string {
	result := map[string]string{}
	if len(obj.Name) > 0 {
		result["name"] = obj.Name
	}
	if len(obj.Namespace) > 0 {
		result["namespace"] = obj.Namespace
	}
	result["replicas"] = string(obj.Status.Replicas)
	if *obj.Spec.Replicas == obj.Status.ReadyReplicas {
		result["ready"] = fmt.Sprintf("%t", true)
	} else {
		result["ready"] = fmt.Sprintf("%t", false)
	}

	return result
}

//ExtractPodConditionMap - Convert PodCondition to Map
func ExtractDeploymentConditionMap(obj *v1beta1.DeploymentCondition) map[string]string {
	result := map[string]string{}
	if len(obj.Name) > 0 {
		result["name"] = obj.Name
	}
	if len(obj.Namespace) > 0 {
		result["namespace"] = obj.Namespace
	}
	if obj.Replicas > 0 {
		result["replicas"] = strconv.FormatInt(int64(obj.Replicas), 10)
	}
	//This need to be string to check empty
	if len(obj.Ready) > 0 {
		result["ready"] = strings.ToLower(obj.Ready)
	}

	return result
}
