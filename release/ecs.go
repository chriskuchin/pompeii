package release

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/chriskuchin/pompeii/config"
	"github.com/prometheus/common/log"
)

func (c *Client) GetCurrentServiceState(service, pool string) *config.ServiceState {
	info := c.GetCurrentServiceInfo(service, pool)

	return &config.ServiceState{
		TaskDef: *info.TaskDefinition,
		Count:   *info.DesiredCount,
	}
}

func (c *Client) GetCurrentServiceInfo(service, pool string) *ecs.Service {
	svc := c.ecsSvc[service]
	input := &ecs.DescribeServicesInput{
		Services: []*string{
			aws.String(c.Config.GetECSService(service, pool)),
		},
		Cluster: aws.String(c.Config.Services[service].ClusterARN),
	}

	result, err := svc.DescribeServices(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			case ecs.ErrCodeClusterNotFoundException:
				fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	log.Debug(result)

	return result.Services[0]
}

func (c *Client) UpdateService(service, pool string, state *config.ServiceState) {
	svc := c.ecsSvc[service]
	input := &ecs.UpdateServiceInput{
		Service:        aws.String(c.Config.Services[service].Canary.Service),
		TaskDefinition: aws.String(state.TaskDef),
		Cluster:        aws.String(c.Config.Services[service].ClusterARN),
		DesiredCount:   aws.Int64(state.Count),
	}

	result, err := svc.UpdateService(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			case ecs.ErrCodeClusterNotFoundException:
				fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
			case ecs.ErrCodeServiceNotFoundException:
				fmt.Println(ecs.ErrCodeServiceNotFoundException, aerr.Error())
			case ecs.ErrCodeServiceNotActiveException:
				fmt.Println(ecs.ErrCodeServiceNotActiveException, aerr.Error())
			case ecs.ErrCodePlatformUnknownException:
				fmt.Println(ecs.ErrCodePlatformUnknownException, aerr.Error())
			case ecs.ErrCodePlatformTaskDefinitionIncompatibilityException:
				fmt.Println(ecs.ErrCodePlatformTaskDefinitionIncompatibilityException, aerr.Error())
			case ecs.ErrCodeAccessDeniedException:
				fmt.Println(ecs.ErrCodeAccessDeniedException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	log.Debug(result)
}

func (c *Client) RunTask(service, task, container string, command []string) (string, error) {
	log.Info(service, task)
	input := &ecs.RunTaskInput{
		Cluster:        aws.String(c.Config.Services[service].ClusterARN),
		TaskDefinition: aws.String(task),
	}

	if len(command) > 0 {
		commandOverride := []*string{}
		for _, piece := range command {
			commandOverride = append(commandOverride, aws.String(piece))
		}

		input.Overrides = &ecs.TaskOverride{
			ContainerOverrides: []*ecs.ContainerOverride{
				{
					Name:    aws.String(container),
					Command: commandOverride,
				},
			},
		}
	}

	return c.runTask(service, input)
}

func (c *Client) runTask(service string, taskInput *ecs.RunTaskInput) (string, error) {
	log.Info(taskInput)
	svc := c.ecsSvc[service]

	result, err := svc.RunTask(taskInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			case ecs.ErrCodeClusterNotFoundException:
				fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
			case ecs.ErrCodeUnsupportedFeatureException:
				fmt.Println(ecs.ErrCodeUnsupportedFeatureException, aerr.Error())
			case ecs.ErrCodePlatformUnknownException:
				fmt.Println(ecs.ErrCodePlatformUnknownException, aerr.Error())
			case ecs.ErrCodePlatformTaskDefinitionIncompatibilityException:
				fmt.Println(ecs.ErrCodePlatformTaskDefinitionIncompatibilityException, aerr.Error())
			case ecs.ErrCodeAccessDeniedException:
				fmt.Println(ecs.ErrCodeAccessDeniedException, aerr.Error())
			case ecs.ErrCodeBlockedException:
				fmt.Println(ecs.ErrCodeBlockedException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "", err
	}

	log.Debug(result)

	return *result.Tasks[0].TaskArn, nil
}

func (c *Client) DescribeTask(service, taskARN string) (*config.TaskState, error) {
	svc := c.ecsSvc[service]
	input := &ecs.DescribeTasksInput{
		Cluster: aws.String(c.Config.Services[service].ClusterARN),
		Tasks: []*string{
			aws.String(taskARN),
		},
	}

	result, err := svc.DescribeTasks(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			case ecs.ErrCodeClusterNotFoundException:
				fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil, err
	}

	log.Debug(result)

	return c.calcuateTaskState(result), nil
}

func (c *Client) calcuateTaskState(taskInfo *ecs.DescribeTasksOutput) *config.TaskState {
	state := &config.TaskState{}
	if len(taskInfo.Failures) > 0 {
		state.Failed = true

	}

	// We currently only handle single tasks
	if len(taskInfo.Tasks) > 1 {
		return nil
	}

	if *taskInfo.Tasks[0].LastStatus == "PENDING" || *taskInfo.Tasks[0].LastStatus == "RUNNING" {
		state.Running = true
	} else if *taskInfo.Tasks[0].LastStatus == "STOPPED" {
		state.Running = false

		if len(taskInfo.Tasks[0].Containers) == 1 {
			log.Info(*&taskInfo.Tasks[0].Containers)
			if *taskInfo.Tasks[0].Containers[0].ExitCode == 0 {
				state.Failed = false
			} else {
				state.Failed = true
			}
		}
	}
	log.Info(state)
	return state
}

func (c *Client) StartAndMonitorTask(service, task, container string, command []string) bool {
	taskARN, _ := c.RunTask(service, task, container, command)
	return c.monitorTaskRun(service, taskARN)
}

func (c *Client) monitorTaskRun(service, taskARN string) bool {
	state, _ := c.DescribeTask(service, taskARN)

	for state != nil && state.Running {
		time.Sleep(1 * time.Minute)

		state, _ = c.DescribeTask(service, taskARN)
	}

	log.Infof("%+v", state)

	return state != nil && state.Failed
}

func (c *Client) RollbackDeployment(service, pool string) {
	serviceDef := c.GetCurrentServiceInfo(service, pool)

	previousTaskDef := &config.ServiceState{}
	newTaskDef := &config.ServiceState{}
	for _, deployment := range serviceDef.Deployments {
		if strings.ToLower(*deployment.Status) == "primary" {
			// in progress deployment
			newTaskDef.TaskDef = *deployment.TaskDefinition
			newTaskDef.Count = *deployment.DesiredCount
		} else if strings.ToLower(*deployment.Status) == "active" {
			// previous deployment
			previousTaskDef.TaskDef = *deployment.TaskDefinition
			previousTaskDef.Count = *deployment.DesiredCount
		}
	}

	if previousTaskDef.TaskDef == "" || newTaskDef.TaskDef == "" {
		log.Errorf("Failed to locate the previous or new version. previousTaskDef: %+v newTaskDef: %+v", previousTaskDef, newTaskDef)
		return
	}

	log.Debugf("Rolling back from %+v to %+v", newTaskDef, previousTaskDef)

	c.UpdateService(service, pool, previousTaskDef)
}

func getActiveDeployment(info *ecs.Service) *ecs.Deployment {
	for _, deployment := range info.Deployments {
		if strings.ToLower(*deployment.Status) == "primary" {
			return deployment
		}
	}

	return nil
}

func (c *Client) MonitorServiceDeployment(service, pool string) bool {
	for {
		info := c.GetCurrentServiceInfo(service, pool)
		log.Debugf("%#v\n", info.Deployments)

		if len(info.Deployments) == 1 {
			return true
		}

		deployment := getActiveDeployment(info)
		if deployment == nil {
			log.Errorf("Failed to locate the Primary deployment: %#v", info.Deployments)
			return false
		}

		if time.Now().Sub(*deployment.CreatedAt) > c.Config.Services[service].Timeout {
			log.Errorf("Deployment Timed out: %v", c.Config.Services[service].Timeout)
			return false
		}

		log.Info("Waiting for deployment to complete, sleeping")
		time.Sleep(30 * time.Second)
	}
}

// Deploy deploys the given service and waits for the deployment to complete
func (c *Client) Deploy(service, pool string, state *config.ServiceState) bool {
	rollbackState := c.GetCurrentServiceState(service, pool)

	c.UpdateService(service, pool, state)

	result := c.MonitorServiceDeployment(service, pool)

	if result == false {
		c.UpdateService(service, pool, rollbackState)
		return false
	}

	return true
}
