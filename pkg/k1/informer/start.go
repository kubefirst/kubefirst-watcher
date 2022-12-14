package informer

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/kubefirst/kubefirst-watcher/pkg/k1/crd"
	"github.com/kubefirst/kubefirst-watcher/pkg/k1/k8s"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/informers"
)

// StartWatcher - starts watcher tooling
var logger *zap.Logger

const (
	StatusCompleted string = "Satisfied"
	StatusTimeout   string = "Timeout"
)

func StartCRDWatcher(clientCrd *crd.CRDClient, loggerIn *zap.Logger) error {
	logger = clientCrd.Logger
	myCRD, err := clientCrd.GetCRD()
	if err != nil {
		return fmt.Errorf("Error getting CRD")
	}
	//Setup channels
	interestingPods := make(chan Condition)
	defer close(interestingPods)
	stopper := make(chan struct{})
	defer close(stopper)

	//Process Conditions into goals
	exitScenario, exitScenarioState, _ := loadExitScenarioFromCRD(myCRD.Spec)
	logger.Info(fmt.Sprintf("%#v", exitScenario))
	logger.Info(fmt.Sprintf("%#v", exitScenarioState))
	//Process Conditions into watchers
	//Start Goals tracker
	go checkConditions(exitScenarioState, clientCrd, interestingPods, stopper)
	//Start Watchers
	//go WatchSecrets(exitScenario.Secrets, interestingPods, stopper)
	startWatchers(exitScenario, interestingPods, stopper)
	//Check Current State - to catch events pre-informers are started
	time.Sleep(time.Duration(exitScenario.Timeout) * time.Second)
	logger.Error("Timeout - Fail to match conditions")
	clientCrd.UpdateStatus(StatusTimeout)
	return fmt.Errorf("timeout - Failed to meet exit condition")
}

func StartWatcher(configFile string, ownerFile string, loggerIn *zap.Logger) error {
	return fmt.Errorf("Disabled mode")
	logger = loggerIn
	//Setup channels
	interestingPods := make(chan Condition)
	defer close(interestingPods)
	stopper := make(chan struct{})
	defer close(stopper)

	//Process Conditions into goals
	exitScenario, exitScenarioState, _ := loadExitScenario(configFile)
	logger.Info(fmt.Sprintf("%#v", exitScenario))
	logger.Info(fmt.Sprintf("%#v", exitScenarioState))
	//Process Conditions into watchers
	//Start Goals tracker
	//go checkConditions(exitScenarioState, ownerFile, interestingPods, stopper)
	//Start Watchers
	//go WatchSecrets(exitScenario.Secrets, interestingPods, stopper)
	startWatchers(exitScenario, interestingPods, stopper)
	//Check Current State - to catch events pre-informers are started
	time.Sleep(time.Duration(exitScenario.Timeout) * time.Second)
	logger.Error("Timeout - Fail to match conditions")
	UpdateStatus(ownerFile, StatusTimeout)
	return fmt.Errorf("timeout - Failed to meet exit condition")
}

func startWatchers(exitScenario *crd.WatcherSpec, interestingPods chan Condition, stopper chan struct{}) {
	clientSet := k8s.GetK8SConfig()
	factory := informers.NewSharedInformerFactory(clientSet, 0)
	if len(exitScenario.Pods) > 0 {
		go WatchPods(exitScenario.Pods, interestingPods, stopper)
	}
	if len(exitScenario.ConfigMaps) > 0 {
		go WatchBasic(exitScenario.ConfigMaps, interestingPods, stopper, factory.Core().V1().ConfigMaps().Informer())
	}
	if len(exitScenario.Secrets) > 0 {
		go WatchBasic(exitScenario.Secrets, interestingPods, stopper, factory.Core().V1().Secrets().Informer())

	}
	if len(exitScenario.Services) > 0 {
		go WatchBasic(exitScenario.Services, interestingPods, stopper, factory.Core().V1().Services().Informer())
	}
	logger.Info("All conditions checkers started")
}

