package handlers

import (
	"main/connection"
	"main/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var user models.User
var book models.BookInventory
var issueRegistery models.IssueRegistery //accept or reject request
var request models.RequestEvents         //request by user to read and search book

func Create_user(c *gin.Context) {

	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := connection.DB.Where("id = ?", user.ID).Find(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Request ID Found"})
		return
	}
	if err := connection.DB.Where("email = ?", user.Email).Find(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Email Exists"})
		return
	}
	err := connection.DB.Create(&user).Error

	if err != nil {
		c.JSON(400, "Error creating user")
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "USER created"})
}

func Get_user(c *gin.Context) {
	var userall []models.User
	err := connection.DB.Find(&userall).Error
	if err != nil {
		c.JSON(400, "User not FOUND in DATA-BASE")
		return
	}
	c.IndentedJSON(http.StatusOK, userall)
}

func Remove_user(c *gin.Context) {

	id := c.Param("id")
	var user models.User
	if err := connection.DB.Where("id = ?", user.ID).Find(&user).Error; err == nil {
		if _, ok := err.(interface{}); ok {
			c.JSON(http.StatusNotFound, gin.H{"message": "User ID not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error occurred while fetching user"})
		return
	}
	err := connection.DB.Delete(&user, id).Error
	if err != nil {
		c.JSON(400, "User not DELETED from DATA-BASE")
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User Deleted successfully"})

}
func Create_book(c *gin.Context) {

	if err := c.ShouldBindJSON(&book); err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := connection.DB.Create(&book).Error

	if err != nil {
		c.JSON(400, "Book not created in Library")
		return
	}
	c.JSON(http.StatusOK, "Book created successfully")

}
func Show_book(c *gin.Context) {
	var allbooks []models.BookInventory
	err := connection.DB.Find(&allbooks).Error
	if err != nil {
		c.JSON(404, "Books not FOUND in DATA-BASE")
		return
	}
	c.IndentedJSON(http.StatusOK, allbooks)
}

func Remove_book(c *gin.Context) {
	isbn, _ := strconv.Atoi(c.Param("isbn")) //int conversion
	if err := connection.DB.Where("isbn = ?", uint(isbn)).Find(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No Book found"})
		return
	}
	if book.ISBN != uint(isbn) {
		c.JSON(http.StatusNotFound, gin.H{"message": "No Book with Given ISBN number"})
		return
	} else if book.AvailableCopies == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Available Copies = 0 can't delete"})
		return
	}
	book.AvailableCopies--
	connection.DB.Save(&book)
}
func Update_book(c *gin.Context) {
	isbn := c.Params.ByName("isbn")

	err := connection.DB.Find(&book, isbn).Error
	if err != nil {
		c.JSON(400, "Book with ISBN number "+isbn+" not found")
		return
	}

	var updatebook models.BookInventory
	if err := c.ShouldBindJSON(&updatebook); err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Set default values for missing fields
	if updatebook.TotalCopies == 0 {
		updatebook.TotalCopies = book.TotalCopies
	}
	if updatebook.AvailableCopies == 0 {
		updatebook.AvailableCopies = book.AvailableCopies
	}
	if updatebook.Authors == "" {
		updatebook.Authors = book.Authors
	}
	if updatebook.LibID == 0 {
		updatebook.LibID = book.LibID
	}
	if updatebook.Publisher == "" {
		updatebook.Publisher = book.Publisher
	}
	if updatebook.Title == "" {
		updatebook.Title = book.Title
	}
	if updatebook.Version == 0 {
		updatebook.Version = book.Version
	}

	connection.DB.Model(&book).Where("isbn =?", isbn).Updates(&updatebook)

	c.JSON(http.StatusOK, "Book with ISBN number "+isbn+" updated")
	return

}

func List_issue_request(c *gin.Context) {
	var requestbookall []models.RequestEvents
	err := connection.DB.Find(&requestbookall).Error
	if err != nil {
		c.JSON(404, "Books not FOUND in DATA-BASE")
		return
	}
	c.IndentedJSON(http.StatusOK, requestbookall)

}

type DetailReqID struct {
	ReqID int `json:"reqid"`
}

func Approve_request(c *gin.Context) {
	var ri DetailReqID
	if err := c.BindJSON(&ri); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var request models.RequestEvents
	if err := connection.DB.Where("req_id = ?", ri.ReqID).Find(&request).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Request ID Not Found"})
		return
	}
	if request.RequestType == "Issued" {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book is already approved"})
		return
	}
	var book models.BookInventory
	if err := connection.DB.Where("isbn = ?", request.BookID).Find(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Book is not available "})
		return
	}
	if book.TotalCopies < 1 {
		c.JSON(http.StatusAccepted, gin.H{"message": "Book is not available "})
		return
	}
	if book.AvailableCopies <= 0 && book.TotalCopies >= 1 {
		c.JSON(http.StatusAccepted, gin.H{"message": " All books are already issued"})
		return
	}
	// after the book approval save the book copies left
	book.AvailableCopies = book.AvailableCopies - 1
	if err := connection.DB.Where("isbn = ?", request.BookID).Save(&book).Error; err != nil {
		c.JSON(http.StatusAccepted, gin.H{"message": "error in updating book inventory"})
		return
	}

	request.ApprovalDate = time.Now().Format("2006-01-02T15:04:05")
	request.ApproverID = user.ID
	request.RequestType = "Approve"
	//saving the RequestEvent table values..
	connection.DB.Where("req_id = ?", ri.ReqID).Save(&request)

	var issueRegistery models.IssueRegistery
	issueRegistery.ISBN = request.BookID
	issueRegistery.ReaderID = request.ReaderID
	issueRegistery.IssueApproverID = user.ID
	issueRegistery.IssueStatus = "Approve"
	issueRegistery.IssueDate = time.Now().Format("2006-01-02T15:04:05")
	issueRegistery.ExpectedReturnDate = time.Now().AddDate(0, 0, 30)

	if err := connection.DB.Create(&issueRegistery).Error; err != nil {
		c.JSON(http.StatusAccepted, gin.H{"message": "Issue registry not saved"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book approved"})
}

func Reject_request(c *gin.Context) {
	var ri DetailReqID
	if err := c.BindJSON(&ri); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var request models.RequestEvents
	if err := connection.DB.Where("req_id = ?", ri.ReqID).Find(&request).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Request ID Not Found"})
		return
	}

	if err := connection.DB.Where("isbn = ?", request.BookID).Find(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No Book with Given BookID found"})
		return
	}
	request.ApprovalDate = time.Now().Format("2006-01-02T15:04:05")
	request.ApproverID = user.ID
	request.RequestType = "Rejected"
	//saving the RequestEvent table values..
	connection.DB.Where("req_id = ?", ri.ReqID).Save(&request)

	var issueRegistery models.IssueRegistery
	issueRegistery.ISBN = request.BookID
	issueRegistery.ReaderID = request.ReaderID
	issueRegistery.IssueApproverID = user.ID
	issueRegistery.IssueStatus = "Rejected"
	issueRegistery.ReturnDate = "none"
	issueRegistery.IssueDate = time.Now().Format("2006-01-02T15:04:05")
	issueRegistery.ExpectedReturnDate = time.Now().AddDate(0, 0, 30)

	if err := connection.DB.Create(&issueRegistery).Error; err != nil {
		c.JSON(http.StatusAccepted, gin.H{"message": "Issue registry not saved"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book Request Rejected!"})
}
