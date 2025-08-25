package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"bookstore-api/controllers"
	"bookstore-api/middlewares"
)

func Register(r *gin.Engine) {
	// 🚀 Default route supaya tidak 404 di /
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "🚀 Bookstore API is running on Railway",
		})
	})

	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/register", controllers.Register)
			users.POST("/login", controllers.Login)
		}

		// protected
		api.Use(middlewares.JWTAuth())

		cats := api.Group("/categories")
		{
			cats.GET("/", controllers.GetCategories)
			cats.POST("/", controllers.CreateCategory)
			cats.GET("/:id", controllers.GetCategoryByID)
			cats.PUT("/:id", controllers.UpdateCategory)
			cats.DELETE("/:id", controllers.DeleteCategory)
			cats.GET("/:id/books", controllers.GetBooksByCategory)
		}

		books := api.Group("/books")
		{
			books.GET("/", controllers.GetBooks)
			books.POST("/", controllers.CreateBook)
			books.GET("/:id", controllers.GetBookByID)
			books.PUT("/:id", controllers.UpdateBook)
			books.DELETE("/:id", controllers.DeleteBook)
		}
	}
}