func checkConditions(goal *ExitScenarioState, clientCrd *crd.CRDClient, in <-chan Condition, stopper chan struct{}) {
	logger.Debug("Started Listener")
	logger.Info(fmt.Sprintf("%#v", goal))
	pendingConditions := len(goal.Conditions)

	for {
		receivedResource := <-in
		fmt.Println("\nInteresting Resource:", receivedResource)
		for key, currentCondition := range goal.Conditions {
			fmt.Println("Key:", key, "Value:", currentCondition)
			if currentCondition.ID == receivedResource.ID &&
				currentCondition.Met == false {
				goal.Conditions[key].Met = true
				logger.Debug("\n Condition  Met:" + fmt.Sprintf("%#v", currentCondition))
				pendingConditions = pendingConditions - 1
				logger.Debug("\n Pending Conditions:" + fmt.Sprintf("%#v", pendingConditions))
				break
			}
		}

		fmt.Println("\n State of Conditions:", goal.Conditions)

		if pendingConditions < 1 {
			logger.Debug("All required objects found, ready to close waiting channels")
			logger.Debug(fmt.Sprintf("%#v", goal.Conditions))
			clientCrd.UpdateStatus(StatusCompleted)
			os.Exit(int(goal.Exit))
		}
	}
}
func loadExitScenarioFromCRD(watcherSpec crd.WatcherSpec) (*crd.WatcherSpec, *ExitScenarioState, error) {

	exitScenarioState, err := processExitScenario(&watcherSpec)
	if err != nil {
		logger.Info(fmt.Sprintf("Error processing Scenario State: %v", err))
		return nil, nil, err
	}
	logger.Info(fmt.Sprintf("Log processing exitScenarioState: %v", exitScenarioState))
	return &watcherSpec, exitScenarioState, nil
}

func loadExitScenario(file string) (*crd.WatcherSpec, *ExitScenarioState, error) {
	exitScenario := &crd.WatcherSpec{}
	logger.Debug("Loading config file:" + file)
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Info(fmt.Sprintf("yamlFile.Get err   #%v ", err))
		return nil, nil, err
	}
	err = yaml.Unmarshal(yamlFile, exitScenario)
	if err != nil {
		logger.Info(fmt.Sprintf("Unmarshal: %v", err))
		return nil, nil, err
	}

	exitScenarioState, err := processExitScenario(exitScenario)
	if err != nil {
		logger.Info(fmt.Sprintf("Error processing Scenario State: %v", err))
		return nil, nil, err
	}
	logger.Info(fmt.Sprintf("Log processing exitScenarioState: %v", exitScenarioState))
	return exitScenario, exitScenarioState, nil
}

func processExitScenario(exitScenario *crd.WatcherSpec) (*ExitScenarioState, error) {
	exitScenarioState := &ExitScenarioState{}
	exitScenarioState.Exit = exitScenario.Exit
	exitScenarioState.Timeout = exitScenario.Timeout
	exitScenarioState.Conditions = []Condition{}

	id := 1
	for k, _ := range exitScenario.Pods {
		exitScenario.Pods[k].ID = id
		exitScenarioState.Conditions = append(exitScenarioState.Conditions, Condition{ID: id, Met: false, Description: fmt.Sprintf("%#v", exitScenario.Pods[k])})
		id++
	}
	for k, _ := range exitScenario.ConfigMaps {
		exitScenario.ConfigMaps[k].ID = id
		exitScenarioState.Conditions = append(exitScenarioState.Conditions, Condition{ID: id, Met: false, Description: fmt.Sprintf("%#v", exitScenario.ConfigMaps[k])})
		id++
	}
	for k, _ := range exitScenario.Secrets {
		exitScenario.Secrets[k].ID = id
		exitScenarioState.Conditions = append(exitScenarioState.Conditions, Condition{ID: id, Met: false, Description: fmt.Sprintf("%#v", exitScenario.Secrets[k])})
		id++
	}
	for k, _ := range exitScenario.Services {
		exitScenario.Services[k].ID = id
		exitScenarioState.Conditions = append(exitScenarioState.Conditions, Condition{ID: id, Met: false, Description: fmt.Sprintf("%#v", exitScenario.Services[k])})
		id++
	}
	for k, _ := range exitScenario.Jobs {
		exitScenario.Jobs[k].ID = id
		exitScenarioState.Conditions = append(exitScenarioState.Conditions, Condition{ID: id, Met: false, Description: fmt.Sprintf("%#v", exitScenario.Jobs[k])})
		id++
	}
	return exitScenarioState, nil
}
