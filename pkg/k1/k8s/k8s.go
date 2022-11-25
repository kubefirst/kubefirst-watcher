package k8s

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetK8SConfig() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	fmt.Println("Config used")
	return clientSet

}
