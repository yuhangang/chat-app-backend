package gemini_service

import (
	"log"
	"sync"
	"time"

	"github.com/google/generative-ai-go/genai"
)

// CachedSession stores a chat session with metadata
type CachedSession struct {
	Session    *genai.ChatSession
	LastAccess time.Time
}

// ConversationCache manages chat sessions with TTL-based cleanup
type ConversationCache struct {
	Sessions    map[string]CachedSession
	TTL         time.Duration
	CleanupTick time.Duration
	mu          sync.RWMutex
}

// NewConversationCache creates a new conversation cache with TTL
func NewConversationCache(ttl, cleanupTick time.Duration) *ConversationCache {
	cache := &ConversationCache{
		Sessions:    make(map[string]CachedSession),
		TTL:         ttl,
		CleanupTick: cleanupTick,
	}

	// Start cleanup goroutine
	go cache.startCleanup()

	return cache
}

// startCleanup runs a ticker to periodically clean up expired sessions
func (c *ConversationCache) startCleanup() {
	ticker := time.NewTicker(c.CleanupTick)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanup()
	}
}

// cleanup removes expired sessions
func (c *ConversationCache) cleanup() {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()

	for id, session := range c.Sessions {
		if now.Sub(session.LastAccess) > c.TTL {
			delete(c.Sessions, id)
			log.Printf("Session expired and removed: %s", id)
		}
	}
}

// GetOrCreateSession gets an existing session or creates a new one
func (c *ConversationCache) GetOrCreateSession(sessionID string, model *genai.GenerativeModel) *genai.ChatSession {
	c.mu.RLock()
	cachedSession, exists := c.Sessions[sessionID]
	c.mu.RUnlock()

	if exists {
		// Update last access time
		c.mu.Lock()
		c.Sessions[sessionID] = CachedSession{
			Session:    cachedSession.Session,
			LastAccess: time.Now(),
		}
		c.mu.Unlock()
		return cachedSession.Session
	}

	// Create a new session
	session := model.StartChat()

	c.mu.Lock()
	c.Sessions[sessionID] = CachedSession{
		Session:    session,
		LastAccess: time.Now(),
	}
	c.mu.Unlock()

	return session
}
