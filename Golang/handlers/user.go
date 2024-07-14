package handlers

import (
	"log"
	"main/connection"
	"main/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RaiseIssueRequest(c *gin.Context) {

	BookID := c.Param("bookid")

	var req struct {
		Email string `json:"email"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "json err"})
		return
	}

	log.Println(req.Email)
	log.Println(BookID)

	var book models.BookInventory
	if err := connection.DB.Where("isbn = ?", BookID).Find(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book is not Available...."})
		return
	}
	if err := connection.DB.Where("email = ?", req.Email).Find(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User is not in DB...."})
		return
	}
	if book.AvailableCopies == 0 || book.AvailableCopies < 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book is not Available Copies=0"})
		return
	}

	var user models.User
	if err := connection.DB.Where("email = ?", req.Email).Find(&user).Error; err != nil {
		// log.Println("inside err")
		return
	}
	if user.Email == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	var requestEvents models.RequestEvents
	requestEvents.BookID = book.ISBN
	requestEvents.ReaderID = user.ID
	requestEvents.RequestDate = time.Now().Format("2006-01-02")
	requestEvents.ApprovalDate = ""
	requestEvents.RequestType = "Pending"

	if err := connection.DB.Create(&requestEvents).Error; err != nil {
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Raise Issue Request"})

}

func Searchbook(c *gin.Context) {
	var books []models.BookInventory
	title := c.Params.ByName("title")
	if err := connection.DB.Where("title LIKE?", "%"+title+"%").Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books": books})
}
