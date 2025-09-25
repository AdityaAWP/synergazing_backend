package controller

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"synergazing.com/synergazing/helper"
	"synergazing.com/synergazing/service"
)

var (
	googleOauthConfig *oauth2.Config
	configOnce        sync.Once
)

func getGoogleOAuthConfig() *oauth2.Config {
	configOnce.Do(func() {
		googleOauthConfig = &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		}

		log.Printf("Google OAuth Config initialized:")
		log.Printf("ClientID: %s", googleOauthConfig.ClientID)
		log.Printf("RedirectURL: %s", googleOauthConfig.RedirectURL)
		log.Printf("Client Secret length: %d", len(googleOauthConfig.ClientSecret))
		if googleOauthConfig.RedirectURL == "" {
			log.Printf("WARNING: GOOGLE_REDIRECT_URI is empty!")
		}
		if googleOauthConfig.ClientID == "" {
			log.Printf("WARNING: GOOGLE_CLIENT_ID is empty!")
		}
		if googleOauthConfig.ClientSecret == "" {
			log.Printf("WARNING: GOOGLE_CLIENT_SECRET is empty!")
		}
	})
	return googleOauthConfig
}

const oauthStateString = "random-string-for-security"

type SocialController struct {
	socialAuthService *service.SocialAuthService
	authService       *service.AuthService
}

func NewSocialController(sas *service.SocialAuthService, as *service.AuthService) *SocialController {
	return &SocialController{
		socialAuthService: sas,
		authService:       as,
	}
}

func (c *SocialController) GoogleLogin(ctx *fiber.Ctx) error {
	config := getGoogleOAuthConfig()
	url := config.AuthCodeURL(oauthStateString)
	log.Printf("Generated Google OAuth URL: %s", url)
	log.Printf("Expected callback URL: %s", config.RedirectURL)
	return ctx.Redirect(url)
}

func (c *SocialController) GoogleCallback(ctx *fiber.Ctx) error {
	log.Printf("GoogleCallback called with URL: %s", ctx.OriginalURL())
	log.Printf("State received: %s", ctx.Query("state"))
	log.Printf("Code received: %s", ctx.Query("code"))

	if ctx.Query("state") != oauthStateString {
		log.Printf("State mismatch! Expected: %s, Got: %s", oauthStateString, ctx.Query("state"))
		errorURL := helper.BuildOAuthErrorURL("invalid_state")
		return ctx.Redirect(errorURL)
	}

	config := getGoogleOAuthConfig()
	token, err := config.Exchange(context.Background(), ctx.Query("code"))
	if err != nil {
		log.Printf("Failed to exchange token: %v", err)
		errorURL := helper.BuildOAuthErrorURL("token_exchange_failed")
		return ctx.Redirect(errorURL)
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		errorURL := helper.BuildOAuthErrorURL("user_info_failed")
		return ctx.Redirect(errorURL)
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Failed to read user info response: %v", err)
		errorURL := helper.BuildOAuthErrorURL("user_info_read_failed")
		return ctx.Redirect(errorURL)
	}

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(contents, &userInfo); err != nil {
		log.Printf("Failed to parse user info: %v", err)
		errorURL := helper.BuildOAuthErrorURL("user_info_parse_failed")
		return ctx.Redirect(errorURL)
	}

	user, err := c.socialAuthService.HandleProviderCallback("google", userInfo.ID, userInfo.Name, userInfo.Email)
	if err != nil {
		log.Printf("Error in HandleProviderCallback: %v", err)
		errorURL := helper.BuildOAuthErrorURL("user_processing_failed")
		return ctx.Redirect(errorURL)
	}

	jwtToken, err := c.authService.GenerateTokenForUser(user.ID, user.Email)
	if err != nil {
		log.Printf("Token generation failed: %v", err)
		errorURL := helper.BuildOAuthErrorURL("token_generation_failed")
		return ctx.Redirect(errorURL)
	}

	// Redirect to frontend with success data
	redirectURL := helper.BuildOAuthSuccessURL(jwtToken, user.ID, user.Name, user.Email)
	log.Printf("Redirecting to frontend: %s", redirectURL)
	return ctx.Redirect(redirectURL)
}

// OAuthSuccess handles successful OAuth redirects with query parameters
func (c *SocialController) OAuthSuccess(ctx *fiber.Ctx) error {
	token := ctx.Query("token")
	userID := ctx.Query("user_id")
	userName := ctx.Query("user_name")
	userEmail := ctx.Query("user_email")

	if token == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing token parameter",
		})
	}

	return ctx.JSON(fiber.Map{
		"success":    true,
		"message":    "OAuth authentication successful",
		"token":      token,
		"user_id":    userID,
		"user_name":  userName,
		"user_email": userEmail,
	})
}

// OAuthError handles OAuth error redirects
func (c *SocialController) OAuthError(ctx *fiber.Ctx) error {
	errorType := ctx.Query("error")
	errorDescription := ctx.Query("error_description", "An error occurred during authentication")

	log.Printf("OAuth Error: %s - %s", errorType, errorDescription)

	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"error":   errorType,
		"message": errorDescription,
	})
}
