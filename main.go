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
	router.PUT("/customers/:id", updateCustomer)
	router.DELETE("/customers/:id", removeCustomersById)

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
	rows, err := db.Query("SELECT id, first_name, last_name, created_at, updated_at, email, is_active FROM customers WHERE is_active=true")
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

// Returns a single customer from the database
func getCustomersById(c *gin.Context) {
	rows, err := db.Query("SELECT id, first_name, last_name, created_at, updated_at, email, is_active FROM customers WHERE id = $1 AND is_active=true", c.Param("id"))
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
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	Email     string `json:"email"`
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

	stmt, err := db.Exec("INSERT INTO customers (first_name, last_name, created_at, updated_at, password, email, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7)", &awesomeCustomer.FirstName, &awesomeCustomer.LastName, time.Now(), time.Now(), hashedPassword, &awesomeCustomer.Email, true)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.RowsAffected()

	c.JSON(http.StatusCreated, gin.H{"Message": "Customer created successfully"})
}

type putRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	IsActive  bool   `json:"isActive"`
}

// Updates an existing customer on database
func updateCustomer(c *gin.Context) {
	var awesomeCustomer putRequest
	if err := c.BindJSON(&awesomeCustomer); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	stmt, err := db.Exec(`UPDATE customers SET first_name=$1, last_name=$2, email=$3, updated_at=$4 WHERE id=$5`, &awesomeCustomer.FirstName, &awesomeCustomer.LastName, &awesomeCustomer.Email, time.Now(), c.Param("id"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.RowsAffected()

	c.JSON(http.StatusOK, gin.H{"Message": "Customer updated successfully"})
}

func removeCustomersById(c *gin.Context) {
	stmt, err := db.Exec("UPDATE customers SET is_active=false WHERE id=$1", c.Param("id"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.RowsAffected()

	c.JSON(http.StatusNoContent, nil)
}
