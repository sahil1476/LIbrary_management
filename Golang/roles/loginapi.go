package roles

import (
	"log"
	"main/connection"
	"main/models"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT")

	var user models.User
	//log.Println(user.Email)
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("login >>------>>")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User

	connection.DB.Where("email = ?", user.Email).First(&existingUser)
	var role = existingUser.Role

	if role != "admin" && role != "user" {
		c.JSON(403, gin.H{"error": "Unauthorized access"})
		return
	}
	// log.Println(userMail, "this isu ewr mail")
	log.Println(existingUser.Role, " is role of current Email")

	log.Println("----", role, "--------")

}
