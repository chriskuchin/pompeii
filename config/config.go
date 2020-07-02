package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/prometheus/common/log"
	"gopkg.in/yaml.v2"
)

type (
	// Config config object
	Config struct {
		ClusterARN string                    `yaml:"cluster-arn"`
		Region     string                    `yaml:"region"`
		Services   map[string]*ServiceConfig `yaml:"services"`
		Workflows  WorkflowConfig            `yaml:"workflows"`
	}

	// ServiceConfig test
	ServiceConfig struct {
		ClusterARN              string        `yaml:"cluster-arn"`
		ListenerARN             string        `yaml:"listener-rule-arn"`
		Region                  string        `yaml:"region"`
		Timeout                 time.Duration `yaml:"deploy-timeout"`
		ValidationTask          string        `yaml:"valdation-task"`
		ValidationTaskContainer string        `yaml:"validation-task-container"`

		Canary  *PoolConfig `yaml:"canary"`
		Primary *PoolConfig `yaml:"primary"`
	}

	// PoolConfig test
	PoolConfig struct {
		TargetGroupARN string `yaml:"tg-arn"`
		Service        string `yaml:"ecs-service"`
	}

	WorkflowConfig map[string][]*Action
)

// NewConfigFromFile loads a yaml file intoa  config struct
func NewConfigFromFile(path string) *Config {
	configFile, err := os.Open(path)
	if err != nil {
		return nil
	}

	rawYaml, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil
	}

	config := &Config{}
	err = yaml.Unmarshal(rawYaml, config)
	if err != nil {
		return nil
	}
	return config
}

func NewConfigFromS3(key, bucket string) *Config {
	svc := s3.New(session.New())
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				fmt.Println(s3.ErrCodeNoSuchKey, aerr.Error())
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

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil
	}

	log.Debugf("BODY####\n%s\n", body)
	config := &Config{}
	yaml.Unmarshal(body, config)

	return config
}

// GetRegion Looks up the region for a service in the config falling back to the root region if undefined
func (c *Config) GetRegion(service string) string {
	if c.Services[service] != nil && c.Services[service].Region != "" {
		return c.Services[service].Region
	}

	return c.Region
}

// GetListenerRuleARN returns the services listener rule arn
func (c *Config) GetListenerRuleARN(service string) string {
	return c.Services[service].ListenerARN
}

// GetCanaryTargetGroupARN returns the canary services target group arn
func (c *Config) GetCanaryTargetGroupARN(service string) string {
	return c.Services[service].Canary.TargetGroupARN
}

// GetPrimaryTargetGroupARN returns the priamry groups target group arn
func (c *Config) GetPrimaryTargetGroupARN(service string) string {
	return c.Services[service].Primary.TargetGroupARN
}

// IsCanaryTargetGroup tests whether a provided arn is the canary pools target group
func (c *Config) IsCanaryTargetGroup(service, arn string) bool {
	return c.GetCanaryTargetGroupARN(service) == arn
}

// IsPrimaryTargetGroup tests if a given arn is the target group arn for a defined service
func (c *Config) IsPrimaryTargetGroup(service, arn string) bool {
	return c.GetPrimaryTargetGroupARN(service) == arn
}

func (c *Config) GetECSService(service, pool string) string {
	if strings.ToLower(pool) == "canary" {
		return c.Services[service].Canary.Service
	} else if strings.ToLower(pool) == "primary" {
		return c.Services[service].Primary.Service
	}

	return ""
}

func (c *Config) GetServiceValidationTask(service string) string {
	return c.Services[service].ValidationTask
}
