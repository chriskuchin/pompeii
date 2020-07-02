package config

type (
	ActionType string

	Action struct {
		Type      ActionType `yaml:"action"`
		Target    string     `yaml:"target"`
		Ratio     int64      `yaml:"ratio"`
		Validator string     `yaml:"validator"`
		Count     int64      `yaml:"count"`
		Task      string     `yaml:"task"`
		Command   []string   `yaml:"command"`
	}

	Workflow struct {
		Config  *Config
		Service string

		Default *ServiceState
		Steps   []*Action
	}
)

const (
	TrafficShift ActionType = "shift"
	ValidatePool ActionType = "validate"
	UpdatePool   ActionType = "update"
)
