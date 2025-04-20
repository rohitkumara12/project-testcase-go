package controllers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"test-case/models"
	"test-case/views"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	var registUser views.UserRegistRequest
	var responsemodel views.Baseresponse

	// Binding JSON request
	err := c.ShouldBindJSON(&registUser)
	if err != nil {
		responsemodel.Status = http.StatusBadRequest
		responsemodel.Data = gin.H{"error": "Invalid input format"}
		c.JSON(http.StatusBadRequest, responsemodel)
		return
	}
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		responsemodel.Status = http.StatusInternalServerError
		responsemodel.Data = gin.H{"error": "Failed to process password"}
		c.JSON(http.StatusInternalServerError, responsemodel)
		return
	}
	// Assign values to the User model

	user := models.DataNewUser{
		Name:     registUser.Name,
		Email:    registUser.Email,
		Password: string(hashedPassword),
	}

	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "failed to create new user"})
		return
	}
	address := models.UserAddressDetails{

		City:    registUser.City,
		Country: registUser.Country,
		Kodepos: registUser.Kodepos,
		Street:  registUser.Street,
	}
	if err := models.DB.Create(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create new address for user"})
	}

	// Success response
	responsemodel.Status = http.StatusOK
	responsemodel.Data = gin.H{
		"message": "Registration successful",
		"user": gin.H{
			"user":    user,
			"address": address,
		},
	}

	c.JSON(http.StatusOK, responsemodel)
}

func Login(c *gin.Context) {
	var responsemodel views.Baseresponse

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		responsemodel.Status = http.StatusUnauthorized
		responsemodel.Data = gin.H{"error": "Authorization header is required"}
		c.JSON(http.StatusUnauthorized, responsemodel)
		return
	}

	// Memastikan format header menggunakan "Basic"
	if !strings.HasPrefix(authHeader, "Basic ") {
		responsemodel.Status = http.StatusUnauthorized
		responsemodel.Data = gin.H{"error": "Invalid Authorization header format"}
		c.JSON(http.StatusUnauthorized, responsemodel)
		return
	}

	// Parsing Basic Auth credentials
	encodedCredentials := strings.TrimPrefix(authHeader, "Basic ")
	decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		responsemodel.Status = http.StatusBadRequest
		responsemodel.Data = gin.H{"error": "Invalid base64 encoding in Authorization header"}
		c.JSON(http.StatusBadRequest, responsemodel)
		return
	}

	credentials := strings.SplitN(string(decodedCredentials), ":", 2)
	if len(credentials) != 2 {
		responsemodel.Status = http.StatusBadRequest
		responsemodel.Data = gin.H{"error": "Invalid credentials format"}
		c.JSON(http.StatusBadRequest, responsemodel)
		return
	}

	email, password := credentials[0], credentials[1]

	if len(credentials) < 2 {
		log.Println("Invalid credentials length")
		responsemodel.Status = http.StatusBadRequest
		responsemodel.Data = gin.H{"error": "Invalid input format"}
		c.JSON(http.StatusBadRequest, responsemodel)
		return
	}

	// Debugging input
	log.Printf("Received email: %s", email)
	log.Printf("Received password: %s", password)
	log.Printf("Content-Type: %s", c.GetHeader("Content-Type"))

	if email == "" {
		responsemodel.Status = http.StatusBadRequest
		responsemodel.Data = gin.H{"error": "email must be fill"}
		c.JSON(http.StatusBadRequest, responsemodel)
		return
	}

	// Mencari user di database berdasarkan email
	var user models.DataNewUser
	if err := models.DB.Where("email = ?", email).First(&user).Error; err != nil {
		fmt.Println("error connected db", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			responsemodel.Status = http.StatusUnauthorized
			responsemodel.Data = gin.H{"error": "Email or password is incorrect"}
		} else {
			log.Printf("Database error: %v", err)
			responsemodel.Status = http.StatusInternalServerError
			responsemodel.Data = gin.H{"error": "An error occurred while fetching user data"}
		}
		responsemodel.Status = http.StatusInternalServerError
		responsemodel.Data = gin.H{"error": "An error occurred while fetching user data"}
		c.JSON(responsemodel.Status, responsemodel)
		return
	}

	// Memverifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println("error compare password", err)
		responsemodel.Status = http.StatusUnauthorized
		responsemodel.Data = gin.H{"error": "Email or password is incorrect"}
		HandleFailedAttempts(email)
		c.JSON(http.StatusUnauthorized, responsemodel)
		return
	}

	// Jika login berhasil
	responsemodel.Status = http.StatusOK
	responsemodel.Data = gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	}
	c.JSON(http.StatusOK, responsemodel)

}

