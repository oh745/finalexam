package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"database/sql"
	"log"
	"os"
	_"github.com/lib/pq"
)

func main() {
	checkDataBaseExist()
	r := setupRouter()
	r.Run(getPort())
}

func setupRouter() *gin.Engine{
	r := gin.Default()
	//r.Use(middleware)
	r.POST("/customers", postCustomerHandler)
	r.GET("/customers", getCustomerHandler)
	r.GET("/customers/:id", getCustomerById)
	r.PUT("/customers/:id",putCustomerById)
	r.DELETE("/customers/:id", deleteCustomerById)
	return r
}

// func middleware(c *gin.Context){
// 	fmt.Println("malaew")
// 	token := c.GetHeader("Authorization")
// 	fmt.Println("token:", token)
// 	if token != "Bearer token 123" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
// 		c.Abort()
// 		return
// 	}
// 	c.Next()
// 	fmt.Println("naja")
// }

func getPort()string{
	return ":2019"
}

func postCustomerHandler(c *gin.Context) {
	t := Customer{}
	
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	db := connectDB(c)
	defer db.Close()

	query := `INSERT INTO Customer (name,email,status) VALUES ($1,$2,$3) RETURNING id`
	var id int
	row := db.QueryRow(query ,t.Name,t.Email,t.Status)
	err := row.Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	t.ID = id
	c.JSON(201, t)
}

func getCustomerHandler(c *gin.Context) {
	db := connectDB(c)
	stmt,_ := db.Prepare("SELECT * FROM Customer")
	rows,_ := stmt.Query()

	todos := []Customer{}

	for rows.Next(){

		t := Customer{}
		err1 := rows.Scan(&t.ID,&t.Name,&t.Email ,&t.Status)
		if err1 != nil{
			c.JSON(http.StatusInternalServerError, err1)
			return
		}
		todos = append(todos,t)
	}

	
	c.JSON(200, todos)
}

func getCustomerById(c *gin.Context) {
	t := Customer{}
	id := c.Param("id")
	db := connectDB(c)
	stmt,err := db.Prepare("SELECT * FROM Customer WHERE id = $1")
	if err != nil{
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	row := stmt.QueryRow(id)

	err1 := row.Scan(&t.ID,&t.Name,&t.Email ,&t.Status)
		if err1 != nil{
			c.JSON(http.StatusInternalServerError, err)
			return
		}

	c.JSON(http.StatusOK, t)
}

func deleteCustomerById(c *gin.Context){
	id := c.Param("id")
	db := connectDB(c)
	defer db.Close()
	stmt := `DELETE FROM Customer WHERE id = $1`
	
	_,err := db.Exec(stmt,id)
	if err != nil{
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, gin.H{
		"message": "customer deleted",
	})
}

func putCustomerById(c *gin.Context){
	t := Customer{}
	
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	db := connectDB(c)
	defer db.Close()

	stmt,err := db.Prepare("UPDATE Customer SET name=$2,email=$3,status=$4 WHERE id=$1")
	if err != nil{
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	row := stmt.QueryRow(t.ID,t.Name,t.Email,t.Status)
	err1 := row.Scan(&t.ID,&t.Name,&t.Email ,&t.Status)
		if err1 != nil{
			c.JSON(http.StatusInternalServerError, err)
			return
		}

	c.JSON(http.StatusOK, t)
}

func checkDataBaseExist(){
	db,err := sql.Open("postgres",os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("cant connect db",err.Error())
	}
	defer db.Close()

	createTb := 
	`CREATE TABLE IF NOT EXISTS CUSTOMER(
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);`
	_,err = db.Exec(createTb)
	if err != nil {
		log.Fatal("cant create table",err.Error())
	}
}

func connectDB(c *gin.Context) *sql.DB {
	uri := os.Getenv("DATABASE_URL")
	db,err := sql.Open("postgres",uri)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return nil
	}
	return db
}

type Customer struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Status string `json:"status"`
}