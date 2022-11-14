package informer

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

var jobGVK = schema.GroupVersionKind{Group: "", Version: "batch/v1", Kind: "Job"}

func WatchJods(condition Condition, interestingPods chan Condition, stopper chan struct{}) {
	logger.Debug(fmt.Sprintf("Started Wacher for %#v", condition))
	clientSet := getK8SConfig()
	factory := informers.NewSharedInformerFactory(clientSet, 0)
	informer := factory.Batch().V1().Jobs().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(*batchv1.Job)
			logger.Debug(fmt.Sprintf("New Job updated:", mObj.GetName(), mObj.GetNamespace(), mObj.Status))
			checkMatchConditionJob(mObj, condition, interestingPods)

		},
		UpdateFunc: func(old, new interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			newObj := new.(*batchv1.Job)
			logger.Debug(fmt.Sprintf("Job updated:", newObj.GetName(), newObj.GetNamespace(), newObj.Status))
			checkMatchConditionJob(newObj, condition, interestingPods)
		},
		DeleteFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(*batchv1.Job)
			logger.Debug(fmt.Sprintf("New Job deleted from Store: %s", mObj.GetName()))
		},
	})
	informer.Run(stopper)
}

func checkMatchConditionJob(obj *batchv1.Job, condition Condition, interestingPods chan Condition) {
	if string(obj.Status) == condition.Phase {
		logger.Debug(fmt.Sprintf("\nInterest Pod event found -  status: ", obj.GetName(), obj.GetNamespace(), obj.Status.Phase))
		foundCondition := Condition{
			Namespace:  obj.GetNamespace(),
			Name:       obj.GetName(),
			Phase:      string(obj.Status),
			APIVersion: condition.APIVersion, //default value is blank - obj.APIVersion
			Kind:       condition.Kind,       //default value is blank - obj.Kind
		}
		logger.Debug(fmt.Sprintf("\nSending Condition -  status: ", foundCondition.Namespace, foundCondition.Name, obj.Status.Phase, obj.APIVersion, obj.Kind))
		interestingPods <- foundCondition
	}
}
