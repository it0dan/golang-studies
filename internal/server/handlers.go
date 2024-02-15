package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/it0dan/golang-studies/internal/models"
	"github.com/it0dan/golang-studies/pkg/db/postgres"
)

// Returns a list of customers from the database
func GetCustomers(c *gin.Context) {
	rows, err := postgres.Db.Query("SELECT id, first_name, last_name, created_at, updated_at, email, is_active FROM customers WHERE is_active=true")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var customers []models.GetCustomer
	for rows.Next() {
		var a models.GetCustomer
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
func GetCustomersById(c *gin.Context) {
	rows, err := postgres.Db.Query("SELECT id, first_name, last_name, created_at, updated_at, email, is_active FROM customers WHERE id = $1 AND is_active=true", c.Param("id"))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var customer []models.GetCustomer
	for rows.Next() {
		var a models.GetCustomer
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

// Creates a new customer on database
func CreateCustomer(c *gin.Context) {
	godotenv.Load()
	var awesomeCustomer models.Customer
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

	stmt, err := postgres.Db.Exec("INSERT INTO customers (first_name, last_name, created_at, updated_at, password, email, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", &awesomeCustomer.FirstName, &awesomeCustomer.LastName, time.Now(), time.Now(), hashedPassword, &awesomeCustomer.Email, true)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.RowsAffected()

	// c.JSON(http.StatusCreated, gin.H{"Message": "Customer created successfully"})
	c.PureJSON(http.StatusCreated, gin.H{"Message": "Customer created successfully", "email": &awesomeCustomer.Email})
}

// Updates an existing customer on database
func UpdateCustomer(c *gin.Context) {
	var awesomeCustomer models.Customer
	if err := c.BindJSON(&awesomeCustomer); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	stmt, err := postgres.Db.Exec(`UPDATE customers SET first_name=$1, last_name=$2, email=$3, updated_at=$4 WHERE id=$5`, &awesomeCustomer.FirstName, &awesomeCustomer.LastName, &awesomeCustomer.Email, time.Now(), c.Param("id"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.RowsAffected()

	c.JSON(http.StatusOK, gin.H{"Message": "Customer updated successfully"})
}

func RemoveCustomersById(c *gin.Context) {
	stmt, err := postgres.Db.Exec("UPDATE customers SET is_active=false WHERE id=$1", c.Param("id"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.RowsAffected()

	c.JSON(http.StatusNoContent, nil)
}
