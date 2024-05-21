package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

// func CreateEventBus(ctx context.Context, client *eventbridge.Client) {
// 	input := &eventbridge.CreateEventBusInput{
// 		Name: aws.String("eventbus-2"),
// 	}

// 	result, err := client.CreateEventBus(ctx, input)
// 	if err != nil {
// 		fmt.Println("Error creating event bus:", err)
// 		return
// 	}

// 	fmt.Println("Event bus ARN:", *result.EventBusArn)
// }

func CreateRule(ctx context.Context, client *eventbridge.Client) {
	eventPattern := `{
        "source": ["aws.ecr"],
        "detail-type": ["ECR Image Action"]
    }`

	input := &eventbridge.PutRuleInput{
		Name:         aws.String("Rule-ECRPushEvent3"),
		EventPattern: aws.String(eventPattern),
		EventBusName: aws.String("default"),
	}

	result, err := client.PutRule(ctx, input)
	if err != nil {
		fmt.Println("Error creating rule:", err)
		return
	}

	fmt.Println("Rule ARN:", *result.RuleArn)
}

func AddTarget(ctx context.Context, client *eventbridge.Client) {
	target := types.Target{
		Arn: aws.String("arn:aws:lambda:us-east-1:975050154225:function:saipujitha"),
		Id:  aws.String("Lambda-1"),
	}

	input := &eventbridge.PutTargetsInput{
		Rule:         aws.String("Rule-ECRPushEvent3"),
		Targets:      []types.Target{target},
		EventBusName: aws.String("default"),
	}

	_, err := client.PutTargets(ctx, input)
	if err != nil {
		fmt.Println("Error adding target to rule:", err)
		return
	}

	fmt.Println("Target added to rule successfully")
}

func main() {
	fmt.Println("AWS EventBridge Setup")

	// Load the AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create an EventBridge client
	client := eventbridge.NewFromConfig(cfg)

	// Create a context
	ctx := context.TODO()

	// CreateEventBus(ctx, client)
	CreateRule(ctx, client)
	AddTarget(ctx, client)
}
