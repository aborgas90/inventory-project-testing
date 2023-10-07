package services

import (
	"errors"

	"inventory-project-testing/database"
	"inventory-project-testing/models"
	"inventory-project-testing/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//Signup return JWT token for the user
func Signup(userInput models.UserRequest) (string, error){
	//create a password using bcrypt Library
	password, err := bcrypt.GenerateFromPassword([]byte(userInput.Password),bcrypt.DefaultCost)

	//if password creation failed, return the errror
	if err != nil {
		return "",err
	}

	//create a new user object
	//this user will be added into the database
	var user models.User = models.User{
		ID: uuid.New().String(),
		Email: userInput.Email,
		Password: string(password),
	} 

	//create a user into the database
	database.DB.Create(&user)

	//generate the JWT token 
	token, err := utils.GeneralNewAccessToken()

	//if generation is failed, return the error
	if err != nil {
		return "", err
	}

	//return the JWT token
	return token, nil
}


func Login(userInput models.UserRequest) (string, error){
	//create a variable called "user"
	var user models.User

	//find the user based on them email
	result := database.DB.First(&user, "email = ?", userInput.Email)

	//if the user is not found, return the error
	if result.RowsAffected == 0 {
		return "",errors.New("Invalid password")
	}

	//generate the JWT token 
	token, err := utils.GeneralNewAccessToken()

	//if gneration is failed, return the error 
	if err != nil {
		return "", err
	}

	//return the JWT token
	return token, nil
}
