package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lazzato/server/config"
	"github.com/lazzato/server/initializers"
	"github.com/lazzato/server/models"
	"github.com/lazzato/server/utils"
	"golang.org/x/oauth2"
)

func GoogleAuth(c *gin.Context) {
	redirectURL := c.Query("redirect")
	if redirectURL != "" {
		// Store redirect URL in cookie (valid for 10 mins)
		c.SetCookie("redirect_after_login", redirectURL, 600, "/", "", false, true)
	}

	// Generate Google OAuth URL and redirect
	url := config.GoogleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(302, url)
}

func GoogleAuthCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	// 1. Exchange code for Google token
	tok, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token exchange failed"})
		return
	}

	// 2. Get user info from Google
	client := config.GoogleOAuthConfig.Client(context.Background(), tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var ui struct {
		Sub     string `json:"sub"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
		Locale  string `json:"locale"`
	}
	if err := json.Unmarshal(body, &ui); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user info"})
		return
	}

	// 3. Check if user exists in DB
	var user models.User
	result := initializers.DB.Where("google_id = ?", ui.Sub).First(&user)

	if result.Error != nil {
		// User not found — create a new owner
		user = models.User{
			Email:    ui.Email,
			Name:     ui.Name,
			GoogleID: &ui.Sub,
			Role:     "owner",
		}
		if err := initializers.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	// 4. Generate JWT tokens
	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create access token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create refresh token"})
		return
	}

	// 5. Set refresh token in secure HttpOnly cookie (valid for 30 days)
	c.SetCookie(
		"refresh_token",  // name
		refreshToken,     // value
		60*60*24*30,      // max-age (30 days)
		"/",              // path
		"",               // domain
		false,            // secure (use true if on HTTPS)
		true,             // httpOnly
	)

	c.SetCookie(
	"access_token",
	accessToken,
	60*30,     // 30 minutes
	"/",
	"",
	false,     // true if using HTTPS
	true,      // HttpOnly
)

	redirectURL, err := c.Cookie("redirect_after_login")
if err != nil || redirectURL == "" {
	redirectURL = "http://localhost:3000/"
}
c.SetCookie("redirect_after_login", "", -1, "/", "", false, true)
c.Redirect(http.StatusFound, redirectURL)

}



func RefreshAccessToken(c *gin.Context) {
	rt, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token"})
		return
	}

	token, err := jwt.Parse(rt, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["sub"].(float64))

	newAccessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}

func GetMeHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	var user models.User
	if err := initializers.DB.First(&user, userID.(uint)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     user.ID,
		"email":  user.Email,
		"name":   user.Name,
		"role":   user.Role,
	})
}

