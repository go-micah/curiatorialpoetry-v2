package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-micah/clevelandart"
	"github.com/go-micah/go-bedrock"
)

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

type Payload struct {
	Poem string
	Id   string
}

func handler(ctx context.Context) (*Payload, error) {

	art, err := clevelandart.GetRandomArtwork(true)
	if err != nil {
		return nil, err
	}

	b, err := JSONMarshal(art.Data)
	if err != nil {
		return nil, err
	}

	styles := []string{
		"free verse",
		"haiku",
		"limerick",
		"elegy",
		"couplet",
		"ballad",
		"sonnet",
		"ode",
		"narrative",
		"prose",
		"epic",
	}

	style := styles[rand.Intn(len(styles))]

	prompt := string(b) + "\n"
	prompt += "Use the above <json> document, to inspire a poem in the " + style + " style.\n"
	prompt += "Return the result between <poem> tags.\n"
	prompt += "Also, give your poem a title betweem <title> tags.\n"
	prompt += "Don't include any explanation or introduction."

	claude := bedrock.AnthropicClaude{
		Region:            "us-east-1",
		ModelId:           "anthropic.claude-v2",
		Prompt:            "Human: \n\nHuman: " + prompt + "\n\nAssistant:",
		MaxTokensToSample: 1000,
		TopP:              0.999,
		TopK:              250,
		Temperature:       1,
		StopSequences:     []string{`"\n\nHuman:\"`},
	}

	resp, err := claude.InvokeModel()
	if err != nil {
		return nil, err
	}

	text, err := claude.GetText(resp)
	if err != nil {
		return nil, err
	}

	text = strings.TrimSpace(text)

	payload := Payload{}
	payload.Poem = text

	payload.Id = fmt.Sprint(time.Now().Unix())

	return &payload, nil
}

func main() {
	lambda.Start(handler)
}
