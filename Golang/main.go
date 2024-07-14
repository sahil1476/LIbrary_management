package main

import (
	"main/connection"
	"main/handlers"
	"main/models"
	"main/roles"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	connection.Connection()

	connection.DB.AutoMigrate(&models.Library{})
	connection.DB.AutoMigrate(&models.BookInventory{})
	connection.DB.AutoMigrate(&models.IssueRegistery{})
	connection.DB.AutoMigrate(&models.RequestEvents{})
	connection.DB.AutoMigrate(&models.User{})

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "OPTIONS"}

	router.Use(cors.New(config))

	//Owner
	router.POST("/createlibrary", handlers.Createlibrary)

	//Admin
	router.POST("/usercreate", handlers.Create_user)
	router.GET("/userlist", handlers.Get_user)
	router.DELETE("/userdelete/:id", handlers.Remove_user)

	router.POST("/createbook", handlers.Create_book)
	router.GET("/showbook", handlers.Show_book)
	router.DELETE("/removebook/:isbn", handlers.Remove_book)
	router.PATCH("/updatebook/:isbn", handlers.Update_book)

	router.POST("/login", roles.Login)

	router.GET("/showrequest", handlers.List_issue_request)
	router.POST("/approverequest", handlers.Approve_request)
	router.POST("/rejectrequest", handlers.Reject_request)
	//User
	router.POST("/requestbook/:bookid", handlers.RaiseIssueRequest)
	router.GET("/search/:title", handlers.Searchbook)
	router.Run("localhost:8080")
	// router.Run("localhost:5500")
}
