package bedrock

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

const defaultRegion = "eu-central-1"

var Client *bedrockruntime.Client

func init() {

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = defaultRegion
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	Client = bedrockruntime.NewFromConfig(cfg)
}

var verbose *bool

func Chat(input string) string {

	input = strings.TrimSpace(input)

	response, err := send(input)

	if err != nil {
		log.Fatal(err)
	}

	return response
}

const claudePromptFormat = "\n\nHuman: %s\n\nAssistant:"

func send(msg string) (string, error) {

	fmt.Println("[sending message]", msg)
	msg = fmt.Sprintf(claudePromptFormat, msg)
	payload := Request{
		Prompt:            msg,
		MaxTokensToSample: 2048,
		Temperature:       0.2,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	output, err := Client.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		Body:        payloadBytes,
		ModelId:     aws.String("anthropic.claude-v2"),
		ContentType: aws.String("application/json"),
	})

	if err != nil {
		fmt.Printf("Error %v\n", err)
		return "", err
	}

	var resp Response

	err = json.Unmarshal(output.Body, &resp)

	if err != nil {
		return "", err
	}

	return resp.Completion, nil
}

//request/response model

type Request struct {
	Prompt            string   `json:"prompt"`
	MaxTokensToSample int      `json:"max_tokens_to_sample"`
	Temperature       float64  `json:"temperature,omitempty"`
	TopP              float64  `json:"top_p,omitempty"`
	TopK              int      `json:"top_k,omitempty"`
	StopSequences     []string `json:"stop_sequences,omitempty"`
}

type Response struct {
	Completion string `json:"completion"`
}
