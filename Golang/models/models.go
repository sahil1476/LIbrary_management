package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Role               string `json:"role"`
	jwt.StandardClaims `json:"-"`
}

type Library struct {
	ID   uint    `gorm:"primaryKey;" json:"id"`
	Name *string `json:"name"`
}
type User struct {
	ID            uint    `gorm:"primaryKey;" json:"id"`
	Name          string  `json:"name"`
	Email         *string `json:"email"`
	ContactNumber uint    `json:"contactnumber"`
	Role          string  `gorm:"default:user;" json:"role"`
	LibID         uint    `gorm:"default:2000;" json:"libid"`
}

// type Claims struct {
// 	Role string `json:"role"'

// }
type BookInventory struct {
	ISBN            uint   `gorm:"primaryKey;" json:"isbn"`
	LibID           uint   `json:"libid"`
	Title           string `json:"title"`
	Authors         string `json:"author"`
	Publisher       string `json:"publisher"`
	Version         int    `json:"version"`
	TotalCopies     uint   `json:"totalcopies"`
	AvailableCopies uint   `json:"availablecopies"`
}
type RequestEvents struct {
	ReqID        uint   `gorm:"primaryKey;" json:"reqid"`
	BookID       uint   `json:"bookid"` //bookid == isbn
	ReaderID     uint   `json:"readerid"`
	RequestDate  string `json:"requestdate"`
	ApprovalDate string `json:"apprivaldate"`
	ApproverID   uint   `json:"approverid"`
	RequestType  string `gorm:"default:Pending;" json:"requesttype"`
}
type IssueRegistery struct {
	IssueID            uint      `gorm:"primaryKey;" json:"issueid"`
	ISBN               uint      `json:"isbn"`
	ReaderID           uint      `json:"readerid"`
	IssueApproverID    uint      `json:"issueapproverid"`
	IssueStatus        string    `gorm:"default:Pending;" json:"issuestatus"`
	IssueDate          string    `json:"issuedate"`
	ExpectedReturnDate time.Time `json:"expectedreturndate"`
	ReturnDate         string    `json:"returndate"`
	ReturnApproverID   uint      `json:"returnapproverid"`
}
