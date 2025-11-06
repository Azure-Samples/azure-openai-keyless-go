package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/azure"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Failed to load .env file: %s\n", err)
		os.Exit(1)
	}

	service := os.Getenv("AZURE_OPENAI_SERVICE") // ex: (AZURE_OPENAI_SERVICE).openai.azure.com
	deployment := os.Getenv("AZURE_OPENAI_GPT_DEPLOYMENT")

	if service == "" || deployment == "" {
		fmt.Printf("AZURE_OPENAI_SERVICE and AZURE_OPENAI_GPT_DEPLOYMENT environment variables are empty. See README.")
		os.Exit(1)
	}

	credential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		fmt.Printf("Failed to create DefaultAzureCredential: %s\n", err)
		os.Exit(1)
	}

	client := openai.NewClient(
		option.WithBaseURL(fmt.Sprintf("https://%s.openai.azure.com/openai/v1", service)),
		azure.WithTokenCredential(credential),
	)

	if err != nil {
		fmt.Printf("Failed to create Azure OpenAI client: %s\n", err)
		os.Exit(1)
	}

	response, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		// For Azure OpenAI, the model parameter must be set to the deployment name
		Model:       deployment,
		Temperature: openai.Float(0.7),
		N:           openai.Int(1),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.AssistantMessage("You are a helpful assistant that makes lots of cat references and uses emojis."),
			openai.UserMessage("Write a haiku about a hungry cat who wants tuna"),
		},
	})

	if err != nil {
		fmt.Printf("Failed to get chat completions: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Response:\n%s\n", response.Choices[0].Message.Content)
}
