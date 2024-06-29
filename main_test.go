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
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestLocalContainer(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "localstack/localstack",
		ExposedPorts: []string{"4566/tcp"},
		WaitingFor:   wait.ForLog("Ready."),
		Env: map[string]string{
			"SERVICES":              "sqs",
			"DEFAULT_REGION":        "us-east-1",
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

	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		log.Println("Não foi possível carregar a configuração padrão. Você configurou suas credenciais AWS?")
		log.Fatal(err)
	}
	sqsClient := sqs.NewFromConfig(sdkConfig)
	sqsClient.Options().EndpointResolverV2.ResolveEndpoint(ctx, sqs.EndpointParameters{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://localhost:4566"),
	})
	queueName := "minha-fila-sqs"
	returnQueue, err := sqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName: &queueName,
	})
	fmt.Println(*returnQueue.QueueUrl)
	t.Setenv("AWS_ENDPOINT", *returnQueue.QueueUrl)
	if err != nil {
		log.Printf("Não foi possível criar a fila: %v\n", err)
		return
	}

	fmt.Println("Fila criada com sucesso!")
	fmt.Println(os.Getenv("AWS_ENDPOINT"))
	_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String("Ola"),
		QueueUrl:    aws.String(os.Getenv("AWS_ENDPOINT")),
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Mensagem enviada com sucesso")
	}
}

