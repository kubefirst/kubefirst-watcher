package informer_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/kubefirst/kubefirst-watcher/pkg/k1/crd"
	"github.com/kubefirst/kubefirst-watcher/pkg/k1/informer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
func Test_Subset(t *testing.T) {
	funk.Contains([]string{"foo", "bar"}, "bar")
	labelAll := map[string]string{"label1": "value1", "label2": "value2", "label3": "value3"}
	labelSubset := map[string]string{"label2": "value2", "label3": "value3"}

	isSubset := funk.Contains(labelAll, labelSubset)
	//funk.Subset(labelSubset, labelAll)
	if !isSubset {
		t.Errorf("Not isSubset, got: %v ", isSubset)
	}

}
*/
var logger *zap.Logger

func TestMain(m *testing.M) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	myLogger, err := config.Build()
	if err != nil {
		fmt.Printf("Error parsing config options - %s", err)
	}
	logger = myLogger
	code := m.Run()
	os.Exit(code)
}

func Test_ExtractBasicConfigurationMatchTest01(t *testing.T) {
	logger.Debug("Test BasicConfiguration Match equals")
	basicConfig := &crd.BasicConfigurationCondition{
		Namespace: "namespaceObj",
		Name:      "nameObj",
	}

	propertyExpected := informer.ExtractBasicConfigurationMap(basicConfig)
	propertyFound := map[string]string{"name": "nameObj", "namespace": "namespaceObj"}

	match, _ := informer.IsMapPresent(propertyFound, propertyExpected)
	if !match {
		t.Errorf("Not match, got: %v ", match)
	}
}

func Test_ExtractBasicConfigurationMatchTest02(t *testing.T) {
	logger.Debug("Test BasicConfiguration Match namespace vs object of same namespace")
	basicConfig := &crd.BasicConfigurationCondition{
		Namespace: "namespaceObj",
	}

	propertyExpected := informer.ExtractBasicConfigurationMap(basicConfig)
	propertyFound := map[string]string{"name": "nameObj", "namespace": "namespaceObj"}

	match, _ := informer.IsMapPresent(propertyFound, propertyExpected)
	if !match {
		t.Errorf("Not match, got: %v ", match)
	}
}

func Test_ExtractBasicConfigurationMatchTest03(t *testing.T) {
	logger.Debug("Test BasicConfiguration Match off different objects")
	basicConfig := &crd.BasicConfigurationCondition{
		Namespace: "namespaceObj",
		Name:      "nameObj",
	}

	propertyExpected := informer.ExtractBasicConfigurationMap(basicConfig)
	propertyFound := map[string]string{"name": "nameObj02", "namespace": "namespaceObj"}

	match, _ := informer.IsMapPresent(propertyFound, propertyExpected)
	if match {
		t.Errorf("match, got: %v ", match)
	}
}

func Test_ExtractJobTest01(t *testing.T) {
	logger.Debug("Test Job Match equals")
	jobConditionConfig := &crd.JobCondition{
		Namespace: "namespaceObj",
		Name:      "nameObj",
		Succeeded: 1,
	}
	jobFound := &batchv1.Job{
		ObjectMeta: v1.ObjectMeta{Name: "nameObj", Namespace: "namespaceObj"},
		Status:     batchv1.JobStatus{Succeeded: 1},
	}

	propertyExpected := informer.ExtractJobConditionMap(jobConditionConfig)
	propertyFound := informer.ExtractJobMap(jobFound)

	match, _ := informer.IsMapPresent(propertyFound, propertyExpected)
	if !match {
		t.Errorf("Not match, got: %v ", match)
	}
}

func Test_ExtractJobTest02(t *testing.T) {
	logger.Debug("Test Job Match namespace vs object of same namespace")
	jobConditionConfig := &crd.JobCondition{
		Namespace: "namespaceObj",
	}
	jobFound := &batchv1.Job{
		ObjectMeta: v1.ObjectMeta{Name: "nameObj", Namespace: "namespaceObj"},
	}

	propertyExpected := informer.ExtractJobConditionMap(jobConditionConfig)
	propertyFound := informer.ExtractJobMap(jobFound)

	match, _ := informer.IsMapPresent(propertyFound, propertyExpected)
	if !match {
		t.Errorf("Not match, got: %v ", match)
	}
}

func Test_ExtractJobTest03(t *testing.T) {
	logger.Debug("Test Job Match off different objects")
	jobConditionConfig := &crd.JobCondition{
		Namespace: "namespaceObj",
		Name:      "nameObj",
	}
	jobFound := &batchv1.Job{
		ObjectMeta: v1.ObjectMeta{Name: "nameObj02", Namespace: "namespaceObj"},
	}

	propertyExpected := informer.ExtractJobConditionMap(jobConditionConfig)
	propertyFound := informer.ExtractJobMap(jobFound)

	match, _ := informer.IsMapPresent(propertyFound, propertyExpected)
	if match {
		t.Errorf("match, got: %v ", match)
	}
}

func Test_ExtractPodTest01(t *testing.T) {
	logger.Debug("Test Pod Match equals")
	podConditionConfig := &crd.PodCondition{
		Namespace: "namespaceObj",
		Name:      "nameObj",
		Phase:     string(corev1.PodRunning),
	}
	podFound := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{Name: "nameObj", Namespace: "namespaceObj"},
		Status:     corev1.PodStatus{Phase: corev1.PodRunning},
	}

	propertyExpected := informer.ExtractPodConditionMap(podConditionConfig)
	propertyFound := informer.ExtractPodMap(podFound)

	match, _ := informer.IsMapPresent(propertyFound, propertyExpected)
	if !match {
		t.Errorf("Not match, got: %v ", match)
	}
}

func Test_ExtractPodTest02(t *testing.T) {
	logger.Debug("Test Pod Match namespace vs object of same namespace")
	podConditionConfig := &crd.PodCondition{
		Namespace: "namespaceObj",
	}
	podFound := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{Name: "nameObj", Namespace: "namespaceObj"},
	}

	propertyExpected := informer.ExtractPodConditionMap(podConditionConfig)
	propertyFound := informer.ExtractPodMap(podFound)

	match, _ := informer.IsMapPresent(propertyFound, propertyExpected)
	if !match {
		t.Errorf("Not match, got: %v ", match)
	}
}

func Test_ExtractPodTest03(t *testing.T) {
	logger.Debug("Test Pod Match off different objects")
	podConditionConfig := &crd.PodCondition{
		Namespace: "namespaceObj",
		Name:      "nameObj",
	}
	podFound := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{Name: "nameObj02", Namespace: "namespaceObj"},
	}

	propertyExpected := informer.ExtractPodConditionMap(podConditionConfig)
	propertyFound := informer.ExtractPodMap(podFound)

	match, _ := informer.IsMapPresent(propertyFound, propertyExpected)
	if match {
		t.Errorf("match, got: %v ", match)
	}
}

func Test_ExtractPodTest04(t *testing.T) {
	logger.Debug("Test Pod Match off different objects different phase")
	podConditionConfig := &crd.PodCondition{
		Namespace: "namespaceObj",
		Name:      "nameObj",
		Phase:     string(corev1.PodRunning),
	}
	podFound := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{Name: "nameObj", Namespace: "namespaceObj"},
		Status:     corev1.PodStatus{Phase: corev1.PodPending},
	}

	propertyExpected := informer.ExtractPodConditionMap(podConditionConfig)
	propertyFound := informer.ExtractPodMap(podFound)

	match, _ := informer.IsMapPresent(propertyFound, propertyExpected)
	if match {
		t.Errorf(" match, got: %v ", match)
	}
}
