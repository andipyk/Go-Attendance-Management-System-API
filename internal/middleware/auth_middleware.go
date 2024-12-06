package middleware

import (
	"golang-tes/internal/domain"
	"golang-tes/internal/utils/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
	}
}

// extractToken extracts the token from the Authorization header
func extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			logger.Warn("Missing authorization header",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": domain.ErrUnauthorized.Error(),
			})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, domain.ErrUnauthorized
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil {
			logger.Error("Failed to parse token",
				zap.Error(err),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": domain.ErrUnauthorized.Error(),
			})
			return
		}

		if !token.Valid {
			logger.Warn("Invalid token",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": domain.ErrUnauthorized.Error(),
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Error("Failed to get token claims",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": domain.ErrUnauthorized.Error(),
			})
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			logger.Error("Invalid user ID in token",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": domain.ErrUnauthorized.Error(),
			})
			return
		}

		// Add user information to context
		c.Set("user_id", userID)
		if role, ok := claims["role"].(string); ok {
			c.Set("user_role", role)
		}

		logger.Debug("Authentication successful",
			zap.String("user_id", userID),
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method))

		c.Next()
	}
}

// AdminRequired middleware checks if the user has admin role
func (m *AuthMiddleware) AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || role != domain.RoleAdmin {
			logger.Warn("Unauthorized access to admin endpoint",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("user_id", c.GetString("user_id")))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "admin access required",
			})
			return
		}
		c.Next()
	}
}
