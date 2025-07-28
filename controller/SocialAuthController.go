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
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid state"})
	}

	config := getGoogleOAuthConfig()
	token, err := config.Exchange(context.Background(), ctx.Query("code"))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to exchange token"})
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user info"})
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read user info response"})
	}

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(contents, &userInfo); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse user info"})
	}

	user, err := c.socialAuthService.HandleProviderCallback("google", userInfo.ID, userInfo.Name, userInfo.Email)
	if err != nil {
		log.Printf("Error in HandleProviderCallback: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process user data"})
	}

	jwtToken, err := c.authService.GenerateTokenForUser(user.ID, user.Email)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"error": "Token generation failed",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Login successful",
		"user":    user,
		"token":   jwtToken,
	})
}
