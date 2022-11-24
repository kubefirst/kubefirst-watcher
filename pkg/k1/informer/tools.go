package informer

import (
	"context"
	"fmt"
	"reflect"

	"github.com/thoas/go-funk"
	api "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func getK8SConfig() *kubernetes.Clientset {
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

func IsMapPresent(sourceMap map[string]string, subsetMap map[string]string) (bool, error) {
	match := true
	keysAll := funk.Keys(sourceMap)
	keysSubset := funk.Keys(subsetMap)

	intersect := funk.Intersect(keysAll, keysSubset)
	if !reflect.DeepEqual(intersect, keysSubset) {
		match = false
	}
	funk.ForEach(intersect, func(x string) {
		if sourceMap[x] != subsetMap[x] {
			match = false
		}
	})

	return match, nil
}

type PatchObject struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

func UpdateStatus(watcherConfig *WatcherConfig) error {
	clientSet := getK8SConfig()
	myPatch := `{"status":{"status":"change"}}`
	object, err := clientSet.RESTClient().
		Patch(api.MergePatchType).
		SubResource("status").
		Namespace("default").
		Resource("watcher").
		Name("watcher-sample-01").
		Body([]byte(myPatch)).
		Do(context.TODO()).
		Get()
	logger.Info(fmt.Sprintf("Update err:  %#v ", object))
	logger.Info(fmt.Sprintf("Update err:  %#v ", err))
	logger.Info(fmt.Sprintf("Update status:  %#v ", watcherConfig))
	return nil
}

/*
result, err6 := tprclient.Patch(api.JSONPatchType).
        Namespace(api.NamespaceDefault).
        Resource("pgupgrades").
        Name("junk").
        Body(patchBytes).
        Do().
        Get()
*/
