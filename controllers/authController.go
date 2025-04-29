package controllers

import (
	"net/http"
	"os"
	"ride-sharing/initializers"
	"ride-sharing/models"
	"ride-sharing/schemas"
	"ride-sharing/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {

	var authInput schemas.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		utils.Error(c.Writer, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	var userFound models.User
	initializers.DB.Where("email=?", authInput.Email).Find(&userFound)

	if userFound.ID != uuid.Nil {
		utils.Error(c.Writer, http.StatusBadRequest, "Email already used", nil)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c.Writer, http.StatusInternalServerError, "Failed to hash password", err.Error())
		return
	}

	user := models.User{
		Email:    authInput.Email,
		Password: string(passwordHash),
	}

	initializers.DB.Create(&user)

	utils.Success(c.Writer, http.StatusCreated, "User created successfully", user, "")

}

func Login(c *gin.Context) {
	var authInput schemas.AuthInput

	if err := c.ShouldBindJSON(&authInput); err != nil {
		utils.Error(c.Writer, http.StatusBadRequest, utils.ErrInvalidInput, err.Error())
		return
	}

	var userFound models.User
	initializers.DB.Where("email=?", authInput.Email).Find(&userFound)

	if userFound.ID == uuid.Nil {
		utils.Error(c.Writer, http.StatusNotFound, "User not found", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password)); err != nil {
		utils.Error(c.Writer, http.StatusUnauthorized, "Invalid password", err.Error())
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		utils.Error(c.Writer, http.StatusInternalServerError, "Failed to generate token", err.Error())
		return
	}

	// Return token as success
	utils.Success(c.Writer, http.StatusOK, utils.AUTH_SUCCESS_LOGIN, gin.H{
		"token": signedToken,
	}, "")
}

func GetUserProfile(c *gin.Context) {

	user, _ := c.Get("currentUser")

	c.JSON(200, gin.H{
		"user": user,
	})
}
