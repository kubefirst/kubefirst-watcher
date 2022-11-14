package informer

type ExitScenario struct {
	Exit       int         `yaml:"exit"`
	Conditions []Condition `yaml:"conditions"`
}

type Condition struct {
	Namespace  string `yaml:"namespace"`
	Name       string `yaml:"name"`
	Phase      string `yaml:"phase"`
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Met        bool   `yaml:"kind"`
}
