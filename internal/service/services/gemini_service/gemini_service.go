package gemini_service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yuhangang/chat-app-backend/types"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiServiceV1 struct {
	client *genai.Client
	cache  *ConversationCache
}

func NewGeminiServiceV1(ctx context.Context) *GeminiServiceV1 {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))

	if err != nil {
		log.Fatal(err)
	}

	// Create cache with 30 minute TTL and 5 minute cleanup interval
	cache := NewConversationCache(30*time.Minute, 5*time.Minute)

	return &GeminiServiceV1{
		client: client,
		cache:  cache,
	}
}

func (s *GeminiServiceV1) CallGemini(ctx context.Context, sessionID string, prompt string, history []*genai.Content) (types.GeminiApiResponse, error) {
	model := s.client.GenerativeModel("gemini-2.0-flash")

	// Configure model response format
	model.ResponseMIMEType = "application/json"

	// TODO: handle message temperature
	//temp := float32(1.0)
	//model.Temperature = &temp
	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"response": {Type: genai.TypeString},
			},
		},
	}

	// Get or create a session from cache
	cs := s.cache.GetOrCreateSession(sessionID, model)

	// If history is provided and different from cached history, update it
	if history != nil && len(history) > 0 {
		cs.History = history
	}

	resp, err := cs.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		return types.GeminiApiResponse{}, fmt.Errorf("error generating content: %w", err)
	}

	var result string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			var response []map[string]string
			err := json.Unmarshal([]byte(txt), &response)
			if err != nil {
				return types.GeminiApiResponse{}, fmt.Errorf("error unmarshaling response: %w", err)
			}
			result = response[0]["response"]
		}
	}

	return types.GeminiApiResponse{
		Response:  result,
		SessionID: sessionID,
	}, nil
}

func (s *GeminiServiceV1) SendFileWithText(ctx context.Context, sessionID string, prompt string, history []*genai.Content, tempFilePath string) (types.GeminiApiResponse, error) {
	log.Println("SendFileWithText")
	uploadedFile, err := s.uploadFile(ctx, tempFilePath)
	if err != nil {
		return types.GeminiApiResponse{}, err
	}
	defer s.client.DeleteFile(ctx, uploadedFile.Name)

	model := s.client.GenerativeModel("gemini-2.0-flash")

	// Configure model response format
	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"response": {Type: genai.TypeString},
			},
		},
	}

	// Get or create a session from cache
	cs := s.cache.GetOrCreateSession(sessionID, model)

	// If history is provided and different from cached history, update it
	if history != nil && len(history) > 0 {
		cs.History = history
	}

	resp, err := cs.SendMessage(ctx, genai.Text(prompt), genai.FileData{URI: uploadedFile.URI})
	if err != nil {
		return types.GeminiApiResponse{}, fmt.Errorf("error generating content: %w", err)
	}

	var result string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			var response []map[string]string
			err := json.Unmarshal([]byte(txt), &response)
			if err != nil {
				return types.GeminiApiResponse{}, fmt.Errorf("error unmarshaling response: %w", err)
			}
			result = response[0]["response"]
		}
	}

	return types.GeminiApiResponse{
		Response:  result,
		SessionID: sessionID,
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
