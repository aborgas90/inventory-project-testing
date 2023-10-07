package main

import (
	"fmt"
	"os"

	"inventory-project-testing/database"
	"inventory-project-testing/routes"
	"inventory-project-testing/utils"

	"github.com/gofiber/fiber/v2"
)

// define the default port of the application
const DEFAULT_PORT = "3000"

// NewFiberApp return fiber application
func NewFiberApp() *fiber.App {
	//create a new fiber application
	var app *fiber.App = fiber.New()

	//define the routes
	routes.SetupRoutes(app)

	return app
}

func main() {
	//create a new fiber application
	var app *fiber.App = NewFiberApp()

	//connect to DB
	database.InitDatabase(utils.GetValue("DB_NAME"))

	//get the application port from the defined PORT variable
	var PORT string = os.Getenv("PORT")

	//if the PORT variable is not assigned
	//used the default port
	if PORT == "" {
		PORT = DEFAULT_PORT
	}

	//start the application
	app.Listen(fmt.Sprintf(":%s", PORT))
}

/*if u want to used to test rest api in GitBash u can use command :
curl -XPOST -H "Content-type: application/json" -d '{"name":"milk","price":299,"quantity":100}' 'http://127.0.0.1:3000/api/v1/items'
*/


/*
	Try to run the REST API application. If the application is running, try to send a request to the application with curl by clicking the “+” button to open a new terminal.
	Request Example
	1. Sign Up
		curl -XPOST -H "Content-type: application/json" -d '{"email":"test@test.com","password":"123123"}' 'http://127.0.0.1:3306/api/v1/signup'
	2. Create a new item with authentication
		curl -XPOST -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTQ2MzMxNjV9.6uIxGUod4aQJTobb6a0hHniX-VqPhBvfEo7B-E627Jo' -H "Content-type: application/json" -d '{"name":"coffee","price":10,"quantity":10}' 'http://127.0.0.1:3306/api/v1/items'
		checking token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTI2MzkwMDh9.Nzcz0X9DzWZ-_WeqnPy8bXpZY94Qs60VykezoL2Q9co
			curl -XPOST -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTQ2MzI5NTB9.eTmnNqcyZWRSIV_b05_FuObsaSK5j6smd34myLs-uMM' -H "Content-type: application/json" -d '{"name":"coffee","price":10,"quantity":10}' 'http://127.0.0.1:3306/api/v1/items'
		eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTI2Mzk0ODl9.O8MW04Oi9Q4f1V5bE4RM54cKTSQMqj0L8YENhI22OYg
			curl -XPOST -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTI2Mzk0ODl9.O8MW04Oi9Q4f1V5bE4RM54cKTSQMqj0L8YENhI22OYg' -H "Content-type: application/json" -d '{"name":"coffee","price":10,"quantity":10}' 'http://127.0.0.1:3306/api/v1/items'
			3. Get all items
		curl -XGET -H "Content-type: application/json" 'http://127.0.0.1:3306/api/v1/items'
*/