package main

import (
	"github.com/gin-gonic/gin"
)

// func Auth() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		apikey := ctx.GetHeader("X-API-Key")
// 		if apikey == "" {
// 			ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
// 			return
// 		}
// 		ctx.Next()

// 	}
// }

type UserController struct{}

func (uc *UserController) GetUserInfo(ctx *gin.Context) {
	userID := ctx.Param("id")
	ctx.JSON(200, gin.H{"id": userID, "name": "John Doe", "email": "johndoe@gmail.com"})
}
func main() {

	router := gin.Default()

	UserController := &UserController{}

	router.GET("/user/:id", UserController.GetUserInfo)

	// router := gin.Default()

	// public := router.Group("/public")
	// public.GET("/info", func(ctx *gin.Context) {
	// 	ctx.String(200, "General information")
	// })

	// public.GET("/products", func(ctx *gin.Context) {
	// 	ctx.String(200, "Public Products list")
	// })

	// private := router.Group("/private")
	// private.Use(Auth())
	// private.GET("/data", func(ctx *gin.Context) {
	// 	ctx.String(200, "Private data")
	// })

	// private.GET("/create", func(ctx *gin.Context) {
	// 	ctx.String(200, "create a new resource")
	// })

	router.Run(":8080")
}
