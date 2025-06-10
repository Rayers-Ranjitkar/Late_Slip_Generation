package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates JWT tokens and sets user_id and role in the context
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
            c.Abort()
            return
        }

        // Extract token (Bearer <token>)
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
            c.Abort()
            return
        }
        tokenString := parts[1]

        // Get JWT secret
        secretKey := []byte(os.Getenv("JWT_SECRET"))
        if len(secretKey) == 0 {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret not configured"})
            c.Abort()
            return
        }

        // Parse and validate token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Verify signing method
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return secretKey, nil
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }

        // Extract claims
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }

        // Set user_id and role in context
        c.Set("user_id", claims["user_id"])
        c.Set("role", claims["role"])
        c.Next()
    }
}

// RequireRole restricts access to a specific role
func RequireRole(role string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("role")
        if !exists || userRole != role {
            c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
            c.Abort()
            return
        }
        c.Next()
    }
}