package gemini_service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/yuhangang/chat-app-backend/types"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiServiceV1 struct {
	client *genai.Client
}

func NewGeminiServiceV1(ctx context.Context) *GeminiServiceV1 {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))

	if err != nil {
		log.Fatal(err)
	}

	return &GeminiServiceV1{
		client: client,
	}
}

func (s *GeminiServiceV1) CallGemini(ctx context.Context, prompt string, history []*genai.Content) (types.GeminiApiResponse, error) {

	model := s.client.GenerativeModel("gemini-2.0-flash")
	// Ask the model to respond with JSON.
	model.ResponseMIMEType = "application/json"
	// Specify the schema.
	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"response": {Type: genai.TypeString},
			},
		},
	}
	cs := model.StartChat()
	cs.History = history

	resp, err := cs.SendMessage(ctx, genai.Text(prompt))

	if err != nil {
		log.Fatalf("Error generating content: %v", err)
	}

	var result string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			var response []map[string]string
			err := json.Unmarshal([]byte(txt), &response)
			if err != nil {
				log.Fatal(err)
			}
			result = response[0]["response"]
		}
	}

	return types.GeminiApiResponse{
		Response: result,
	}, nil

}

func (s *GeminiServiceV1) SendFileWithText(ctx context.Context, prompt string, history []*genai.Content, tempFilePath string) (types.GeminiApiResponse, error) {
	log.Println("SendFileWithText")
	uploadedFile, err := s.uploadFile(ctx, tempFilePath)
	if err != nil {
		return types.GeminiApiResponse{}, err
	}
	defer s.client.DeleteFile(ctx, uploadedFile.Name)

	model := s.client.GenerativeModel("gemini-2.0-flash")
	// Ask the model to respond with JSON.
	model.ResponseMIMEType = "application/json"
	// Specify the schema.
	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"response": {Type: genai.TypeString},
			},
		},
	}
	cs := model.StartChat()
	cs.History = history

	resp, err := cs.SendMessage(ctx, genai.Text(prompt), genai.FileData{URI: uploadedFile.URI})

	if err != nil {
		log.Fatalf("Error generating content: %v", err)
	}

	var result string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			var response []map[string]string
			err := json.Unmarshal([]byte(txt), &response)
			if err != nil {
				log.Fatal(err)
			}
			result = response[0]["response"]
		}
	}

	return types.GeminiApiResponse{
		Response: result,
	}, nil
}

// Helper function to upload the file and return the uploaded file details
func (s *GeminiServiceV1) uploadFile(ctx context.Context, tempFilePath string) (*genai.File, error) {
	file, err := s.client.UploadFileFromPath(ctx, tempFilePath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	return file, nil
}
