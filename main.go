package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func main() {
	godotenv.Load()
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_DB")
	ipAdd := os.Getenv("DB_IP")
	serverAdd := os.Getenv("SERVER_ADDR")

	var err error
	db, err = sql.Open("postgres", "postgres://"+dbUser+":"+dbPass+"@"+ipAdd+"/"+dbName+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.GET("/customers", getCustomers)
	router.GET("/customers/:id", getCustomersById)
	router.POST("/customers", createCustomer)
	router.PATCH("/customers/:id", updateCustomer)

	router.Run(serverAdd)
}

type getResponse struct {
	ID        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"isActive"`
}

// Returns a list of customers from the database
func getCustomers(c *gin.Context) {
	rows, err := db.Query("SELECT id, firstName, lastName, createdAt, updatedAt, email, isActive FROM customers")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var customers []getResponse
	for rows.Next() {
		var a getResponse
		err := rows.Scan(&a.ID, &a.FirstName, &a.LastName, &a.CreatedAt, &a.UpdatedAt, &a.Email, &a.IsActive)
		if err != nil {
			log.Fatal(err)
		}
		customers = append(customers, a)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, customers)
}

func getCustomersById(c *gin.Context) {
	rows, err := db.Query("SELECT id, firstName, lastName, createdAt, updatedAt, email, isActive FROM customers WHERE id = $1", c.Param("id"))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var customer []getResponse
	for rows.Next() {
		var a getResponse
		err := rows.Scan(&a.ID, &a.FirstName, &a.LastName, &a.CreatedAt, &a.UpdatedAt, &a.Email, &a.IsActive)
		if err != nil {
			log.Fatal(err)
		}
		customer = append(customer, a)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, customer)
}

type postRequest struct {
	ID        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"isActive"`
}

// Creates a new customer on database
func createCustomer(c *gin.Context) {
	var awesomeCustomer postRequest
	if err := c.BindJSON(&awesomeCustomer); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	password := string(awesomeCustomer.Password)

	// Hashing the passowrd using Salt with SALT_SECRET env var
	saltSecret := os.Getenv("SALT_SECRET")
	salt := []byte(password + saltSecret)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(salt), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalln(err)
	}
	hashedPassword := string(hashedPasswordBytes)

	stmt, err := db.Prepare("INSERT INTO customers (firstName, lastName, createdAt, updatedAt, password, email, isActive) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(awesomeCustomer.FirstName, awesomeCustomer.LastName, awesomeCustomer.CreatedAt, awesomeCustomer.UpdatedAt, hashedPassword, awesomeCustomer.Email, awesomeCustomer.IsActive); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusCreated, gin.H{"Message": "Customer created successfully", "Customer Email": awesomeCustomer.Email})
}

// Updates an existing customer on database
func updateCustomer(c *gin.Context) {
	var awesomeCustomer postRequest
	if err := c.BindJSON(&awesomeCustomer); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	stmt, err := db.Prepare("UPDATE customers SET lastName=$1 WHERE id=$2;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(awesomeCustomer.LastName); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusCreated, gin.H{"Message": "Customer updated successfully"})
}
