package main

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"inventory-project-testing/database"
	"inventory-project-testing/models"
	"inventory-project-testing/utils"
	"github.com/steinfletcher/apitest"
)


// newApp returns application
func newApp() *fiber.App {
    // create a new application
    var app *fiber.App = NewFiberApp()


    // connect to the testing database
    database.InitDatabase(utils.GetValue("DB_NAME"))

    // return the application
    return app
}

// getItem returns sample data for item entity
func getItem() models.Item {
    // connect to the test database
    database.InitDatabase(utils.GetValue("DB_NAME"))


    // seed the item data to the test database
    // the sample data is stored in "item" variable
    item, err := database.SeedItem()
    if err != nil {
        panic(err)
    }


    // return the sample data for item entity
    return item
}

// cleanup performs clean up operation after the testing is completed
func cleanup(res *http.Response, req *http.Request, apiTest *apitest.APITest) {
    if http.StatusOK == res.StatusCode || http.StatusCreated == res.StatusCode {
        database.CleanSeeders()
    }
}

func TestSignup_Success(t *testing.T) {
    // create a sample data for user
    userData, err := utils.CreateFaker[models.User]()


    if err != nil {
        panic(err)
    }


    // create a sign up request
    // the request values are filled
    // with the value from the sample data
    var userRequest *models.UserRequest = &models.UserRequest{
        Email:    userData.Email,
        Password: userData.Password,
    }


    // create a test
    apitest.New().
        // Debug().
        // run the cleanup() function after the test is finished
        Observe(cleanup).
        // add an application to be tested
        HandlerFunc(FiberToHandlerFunc(newApp())).
        // send a POST request for sign up
        Post("/api/v1/signup").
        // set the request body
        JSON(userRequest).
        // expect the response status code is equals200
        Expect(t).
        Status(http.StatusOK).
        End()
}

func TestSignup_ValidationFailed(t *testing.T) {
    // create an empty request
    var userRequest *models.UserRequest = &models.UserRequest{
        Email:    "",
        Password: "",
    }


    // create a test
    apitest.New().
        // add an application to be tested
        HandlerFunc(FiberToHandlerFunc(newApp())).
        // send a POST request for sign up
        Post("/api/v1/signup").
        // set the request body
        JSON(userRequest).
        // expect the response status code is equals 400
        Expect(t).
        Status(http.StatusBadRequest).
        End()
}

func TestLogin_Success(t *testing.T) {
    // connect to the test database
    database.InitDatabase(utils.GetValue("DB_NAME"))


    // seed the sample data for user entity
    // after the data is seeded,
    // the sample data is returned to the "user" variable
    user, err := database.SeedUser()
    if err != nil {
        panic(err)
    }


    // create a request body for login
    var userRequest *models.UserRequest = &models.UserRequest{
        Email:    user.Email,
        Password: user.Password,
    }


    // create a test
    apitest.New().
        // run the cleanup() function after the test is finished
        Observe(cleanup).
         // add an application to be tested
        HandlerFunc(FiberToHandlerFunc(newApp())).
        // send a POST request for login
        Post("/api/v1/login").
        // set the request body
        JSON(userRequest).
        // expect the response status code is equals 200
        Expect(t).
        Status(http.StatusOK).
        End()
}

func TestLogin_Failed(t *testing.T) {
    // create a request body
    var userRequest *models.UserRequest = &models.UserRequest{
        Email:    "notfound@mail.com",
        Password: "123123",
    }


    // create a test
    apitest.New().
        // add an application to be tested
        HandlerFunc(FiberToHandlerFunc(newApp())).
        // send a POST request for login
        Post("/api/v1/login").
        // set the request body
        JSON(userRequest).
        // expect the response status code is equals 500
        Expect(t).
        Status(http.StatusInternalServerError).
        End()
}

// FiberToHandlerFunc returns handler function from the Fiber application
func FiberToHandlerFunc(app *fiber.App) http.HandlerFunc {
    // return the handler function
    return func(w http.ResponseWriter, r *http.Request) {
        // test the application with fiber
        resp, err := app.Test(r)
        if err != nil {
            panic(err)
        }


        // copy the headers
        for k, vv := range resp.Header {
            for _, v := range vv {
                w.Header().Add(k, v)
            }
        }
        w.WriteHeader(resp.StatusCode)


        if _, err := io.Copy(w, resp.Body); err != nil {
            panic(err)
        }
    }
}

func TestGetItem_Success(t *testing.T) {
    // get the sample data for item entity
    var item models.Item = getItem()

    // create a test
    apitest.New().
        // run the cleanup() function after the test is finished
        Observe(cleanup).
        // add an application to be tested
        HandlerFunc(FiberToHandlerFunc(newApp())).
        // send a GET request to get item data by ID
        Get("/api/v1/items/" + item.ID).
        // expect the response status code is equals 200
        Expect(t).
        Status(http.StatusOK).
        End()
}

func TestGetItem_NotFound(t *testing.T) {
    apitest.New().
        // add an application to be tested
        HandlerFunc(FiberToHandlerFunc(newApp())).
        // send GET request to get the item data by ID
        Get("/api/v1/items/0").
        // expect the response status code is equals 404
        Expect(t).
        Status(http.StatusNotFound).
        End()
}

