package database

import (
	"fmt"
	"errors"
	
	"inventory-project-testing/models"
	"inventory-project-testing/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// create a variable to store the database instance
var DB *gorm.DB

// InitDatabase creates a connection to the database
func InitDatabase(dbName string) {

    // initialize some variables
    // for the MySQL data source
	var (
		databaseUser     string = utils.GetValue("DB_USER")
		databasePassword string = utils.GetValue("DB_PASSWORD")
		databaseHost     string = utils.GetValue("DB_HOST")
		databasePort     string = utils.GetValue("DB_PORT")
		databaseName     string = dbName
	)

    // declare the data source for MySQL
	var dataSource string = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", databaseUser, databasePassword, databaseHost, databasePort, databaseName)

    // create a variable to store an error
	var err error

    // create a connection to the database
	DB, err = gorm.Open(mysql.Open(dataSource), &gorm.Config{})

    // if connection fails, print out the errors
	if err != nil {
		panic(err.Error())
	}

    // if connection is successful, print out this message
	fmt.Println("Connected to the database")

	DB.AutoMigrate(&models.User{}, &models.Item{})
}


//SeedItem returns recently created items from the database
func SeedItem()(models.Item, error){
	//create a sample data for item 
	item, err := utils.CreateFaker[models.Item]()
	if err != nil{
		return models.Item{}, nil
	}

	//insert the sample data into the database 
	DB.Create(&item)
	fmt.Println("Item seeded to the database")

	//return recently created item
	return item, nil
}


//SeedUser returns recently created user from the database 
func SeedUser()(models.User, error){
	//create a sample data for user
	user, err := utils.CreateFaker[models.User]()
	if err != nil{
		return models.User{}, err
	}

	//create a password with bcrypt 
	//this password is stored to the database
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil{
		return models.User{}, err
	}

	//create a variable called "inputUser"
	//this variable is used to store the user sample data
	//into the database

	var inputUser models.User = models.User{
		ID: user.ID,
		Email: user.Email,
		Password: string(password),
	}

	//insert the user sample data into the database 
	DB.Create(&inputUser)
	fmt.Println("User seeded to the database")

	//return the user sample data
	return user,nil
}

// CleanSeeders performs clean up mechanism after testing
func CleanSeeders() {
    // remove all data inside items table
    itemResult := DB.Exec("TRUNCATE items")
    // remove all data inside users table
    userResult := DB.Exec("TRUNCATE users")


    // check if the operation is failed
    var isFailed bool = itemResult.Error != nil || userResult.Error != nil


    // if operation is failed, return an error
    if isFailed {
        panic(errors.New("error when cleaning up seeders"))
    }


    fmt.Println("Seeders are cleaned up successfully")
}