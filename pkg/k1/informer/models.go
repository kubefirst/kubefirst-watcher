package informer

import "github.com/6za/k1-watcher/pkg/k1/crd"

type WatcherConfig struct {
	CrdName      string `yaml:"crdname,omitempty"`
	CrdNamespace string `yaml:"crdnamespace,omitempty"`
	Kind         string `yaml:"kind,omitempty"`
	APIVersion   string `yaml:"apiversion,omitempty"`
	Group        string `yaml:"group,omitempty"`
}

type ExitScenario struct {
	Exit       int32                             `yaml:"exit"`
	Timeout    int32                             `yaml:"timeout"`
	Pods       []PodCondition                    `yaml:"pods"`
	ConfigMaps []crd.BasicConfigurationCondition `yaml:"configmaps"`
	Secrets    []crd.BasicConfigurationCondition `yaml:"secrets"`
	Services   []crd.BasicConfigurationCondition `yaml:"services"`
}

type PodCondition struct {
	ID         int               `yaml:"id"`
	Namespace  string            `yaml:"namespace"`
	Name       string            `yaml:"name"`
	Phase      string            `yaml:"phase"`
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Labels     map[string]string `yaml:"labels"`
}

type ExitScenarioState struct {
	Exit       int32       `yaml:"exit"`
	Timeout    int32       `yaml:"timeout"`
	Conditions []Condition `yaml:"conditions"`
}

type Condition struct {
	Selector    string `yaml:"selector"`
	ID          int    `yaml:"id"`
	Met         bool   `yaml:"met"`
	Description string `yaml:"description"`
}

type BasicConfiguration struct {
	Namespace string            `yaml:"namespace"`
	Name      string            `yaml:"name"`
	Labels    map[string]string `yaml:"labels"`
}

type BasicK8s interface {
	GetNamespace() string
	GetName() string
}
