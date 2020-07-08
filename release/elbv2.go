package release

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/chriskuchin/pompeii/config"
	"github.com/prometheus/common/log"
)

// GetCurrentWeights returns the current weights for the primary and canary groups
func (c *Client) GetCurrentWeights(service string) *config.ServiceWeights {
	rule := c.getCurrentRule(service)
	weights := &config.ServiceWeights{}
	for _, action := range rule.Actions {
		if *action.Type == "forward" {
			for _, tg := range action.ForwardConfig.TargetGroups {
				if *tg.TargetGroupArn == c.Config.GetCanaryTargetGroupARN(service) {
					weights.Canary = *tg.Weight
				} else {
					weights.Primary = *tg.Weight
				}
			}
		}
	}

	return weights
}

func (c *Client) getCurrentRule(service string) *elbv2.Rule {
	svc := c.elbv2Svc[service]
	input := &elbv2.DescribeRulesInput{
		RuleArns: []*string{
			aws.String(c.Config.GetListenerRuleARN(service)),
		},
	}

	result, err := svc.DescribeRules(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elbv2.ErrCodeListenerNotFoundException:
				fmt.Println(elbv2.ErrCodeListenerNotFoundException, aerr.Error())
			case elbv2.ErrCodeRuleNotFoundException:
				fmt.Println(elbv2.ErrCodeRuleNotFoundException, aerr.Error())
			case elbv2.ErrCodeUnsupportedProtocolException:
				fmt.Println(elbv2.ErrCodeUnsupportedProtocolException, aerr.Error())
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

	return result.Rules[0]
}

// UpdateWeights updates the given services weight settings
func (c *Client) UpdateWeights(service string, weights *config.ServiceWeights) error {
	svc := c.elbv2Svc[service]
	rule := c.getCurrentRule(service)

	for _, action := range rule.Actions {
		if *action.Type == "forward" {
			for _, target := range action.ForwardConfig.TargetGroups {
				if *target.TargetGroupArn == c.Config.GetCanaryTargetGroupARN(service) {
					target.SetWeight(weights.Canary)
				} else {
					target.SetWeight(weights.Primary)
				}
			}
		}
	}

	log.Debug(rule)

	input := &elbv2.ModifyRuleInput{
		Actions: rule.Actions,
		RuleArn: aws.String(c.Config.GetListenerRuleARN(service)),
	}

	log.Info("Actions: ", rule.Actions)

	result, err := svc.ModifyRule(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elbv2.ErrCodeTargetGroupAssociationLimitException:
				fmt.Println(elbv2.ErrCodeTargetGroupAssociationLimitException, aerr.Error())
			case elbv2.ErrCodeIncompatibleProtocolsException:
				fmt.Println(elbv2.ErrCodeIncompatibleProtocolsException, aerr.Error())
			case elbv2.ErrCodeRuleNotFoundException:
				fmt.Println(elbv2.ErrCodeRuleNotFoundException, aerr.Error())
			case elbv2.ErrCodeOperationNotPermittedException:
				fmt.Println(elbv2.ErrCodeOperationNotPermittedException, aerr.Error())
			case elbv2.ErrCodeTooManyRegistrationsForTargetIdException:
				fmt.Println(elbv2.ErrCodeTooManyRegistrationsForTargetIdException, aerr.Error())
			case elbv2.ErrCodeTooManyTargetsException:
				fmt.Println(elbv2.ErrCodeTooManyTargetsException, aerr.Error())
			case elbv2.ErrCodeTargetGroupNotFoundException:
				fmt.Println(elbv2.ErrCodeTargetGroupNotFoundException, aerr.Error())
			case elbv2.ErrCodeUnsupportedProtocolException:
				fmt.Println(elbv2.ErrCodeUnsupportedProtocolException, aerr.Error())
			case elbv2.ErrCodeTooManyActionsException:
				fmt.Println(elbv2.ErrCodeTooManyActionsException, aerr.Error())
			case elbv2.ErrCodeInvalidLoadBalancerActionException:
				fmt.Println(elbv2.ErrCodeInvalidLoadBalancerActionException, aerr.Error())
			case elbv2.ErrCodeTooManyUniqueTargetGroupsPerLoadBalancerException:
				fmt.Println(elbv2.ErrCodeTooManyUniqueTargetGroupsPerLoadBalancerException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return err
	}

	log.Debug(result)
	return nil
}
