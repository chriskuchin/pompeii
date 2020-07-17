package workflow

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/chriskuchin/pompeii/config"
	"github.com/chriskuchin/pompeii/release"
	"github.com/prometheus/common/log"
)

type (
	Processor struct {
		workflow *config.Workflow
		client   *release.Client

		checkpoint *Checkpoint
	}

	Checkpoint struct {
		Canary  *config.ServiceState
		Primary *config.ServiceState
		Weights *config.ServiceWeights
	}
)

func ProcessWorkflow(workflow *config.Workflow) error {
	processor := &Processor{
		workflow:   workflow,
		client:     release.NewClient(workflow.Config),
		checkpoint: &Checkpoint{},
	}

	processor.getInitialCheckpoint()

	for _, action := range workflow.Steps {
		switch action.Type {
		case config.UpdatePool:
			log.Infof("Update %s pool", action.Target)
			proceed := processor.handleUpdateAction(action)

			if !proceed {
				log.Error("Pool Update Failed!!")
				processor.rollbackToLatestCheckpoint()
				return fmt.Errorf("Failed to update pool: rolled back")
			}

		case config.TrafficShift:
			log.Infof("Shift traffic to pool: %s weight: %d", action.Target, action.Ratio)
			proceed := processor.handleShiftAction(action)

			if !proceed {
				log.Error("Traffic Shift Failed!!")
				processor.rollbackToLatestCheckpoint()
				return fmt.Errorf("Failed to shift traffic: rolled back")
			}

		case config.ValidatePool:
			log.Infof("Validate: %+v", action)
			proceed := processor.handleValidationAction(action)

			log.Info(proceed)
			if !proceed {
				log.Error("Validation Failed!!")
				processor.rollbackToLatestCheckpoint()
				return fmt.Errorf("Failed to validate pools: rolled back")
			}

		default:
			log.Errorf("Undefined ActionType: %s", action.Type)
		}
	}

	return nil
}

func (p *Processor) getInitialCheckpoint() {
	p.checkpoint.Weights = p.client.GetCurrentWeights(p.workflow.Service)

	p.checkpoint.Canary = p.client.GetCurrentServiceState(p.workflow.Service, "canary")
	p.checkpoint.Primary = p.client.GetCurrentServiceState(p.workflow.Service, "primary")
}

func (p *Processor) rollbackToLatestCheckpoint() {
	p.client.UpdateWeights(p.workflow.Service, p.checkpoint.Weights)

	p.client.Deploy(p.workflow.Service, "canary", p.checkpoint.Canary)
	p.client.Deploy(p.workflow.Service, "primary", p.checkpoint.Primary)

}

func (p *Processor) getUpdateActionServiceState(action *config.Action) *config.ServiceState {
	result := p.workflow.Default

	if action.Count != 0 {
		result.Count = action.Count
	}

	if action.Task != "" {
		result.TaskDef = action.Task
	}

	return result

}

func (p *Processor) handleUpdateAction(action *config.Action) bool {
	return p.client.Deploy(p.workflow.Service, action.Target, p.getUpdateActionServiceState(action))
}

func (p *Processor) handleShiftAction(action *config.Action) bool {
	weights := &config.ServiceWeights{}

	if action.Target == "canary" {
		weights.Canary = action.Ratio
		weights.Primary = 100 - action.Ratio
	} else {
		weights.Primary = action.Ratio
		weights.Canary = 100 - action.Ratio
	}

	err := p.client.UpdateWeights(p.workflow.Service, weights)

	return err == nil
}

func (p *Processor) handleValidationAction(action *config.Action) bool {

	switch action.Target {
	case "prompt":
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Does the current system state pass validation (y/n)? ")
		answer, err := reader.ReadString('\n')
		if err != nil {
			log.Error(err)
		}

		if strings.ToLower(strings.Trim(answer, "\n")) == "y" {
			log.Info("Continue")
			return true
		}

		log.Info("Cancel")
		return false

	case "task":
		return !p.client.StartAndMonitorTask(p.workflow.Service, p.client.Config.GetServiceValidationTask(p.workflow.Service), p.client.Config.Services[p.workflow.Service].ValidationTaskContainer, action.Command)
	default:
		return false
	}
}
