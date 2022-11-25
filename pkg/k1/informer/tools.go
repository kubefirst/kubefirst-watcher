package informer

import (
	"context"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v2"
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

func UpdateStatus(ownerFile string, status string) error {
	if ownerFile == "" {
		logger.Info(fmt.Sprintf("No owner file provided, skip CRD update.   #%s ", ownerFile))
		return nil
	}
	watcherConfig, err := loadWatcherConfig(ownerFile)
	if err != nil {
		logger.Info(fmt.Sprintf("Error processing owner file   #%v ", err))
		return err
	}
	logger.Debug(fmt.Sprintf("Watcher Config: #%v ", watcherConfig))
	clientSet := getK8SConfig()
	myPatch := fmt.Sprintf(`{"status":{"status":"%s"}}`, status)
	logger.Debug(fmt.Sprintf("Watcher Patch: #%v ", myPatch))
	_, err = clientSet.RESTClient().
		Patch(api.MergePatchType).
		AbsPath("/apis/" + watcherConfig.APIVersion).
		SubResource("status").
		Namespace(watcherConfig.CrdNamespace).
		Resource("watchers").
		Name(watcherConfig.CrdName).
		Body([]byte(myPatch)).
		DoRaw(context.TODO())
	if err != nil {
		logger.Info(fmt.Sprintf("Error updating CRD   #%v ", err))
		return err
	}
	logger.Info(fmt.Sprintf("Update status:  %#v ", watcherConfig))
	return nil
}

func loadWatcherConfig(file string) (*WatcherConfig, error) {
	watcherConfig := &WatcherConfig{}
	logger.Debug("Loading config file:" + file)
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Info(fmt.Sprintf("yamlFile.Get err   #%v ", err))
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, watcherConfig)
	if err != nil {
		logger.Info(fmt.Sprintf("Unmarshal: %v", err))
		return nil, err
	}

	return watcherConfig, nil
}
