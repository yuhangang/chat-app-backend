package types

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/generative-ai-go/genai"
)

type ChatMessages struct {
	Messages string
}

type HttpHandlerV1 interface {
	RegisterRouter(router *http.ServeMux)
	CreateMessage(w http.ResponseWriter, r *http.Request)
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

type JwtPayload struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JwtService interface {
	GenerateTokens(userID uint) (JwtPayload, error)
	ValidateAccessToken(tokenString string) (Claims, error)
	ValidateRefreshToken(tokenString string) (Claims, error)
	RefreshAccessToken(refreshToken string) (string, error)
}

type HttpServiceV1 interface {
	CallGemini(context.Context, string, string, []*genai.Content) (GeminiApiResponse, error)
	SendFileWithText(ctx context.Context, string, prompt string, history []*genai.Content, tempFilePath string) (GeminiApiResponse, error)
}

type GeminiApiResponse struct {
	Response  string `json:"response"`
	SessionID string `json:"session_id"`
}
