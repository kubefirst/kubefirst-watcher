package informer

type WatcherConfig struct {
	CrdName      string `yaml:"crdname,omitempty"`
	CrdNamespace string `yaml:"crdnamespace,omitempty"`
	Kind         string `yaml:"kind,omitempty"`
	APIVersion   string `yaml:"apiversion,omitempty"`
	Group        string `yaml:"group,omitempty"`
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

type BasicK8s interface {
	GetNamespace() string
	GetName() string
}
