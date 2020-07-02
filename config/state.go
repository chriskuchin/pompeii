package config

type (
	ServiceState struct {
		Count   int64
		TaskDef string
	}

	TaskState struct {
		Running bool
		Failed  bool
	}

	// ServiceWeights the listener/target group weights
	ServiceWeights struct {
		Canary  int64
		Primary int64
	}
)
