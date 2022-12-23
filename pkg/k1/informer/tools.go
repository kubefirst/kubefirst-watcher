package informer

import (
	"fmt"
	"reflect"

	"github.com/thoas/go-funk"
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

func IsMapPresent(sourceMap *map[string]string, subsetMap *map[string]string) (bool, error) {
	match := true
	keysAll := funk.Keys(sourceMap)
	keysSubset := funk.Keys(subsetMap)

	intersect := funk.Intersect(keysAll, keysSubset)
	if !reflect.DeepEqual(intersect, keysSubset) {
		match = false
	}
	funk.ForEach(intersect, func(x string) {
		if (*sourceMap)[x] != (*subsetMap)[x] {
			match = false
		}
	})

	return match, nil
}

// MatchesGeneric Verify if object found matches expected conditions
// Isolated to help on re-use of the logic
func MatchesGeneric(propertyFound *map[string]string, labelsFound *map[string]string, propertyExpected *map[string]string, labelsExpected *map[string]string) bool {
	matchCore, _ := IsMapPresent(propertyFound, propertyExpected)
	if !matchCore {
		return false
	}
	matchLabels, _ := IsMapPresent(labelsFound, labelsExpected)
	if !matchLabels {
		return false
	}
	return true
}
