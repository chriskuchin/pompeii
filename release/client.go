package release

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/chriskuchin/pompeii/config"
)

type (
	// Client test
	Client struct {
		Config   *config.Config
		elbv2Svc map[string]*elbv2.ELBV2
		ecsSvc   map[string]*ecs.ECS
	}
)

// NewClient returns an awsClient
func NewClient(config *config.Config) *Client {
	client := &Client{
		Config:   config,
		elbv2Svc: map[string]*elbv2.ELBV2{},
		ecsSvc:   map[string]*ecs.ECS{},
	}

	for service := range config.Services {
		session := session.New(&aws.Config{
			Region: aws.String(config.GetRegion(service)),
		})

		client.elbv2Svc[service] = elbv2.New(session)
		client.ecsSvc[service] = ecs.New(session)
	}

	return client
}
