package main

import (
	"fmt"
	"os"

	"github.com/chriskuchin/pompeii/config"
	"github.com/chriskuchin/pompeii/workflow"
	"github.com/urfave/cli/v2"
)

type ()

var (
	validPools = []string{"canary", "primary"}
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config-file",
				Value: "config.yml",
			},
			&cli.StringFlag{
				Name:     "service",
				Required: true,
			},
			&cli.StringFlag{
				Name: "s3-bucket",
			},
		},
		Commands: []*cli.Command{
			{
				Name: "deploy",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "workflow",
						Value: "default",
					},
					&cli.StringFlag{
						Name:     "task-def",
						Required: true,
					},
					&cli.Int64Flag{
						Name:  "count",
						Value: 2,
					},
				},
				Action: func(c *cli.Context) error {
					settings := initClient(c)
					fmt.Println(settings)
					return workflow.ProcessWorkflow(&config.Workflow{
						Config:  settings,
						Service: c.String("service"),
						Steps:   settings.Workflows[c.String("workflow")],
						Default: &config.ServiceState{
							TaskDef: c.String("task-def"),
							Count:   c.Int64("count"),
						},
					})
				},
			},
		},
	}

	app.Run(os.Args)
}

func initClient(c *cli.Context) *config.Config {
	filePath := c.String("config-file")
	s3Bucket := c.String("s3-bucket")

	clientConfig := &config.Config{}
	if s3Bucket != "" {
		clientConfig = config.NewConfigFromS3(filePath, s3Bucket)
	} else {
		clientConfig = config.NewConfigFromFile(filePath)

	}

	return clientConfig
}

// 1. shift all traffic to primary pool
// 2. deploy new code to canary pool
// 3. test canary pool
// 4. shift a small fraction of traffic to canary pool
// 5. watch metrics
// 6. shift 50%
// 7. shift 100%
// 8. deploy to primary pool
// 9. shift to 50/50
// 10. done

/// shiftTraffic adjusts the relative weights of the canary and primary pool
// inputs canaryWeight/ratio primary weight/ratio
func shiftTraffic() {
	// read current weights
	// update current rules to new ratio
	// done
}

// validate traffic... can be a query against prom or running an ecs task and checking output
func validateTraffic() {
	// ecs task
	// launch task
	// poll task
	// collect exit code
	// return
}

func updatePool() {
	// update ecs service
}
