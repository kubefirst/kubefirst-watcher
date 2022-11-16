package informer

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func WatchPods(conditions []PodCondition, interestingPods chan Condition, stopper chan struct{}) {
	logger.Debug(fmt.Sprintf("Started Wacher for %#v", conditions))
	clientSet := getK8SConfig()
	factory := informers.NewSharedInformerFactory(clientSet, 0)
	informer := factory.Core().V1().Pods().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(*corev1.Pod)
			logger.Debug(fmt.Sprintf("New Pod updated:", mObj.GetName(), mObj.GetNamespace(), mObj.Status.Phase))
			checkMatchConditionPod(mObj, conditions, interestingPods)

		},
		UpdateFunc: func(old, new interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			newObj := new.(*corev1.Pod)
			logger.Debug(fmt.Sprintf("Pod updated:", newObj.GetName(), newObj.GetNamespace(), newObj.Status.Phase))
			checkMatchConditionPod(newObj, conditions, interestingPods)
		},
		DeleteFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(*corev1.Pod)
			logger.Debug(fmt.Sprintf("New Pod deleted from Store: %s", mObj.GetName()))
		},
	})
	informer.Run(stopper)
}

func checkMatchConditionPod(obj *corev1.Pod, conditions []PodCondition, matchCondition chan Condition) {
	//check on conditions list if there is a match
	for k, _ := range conditions {
		if obj.Namespace == conditions[k].Namespace &&
			obj.Name == conditions[k].Name &&
			string(obj.Status.Phase) == conditions[k].Phase {
			logger.Debug(fmt.Sprintf("Interest Pod event found -  status: ", obj.GetName(), obj.GetNamespace(), obj.Status.Phase))
			foundCondition := Condition{
				ID:  conditions[k].ID,
				Met: true,
			}
			logger.Debug(fmt.Sprintf("Sending Condition -  status:  %#v ", foundCondition))
			matchCondition <- foundCondition
			//Remove Condition found
			//https://github.com/golang/go/wiki/SliceTricks
			conditions = append(conditions[:k], conditions[k+1:]...)
			//This need to be global, as this checks may run in parallel.
			//TODO: need to find an list that is thread safe
			logger.Debug(fmt.Sprintf("Remaning Condition -  status:  %#v ", foundCondition))

		}
	}
}
