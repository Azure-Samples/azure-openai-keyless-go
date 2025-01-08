package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/joho/godotenv"
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

	client, err := azopenai.NewClient(
		fmt.Sprintf("https://%s.openai.azure.com", service),
		credential,
		nil)

	if err != nil {
		fmt.Printf("Failed to create Azure OpenAI client: %s\n", err)
		os.Exit(1)
	}

	response, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		// For Azure OpenAI, the model parameter must be set to the deployment name
		DeploymentName: &deployment,
		Temperature:    to.Ptr[float32](0.7),
		N:              to.Ptr[int32](1),
		Messages: []azopenai.ChatRequestMessageClassification{
			&azopenai.ChatRequestAssistantMessage{
				Content: azopenai.NewChatRequestAssistantMessageContent("You are a helpful assistant that makes lots of cat references and uses emojis."),
			},
			&azopenai.ChatRequestUserMessage{
				Content: azopenai.NewChatRequestUserMessageContent("Write a haiku about a hungry cat who wants tuna"),
			},
		},
	}, nil)

	if err != nil {
		fmt.Printf("Failed to get chat completions: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Response:\n%s\n", *response.Choices[0].Message.Content)
}
