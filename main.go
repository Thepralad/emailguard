package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// EmailValidationResponse represents the API response structure
type EmailValidationResponse struct {
	Email        string `json:"email"`
	IsValid      bool   `json:"is_valid"`
	Reason       string `json:"reason,omitempty"`
	Provider     string `json:"provider,omitempty"`
	IsFree       bool   `json:"is_free"`
	IsDisposable bool   `json:"is_disposable"`
}

// Add this list at the top of the file with other global variables
var disposableDomains = map[string]bool{
	"apklamp.com":       true,
	"temp-mail.org":     true,
	"tempmail.com":      true,
	"guerrillamail.com": true,
	// Add more known disposable domains
}

// BasicEmailValidation performs initial syntax validation
func BasicEmailValidation(email string) (bool, string) {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return false, "Invalid email format"
	}
	return true, ""
}

// ValidateEmailHandler handles the email validation endpoint
func ValidateEmailHandler(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email parameter is required"})
		return
	}

	// Get domain from email
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		c.JSON(http.StatusOK, EmailValidationResponse{
			Email:   email,
			IsValid: false,
			Reason:  "Invalid email format",
		})
		return
	}
	domain := parts[1]

	// Check against known disposable domains
	isDisposable := disposableDomains[domain]

	// Basic validation
	if valid, reason := BasicEmailValidation(email); !valid {
		c.JSON(http.StatusOK, EmailValidationResponse{
			Email:   email,
			IsValid: false,
			Reason:  reason,
		})
		return
	}

	// Initialize response
	response := EmailValidationResponse{
		Email:   email,
		IsValid: true,
	}

	// Mailboxlayer API validation
	mailboxLayerURL := fmt.Sprintf("http://apilayer.net/api/check?access_key=%s&email=%s",
		os.Getenv("MAILBOXLAYER_API_KEY"), email)

	resp, err := http.Get(mailboxLayerURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error validating email"})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing response"})
		return
	}

	// Update response with Mailboxlayer results
	if score, ok := result["score"].(float64); ok {
		response.IsValid = score > 0.7
	}
	if free, ok := result["free"].(bool); ok {
		response.IsFree = free
	}

	// Update the response to use our disposable check
	response.IsDisposable = isDisposable || response.IsDisposable // Combine both checks

	c.JSON(http.StatusOK, response)
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Set up Gin router
	router := gin.Default()

	// Serve static files - Add these lines
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Define routes
	router.GET("/validate", ValidateEmailHandler)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
