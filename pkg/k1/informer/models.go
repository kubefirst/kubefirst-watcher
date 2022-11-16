package informer

type ExitScenario struct {
	Exit       int                           `yaml:"exit"`
	Timeout    int                           `yaml:"timeout"`
	Pods       []PodCondition                `yaml:"pods"`
	ConfigMaps []BasicConfigurationCondition `yaml:"configmaps"`
	Secrets    []BasicConfigurationCondition `yaml:"secrets"`
	Services   []BasicConfigurationCondition `yaml:"services"`
}

type PodCondition struct {
	ID         int    `yaml:"id"`
	Namespace  string `yaml:"namespace"`
	Name       string `yaml:"name"`
	Phase      string `yaml:"phase"`
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
}

type BasicConfigurationCondition struct {
	ID         int    `yaml:"id"`
	Namespace  string `yaml:"namespace"`
	Name       string `yaml:"name"`
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
}

type ExitScenarioState struct {
	Exit       int         `yaml:"exit"`
	Timeout    int         `yaml:"timeout"`
	Conditions []Condition `yaml:"conditions"`
}

type Condition struct {
	Selector string `yaml:"selector"`
	ID       int    `yaml:"id"`
	Met      bool   `yaml:"met"`
}

type BasicConfiguration struct {
	Namespace string `yaml:"namespace"`
	Name      string `yaml:"name"`
}
