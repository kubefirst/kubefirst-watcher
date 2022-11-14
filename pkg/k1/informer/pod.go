package informer

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func WatchPods(condition Condition, interestingPods chan Condition, stopper chan struct{}) {
	clientSet := getK8SConfig()
	factory := informers.NewSharedInformerFactory(clientSet, 0)
	informer := factory.Core().V1().Pods().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(*corev1.Pod)
			//fmt.Printf("New Pod Added to Store: %s - %v", mObj.Name, mObj.ObjectMeta.Labels)
			fmt.Println("\nNew Pod updated:", mObj.GetName(), mObj.GetNamespace(), mObj.Status.Phase)
			checkMatchConditionPod(mObj, condition, interestingPods)

		},
		UpdateFunc: func(old, new interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			newObj := new.(*corev1.Pod)
			fmt.Println("\nPod updated:", newObj.GetName(), newObj.GetNamespace(), newObj.Status.Phase)
			checkMatchConditionPod(newObj, condition, interestingPods)
		},
		DeleteFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(*corev1.Pod)
			fmt.Printf("\nNew Pod deleted from Store: %s", mObj.GetName())
		},
	})
	informer.Run(stopper)
}

func checkMatchConditionPod(obj *corev1.Pod, condition Condition, interestingPods chan Condition) {
	if string(obj.Status.Phase) == condition.Phase {
		fmt.Println("\nInterest Pod event found -  status: ", obj.GetName(), obj.GetNamespace(), obj.Status.Phase)
		foundCondition := Condition{
			Namespace:  obj.GetNamespace(),
			Name:       obj.GetName(),
			Phase:      string(obj.Status.Phase),
			APIVersion: condition.APIVersion, //default value is blank - obj.APIVersion
			Kind:       condition.Kind,       //default value is blank - obj.Kind
		}
		fmt.Println("\nSending Condition -  status: ", foundCondition.Namespace, foundCondition.Name, obj.Status.Phase, obj.APIVersion, obj.Kind)
		interestingPods <- foundCondition
	}
}
