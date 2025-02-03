package auth

import (
	"net/http"

	"github.com/dmdhrumilmistry/defect-detect/pkg/config"
	"github.com/dmdhrumilmistry/defect-detect/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/rs/zerolog/log"
)

type AuthHandler struct {
	store types.AuthStore
}

func NewAuthHandler(store types.AuthStore) *AuthHandler {
	handler := &AuthHandler{
		store: store,
	}

	handler.InitAuth()

	return handler
}

// InitAuth initializes the authentication providers
func (a *AuthHandler) InitAuth() {
	log.Info().Msg("Initializing Auth Providers")
	domainUrl := "http://localhost:8080"
	googleRedirectUri := "/auth/google/callback"

	googleRedirectUrl := domainUrl + googleRedirectUri
	log.Info().Msgf("Google Redirect Url: %s", googleRedirectUrl)

	goth.UseProviders(
		google.New(
			config.DefaultConfig.GoogleClientId,
			config.DefaultConfig.GoogleClientSecret,
			googleRedirectUrl, // Redirect URL
			"email", "profile",
		),
	)
	log.Info().Msg("Initialized Auth Providers Successfully")
}

func (a *AuthHandler) RegisterRoutes(r *gin.Engine) {
	// Google auth
	r.GET("/auth/", a.GoogleAuthHandler) // GET http://domain:8080/auth/?provider=google
	r.GET("/auth/google/callback", a.GoogleCallbackHandler)

	log.Info().Msg("Auth Providers routes registered")
}

// AuthHandler redirects users to Google login
func (a *AuthHandler) GoogleAuthHandler(c *gin.Context) {
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// CallbackHandler handles Google auth callback
func (a *AuthHandler) GoogleCallbackHandler(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Error().Err(err).Msg("failed to complete google oauth flow")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	// DEBUG: log user details
	log.Info().Any("user", user).Msg("")

	// create user if not exists and return paseto token
	c.JSON(http.StatusAccepted, gin.H{"user": user})
}