// getJWTToken returns bearer token with JWT
func getJWTToken(t *testing.T) string {
    // connect to the test database
    database.InitDatabase(utils.GetValue("DB_NAME"))
    // insert a sample data for user into the database
    // the inserted sample data is returned into the "user variable"
    user, err := database.SeedUser()
    if err != nil {
        panic(err)
    }

    // create a request for login
    var userRequest *models.UserRequest = &models.UserRequest{
        Email:    user.Email,
        Password: user.Password,
    }

    // get the response from the login request
    var resp *http.Response = apitest.New().
        HandlerFunc(FiberToHandlerFunc(newApp())).
        Post("/api/v1/login").
        JSON(userRequest).
        Expect(t).
        Status(http.StatusOK).
        End().Response

    // create a variable called "response"
    // to store the response body from the login request
    var response *models.Response[string] = &models.Response[string]{}

    // decode the response body into the "response" variable
    json.NewDecoder(resp.Body).Decode(&response)

    // get the JWT token
    var token string = response.Data

    // create a bearer token
    var JWT_TOKEN = "Bearer " + token

    // return the bearer token with JWT
    return JWT_TOKEN
}

func TestCreateItem_Success(t *testing.T) {
    // create a sample data for item
    itemData, err := utils.CreateFaker[models.Item]()
    if err != nil {
        panic(err)
    }

    // create a request body to create a new item
    // the request body is filled with the value
    // from the sample data
    var itemRequest *models.ItemRequest = &models.ItemRequest{
        Name:     itemData.Name,
        Price:    itemData.Price,
        Quantity: itemData.Quantity,
    }

    // get the JWT token for authentication
    var token string = getJWTToken(t)

    // create a test
    apitest.New().
        // run the cleanup() function after the test is finished
        Observe(cleanup).
        // add an application to be tested
        HandlerFunc(FiberToHandlerFunc(newApp())).
        // send a POST request for creating a new item
        Post("/api/v1/items").
        // attach the JWT token into Authorization header
        Header("Authorization", token).
        // set the request body
        JSON(itemRequest).
        // expect the response status code is equals 201
        Expect(t).
        Status(http.StatusCreated).
        End()
}

func TestCreateItem_ValidationFailed(t *testing.T) {
    // create an empty request
    var itemRequest *models.ItemRequest = &models.ItemRequest{
        Name:     "",
        Price:    0,
        Quantity: 0,
    }

    // get the JWT token for authentication
    var token string = getJWTToken(t)

    // create a test
    apitest.New().
        // Debug().
        // add an application to be tested
        HandlerFunc(FiberToHandlerFunc(newApp())).
        // send a POST request for creating a new item
        Post("/api/v1/items").
        // attach the JWT token into Authorization header
        Header("Authorization", token).
        // set the request body
        JSON(itemRequest).
        // expect the response status code is equals 400
        Expect(t).
        Status(http.StatusBadRequest).
        End()
}

func TestUpdateItem_Success(t *testing.T) {
    // get the sample data for item
    var item models.Item = getItem()

    // create a request body to update an item
    var itemRequest *models.ItemRequest = &models.ItemRequest{
        Name:     item.Name,
        Price:    item.Price,
        Quantity: item.Quantity,
    }

    // get the JWT token for authentication
    var token string = getJWTToken(t)

    // create a test
    apitest.New().
        // run the cleanup() function after the test is finished
        Observe(cleanup).
        // add an application to be tested
        HandlerFunc(FiberToHandlerFunc(newApp())).
        // send a PUT request for updating an item
        Put("/api/v1/items/"+item.ID).
        // attach the JWT token into Authorization header
        Header("Authorization", token).
        // set the request body
        JSON(itemRequest).
        // expect the response status code is equals 200
        Expect(t).
        Status(http.StatusOK).
        End()
}

func TestUpdateItem_Failed(t *testing.T) {
    // create a request
    var itemRequest *models.ItemRequest = &models.ItemRequest{
        Name:     "changed",
        Price:    10,
        Quantity: 10,
    }

    // get the JWT token for authentication
    var token string = getJWTToken(t)

    // create a test
    apitest.New().
        // add an application to be tested
        HandlerFunc(FiberToHandlerFunc(newApp())).
        // send a PUT request for updating an item
        Put("/api/v1/items/0").
        // attach the JWT token into Authorization header
        Header("Authorization", token).
        // set the request body
        JSON(itemRequest).
        // expect the response status code is equals 404
        Expect(t).
        Status(http.StatusNotFound).
        End()
}


func TestDeleteItem_Success(t *testing.T) {
    // create a sample data for item
    // this sample data will be deleted
    var item models.Item = getItem()

    // get the JWT token for authentication
    var token string = getJWTToken(t)

    // create a test
    apitest.New().
        // add an application to be tested
        HandlerFunc(FiberToHandlerFunc(newApp())).
        // send a DELETE request for deleting an item
        // the item ID is picked from the sample data
        Delete("/api/v1/items/"+item.ID).
        // attach the JWT token into Authorization header
        Header("Authorization", token).
        // expect the response status code is equals 200
        Expect(t).
        Status(http.StatusOK).
        End()
}