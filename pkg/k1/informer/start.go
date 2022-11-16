package informer

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// StartWatcher - starts watcher tooling
var logger *zap.Logger

func StartWatcher(configFile string, loggerIn *zap.Logger) error {
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

	go checkConditions(exitScenarioState, interestingPods, stopper)
	//Start Watchers
	if len(exitScenario.Pods) > 0 {
		go WatchPods(exitScenario.Pods, interestingPods, stopper)
	}
	/*
		for k, _ := range exitScenarioState.Conditions {
			conditionGVK := schema.FromAPIVersionAndKind(exitScenario.Conditions[k].APIVersion, exitScenario.Conditions[k].Kind)
			switch conditionGVK {
			case podGVK:
				//go WatchPods(exitScenario.Conditions[k], interestingPods, stopper)
				//		case "linux":
				//			fmt.Println("Linux.")
			default:
				logger.Error(fmt.Sprintf("Error %#v", conditionGVK))
				return fmt.Errorf("err - unkwon checker for GroupVersionKind")
			}

		}
	*/
	logger.Info("All conditions checkers started")
	//Check Current State - to catch events pre-informers are started
	time.Sleep(time.Duration(exitScenario.Timeout) * time.Second)
	return fmt.Errorf("timeout - Failed to meet exit condition")
}

func checkConditions(goal *ExitScenarioState, in <-chan Condition, stopper chan struct{}) {
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
			os.Exit(goal.Exit)
		}
	}
}

func loadExitScenario(file string) (*ExitScenario, *ExitScenarioState, error) {
	exitScenario := &ExitScenario{}
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

func processExitScenario(exitScenario *ExitScenario) (*ExitScenarioState, error) {
	exitScenarioState := &ExitScenarioState{}
	exitScenarioState.Exit = exitScenario.Exit
	exitScenarioState.Timeout = exitScenario.Timeout
	exitScenarioState.Conditions = []Condition{}

	id := 1
	for k, _ := range exitScenario.Pods {
		exitScenario.Pods[k].ID = id
		exitScenarioState.Conditions = append(exitScenarioState.Conditions, Condition{ID: id, Met: false})
		id++
	}
	for k, _ := range exitScenario.ConfigMaps {
		exitScenario.ConfigMaps[k].ID = id
		exitScenarioState.Conditions = append(exitScenarioState.Conditions, Condition{ID: id, Met: false})
		id++
	}
	for k, _ := range exitScenario.Secrets {
		exitScenario.Secrets[k].ID = id
		exitScenarioState.Conditions = append(exitScenarioState.Conditions, Condition{ID: id, Met: false})
		id++
	}
	for k, _ := range exitScenario.Services {
		exitScenario.Services[k].ID = id
		exitScenarioState.Conditions = append(exitScenarioState.Conditions, Condition{ID: id, Met: false})
		id++
	}
	return exitScenarioState, nil
}
