package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	REGION       string = "us-east-1"
	ENDPOINT_ENV string = "AWS_ENDPOINT"
	QUEUE_NAME          = "notificacao-sqs"
)

func TestAll(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "localstack/localstack",
		ExposedPorts: []string{"4566/tcp"},
		WaitingFor:   wait.ForLog("Ready."),
		Env: map[string]string{
			"SERVICES":              "sqs",
			"DEFAULT_REGION":        REGION,
			"AWS_ACCESS_KEY_ID":     "test",
			"AWS_SECRET_ACCESS_KEY": "test",
		},
	}
	localstackContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer localstackContainer.Terminate(ctx)
	sqsClient := prepareSqsClient()
	returnQueue, err := sqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName: aws.String(QUEUE_NAME),
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(*returnQueue.QueueUrl)
	t.Setenv(ENDPOINT_ENV, *returnQueue.QueueUrl)
	t.Run("Sucess", LocalContainerSimplified)
	t.Run("Error", LocalContainerErrorSimplified)
}

func prepareSqsClient() sqs.Client {
	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(REGION))
	if err != nil {
		log.Println("Não foi possível carregar a configuração padrão. Você configurou suas credenciais AWS?")
		log.Fatal(err)
	}
	sqsClient := sqs.NewFromConfig(sdkConfig)
	sqsClient.Options().EndpointResolverV2.ResolveEndpoint(ctx, sqs.EndpointParameters{
		Region:   aws.String(REGION),
		Endpoint: aws.String("http://localhost:4566"),
	})
	return *sqsClient
}

func LocalContainerSimplified(t *testing.T) {
	ctx := context.Background()
	sqsClient := prepareSqsClient()
	_, err := sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String("Ola"),
		QueueUrl:    aws.String(os.Getenv(ENDPOINT_ENV)),
	})
	assert.NoError(t, err)
}

func LocalContainerErrorSimplified(t *testing.T) {
	ctx := context.Background()
	sqsClient := prepareSqsClient()
	_, err := sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String("Ola"),
		QueueUrl:    aws.String("http:"),
	})
	assert.Error(t, err)
}

func TestLocalContainer(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "localstack/localstack",
		ExposedPorts: []string{"4566/tcp"},
		WaitingFor:   wait.ForLog("Ready."),
		Env: map[string]string{
			"SERVICES":              "sqs",
			"DEFAULT_REGION":        REGION,
			"AWS_ACCESS_KEY_ID":     "test",
			"AWS_SECRET_ACCESS_KEY": "test",
		},
	}
	localstackContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer localstackContainer.Terminate(ctx)

	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(REGION))
	if err != nil {
		log.Println("Não foi possível carregar a configuração padrão. Você configurou suas credenciais AWS?")
		log.Fatal(err)
	}
	sqsClient := sqs.NewFromConfig(sdkConfig)
	sqsClient.Options().EndpointResolverV2.ResolveEndpoint(ctx, sqs.EndpointParameters{
		Region:   aws.String(REGION),
		Endpoint: aws.String("http://localhost:4566"),
	})
	returnQueue, err := sqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName: aws.String(QUEUE_NAME),
	})
	fmt.Println(*returnQueue.QueueUrl)
	t.Setenv(ENDPOINT_ENV, *returnQueue.QueueUrl)
	if err != nil {
		log.Printf("Não foi possível criar a fila: %v\n", err)
		return
	}

	fmt.Println("Fila criada com sucesso!")
	fmt.Println(os.Getenv(ENDPOINT_ENV))
	_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String("Ola"),
		QueueUrl:    aws.String(os.Getenv(ENDPOINT_ENV)),
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Mensagem enviada com sucesso")
	}
	msg, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(os.Getenv(ENDPOINT_ENV)),
		MaxNumberOfMessages: *aws.Int32(2),
	})
	if err != nil {
		t.Error("failed to receive message")
	}
	for _, msg := range msg.Messages {
		fmt.Println(*msg.Body)
	}
	assert.NoError(t, err)
}

func TestLocalContainerError(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "localstack/localstack",
		ExposedPorts: []string{"4566/tcp"},
		WaitingFor:   wait.ForLog("Ready."),
		Env: map[string]string{
			"SERVICES":              "sqs",
			"DEFAULT_REGION":        REGION,
			"AWS_ACCESS_KEY_ID":     "test",
			"AWS_SECRET_ACCESS_KEY": "test",
		},
	}
	localstackContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer localstackContainer.Terminate(ctx)

	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(REGION))
	if err != nil {
		log.Println("Não foi possível carregar a configuração padrão. Você configurou suas credenciais AWS?")
		log.Fatal(err)
	}
	sqsClient := sqs.NewFromConfig(sdkConfig)
	sqsClient.Options().EndpointResolverV2.ResolveEndpoint(ctx, sqs.EndpointParameters{
		Region:   aws.String(REGION),
		Endpoint: aws.String("http://localhost:4566"),
	})
	returnQueue, err := sqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName: aws.String(QUEUE_NAME),
	})
	fmt.Println(*returnQueue.QueueUrl)
	if err != nil {
		log.Printf("Não foi possível criar a fila: %v\n", err)
		return
	}

	fmt.Println("Fila criada com sucesso!")
	fmt.Println(os.Getenv(ENDPOINT_ENV))
	_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String("Ola"),
		QueueUrl:    aws.String(os.Getenv(ENDPOINT_ENV)),
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Mensagem enviada com sucesso")
	}
	assert.Error(t, err)
}
