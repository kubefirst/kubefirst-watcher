package informer

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

// TODO: Make this more generic

func WatchConfigMap(conditions []BasicConfigurationCondition, matchConditions chan Condition, stopper chan struct{}) {
	logger.Debug(fmt.Sprintf("Started Wacher for %#v", conditions))
	clientSet := getK8SConfig()
	factory := informers.NewSharedInformerFactory(clientSet, 0)
	informer := factory.Core().V1().ConfigMaps().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(*metav1.ObjectMeta)
			logger.Debug(fmt.Sprintf("New Pod updated:", mObj.GetName(), mObj.GetNamespace()))
			checkMatchBasicConfigurationCondition(&BasicConfiguration{Namespace: mObj.GetNamespace(), Name: mObj.GetName()}, conditions, matchConditions)

		},
		UpdateFunc: func(old, new interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			newObj := new.(*metav1.ObjectMeta)
			logger.Debug(fmt.Sprintf("Pod updated:", newObj.GetName(), newObj.GetNamespace()))
			checkMatchBasicConfigurationCondition(&BasicConfiguration{Namespace: newObj.GetNamespace(), Name: newObj.GetName()}, conditions, matchConditions)
		},
		DeleteFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(*metav1.ObjectMeta)
			logger.Debug(fmt.Sprintf("New Pod deleted from Store: %s", mObj.GetName()))
		},
	})
	informer.Run(stopper)
}

func WatchSecrets(conditions []BasicConfigurationCondition, matchConditions chan Condition, stopper chan struct{}) {
	logger.Debug(fmt.Sprintf("Started Wacher for %#v", conditions))
	clientSet := getK8SConfig()
	factory := informers.NewSharedInformerFactory(clientSet, 0)
	informer := factory.Core().V1().Secrets().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(*corev1.Secret)
			logger.Debug(fmt.Sprintf("New Pod updated:", mObj.GetName(), mObj.GetNamespace()))
			checkMatchBasicConfigurationCondition(&BasicConfiguration{Namespace: mObj.GetNamespace(), Name: mObj.GetName()}, conditions, matchConditions)

		},
		UpdateFunc: func(old, new interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			newObj := new.(*corev1.Secret)
			logger.Debug(fmt.Sprintf("Pod updated:", newObj.GetName(), newObj.GetNamespace()))
			checkMatchBasicConfigurationCondition(&BasicConfiguration{Namespace: newObj.GetNamespace(), Name: newObj.GetName()}, conditions, matchConditions)
		},
		DeleteFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(*corev1.Secret)
			logger.Debug(fmt.Sprintf("New Pod deleted from Store: %s", mObj.GetName()))
		},
	})
	informer.Run(stopper)
}

func WatchServices(conditions []BasicConfigurationCondition, matchConditions chan Condition, stopper chan struct{}) {
	logger.Debug(fmt.Sprintf("Started Wacher for %#v", conditions))
	clientSet := getK8SConfig()
	factory := informers.NewSharedInformerFactory(clientSet, 0)
	informer := factory.Core().V1().Services().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(*corev1.Service)
			logger.Debug(fmt.Sprintf("New Pod updated:", mObj.GetName(), mObj.GetNamespace()))
			checkMatchBasicConfigurationCondition(&BasicConfiguration{Namespace: mObj.GetNamespace(), Name: mObj.GetName()}, conditions, matchConditions)

		},
		UpdateFunc: func(old, new interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			newObj := new.(*corev1.Service)
			logger.Debug(fmt.Sprintf("Pod updated:", newObj.GetName(), newObj.GetNamespace()))
			checkMatchBasicConfigurationCondition(&BasicConfiguration{Namespace: newObj.GetNamespace(), Name: newObj.GetName()}, conditions, matchConditions)
		},
		DeleteFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(*corev1.Service)
			logger.Debug(fmt.Sprintf("New Pod deleted from Store: %s", mObj.GetName()))
		},
	})
	informer.Run(stopper)
}

func checkMatchBasicConfigurationCondition(obj *BasicConfiguration, conditions []BasicConfigurationCondition, matchCondition chan Condition) {
	//check on conditions list if there is a match
	for k, _ := range conditions {
		if obj.Namespace == conditions[k].Namespace &&
			obj.Name == conditions[k].Name {
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