func HandleFailedAttempts(email string) {
	var LoginAttempts models.User
	//cari data user di db
	err := models.DB.First(&LoginAttempts, "email = ?", email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//jika user tidak ada,tambahkan baru ke database
			LoginAttempts = models.User{
				Email: email,
				Count: 1,
			}
		}
		models.DB.Create(&LoginAttempts)
		fmt.Printf("First failed attempt for %s\n", email)
		return
	} else {
		log.Fatalf("failed to query database: %v", err)
	}
	LoginAttempts.Count++

	//cek apakah akun harus dikunci
	if LoginAttempts.Count >= 3 {
		LoginAttempts.LockUntil = time.Now().Add(5 * time.Second)
		log.Printf("User %s locked out until %v", email, LoginAttempts.LockUntil)
	}
	models.DB.Save(&LoginAttempts)
}

func GetUser(c *gin.Context) {
	var responsemodel views.Baseresponse
	var responsdata models.UserAddressDetails

	iduser, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responsemodel.Status = http.StatusBadRequest
		responsemodel.Error = append(responsemodel.Error, "invalid user id")
		c.JSON(http.StatusBadRequest, responsemodel)
		return
	}
	result := models.DB.Find(&responsdata, iduser)
	if result.Error != nil {
		responsemodel.Status = http.StatusBadRequest
		responsemodel.Error = append(responsemodel.Error, "user not found")
		c.JSON(http.StatusBadRequest, responsemodel)
		return
	}

	responsemodel.Status = http.StatusOK
	responsemodel.Data = responsdata
	c.JSON(http.StatusOK, responsemodel)

}

func UpdateAddressUser(c *gin.Context) {
	var responsemodel views.Baseresponse
	var requestUser views.UserRegistRequest

	err := c.ShouldBindJSON(&requestUser)
	if err != nil {
		responsemodel.Status = http.StatusBadRequest
		responsemodel.Error = append(responsemodel.Error, err.Error())
		c.JSON(http.StatusBadRequest, responsemodel)
		return
	}
	UserId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responsemodel.Status = http.StatusBadRequest
		responsemodel.Error = append(responsemodel.Error, err.Error())
		c.JSON(http.StatusBadRequest, responsemodel)
		return
	}

	var Userdata models.UserAddressDetails
	result := models.DB.Find(&Userdata, UserId)
	if result.Error != nil {
		responsemodel.Status = http.StatusBadRequest
		responsemodel.Error = append(responsemodel.Error, result.Error.Error())
		c.JSON(http.StatusBadRequest, responsemodel)
		return
	}

	if requestUser.City != "" {
		Userdata.City = requestUser.City
	}

	if requestUser.Street != "" {
		Userdata.Street = requestUser.Street
	}

	if requestUser.Country != "" {
		Userdata.Country = requestUser.Country
	}
	if requestUser.Kodepos != "" {
		Userdata.Kodepos = requestUser.Kodepos
	}

	resultdata := models.DB.Save(&Userdata)
	if resultdata.Error != nil {
		responsemodel.Status = http.StatusBadRequest
		responsemodel.Error = append(responsemodel.Error, resultdata.Error.Error())
		return
	}
	responsemodel.Status = http.StatusOK
	responsemodel.Data = Userdata
	c.JSON(http.StatusOK, Userdata)
}
