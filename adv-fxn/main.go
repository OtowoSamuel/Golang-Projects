package main

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User created successfully"})
}
func Search(c *gin.Context) {
	query := c.Query("query")
	limit := c.DefaultQuery("limit", "10")
	c.JSON(200, gin.H{"query": query, "limit": limit})
}

// Register a new user
func Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Hash password (use bcrypt in production)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hashedPassword)

	if err := db.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "Email already exists"})
		return
	}
	c.JSON(200, gin.H{"message": "User registered"})
}

// Login and generate JWT token
func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("your-secret-key"))

	c.JSON(200, gin.H{"token": tokenString})
}

// Middleware to protect routes

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		userID := uint(claims["sub"].(float64))

		var user User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(401, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Process request

		// Check for errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			statusCode := c.Writer.Status()

			// Custom error responses
			switch {
			case errors.Is(err.Err, gorm.ErrRecordNotFound):
				c.JSON(404, gin.H{"error": "Resource not found"})
			case statusCode >= 500:
				c.JSON(500, gin.H{"error": "Internal server error"})
			default:
				c.JSON(statusCode, gin.H{"error": err.Error()})
			}
		}
	}
}

func SetupLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}

func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		logger.Info("Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Since(start)),
		)
	}
}

var c = cache.New(5*time.Minute, 10*time.Minute) // Default 5-min TTL

func CacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		key := c.Request.URL.String()
		if cached, found := c.Get(key); found {
			c.JSON(200, cached)
			c.Abort()
			return
		}

		c.Next()

		// Cache the response
		if c.Writer.Status() == 200 {
			c.Set(key, c.Keys["response"], cache.DefaultExpiration)
		}
	}
}

func main() {
	logger := SetupLogger()
	defer logger.Sync()

	router := gin.Default()
	router.Use(RequestLogger(logger))

	router.Use(CacheMiddleware())

	router.GET("/todos", func(c *gin.Context) {
		var todos []Todo
		db.Find(&todos)
		c.Keys["response"] = todos // Store for caching
		c.JSON(200, todos)
	})
}
