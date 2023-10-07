package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

//GetValue return configuration value based on a given key from the .env file
func GetValue(key string) string {
	//load the .env file
	err := godotenv.Load(".env")
	if err != nil{
		log.Fatalf("Error loading .env file\n")
	}

	//return the value based on agiven key
	return os.Getenv(key)
}

/*
	In this directory, some helpers are created. 
	The first helper is for reading database credentials or configurations from the 
	.env file. Inside the utils directory, a new file called utils.go is created.
*/