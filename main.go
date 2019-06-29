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

	uri := os.Getenv("DATABASE_URL")
	db,err := sql.Open("postgres",uri)
	if err != nil {
		log.Fatal("fatal",err.Error())
	}
	defer db.Close()

	query := `INSERT INTO Customer (name,email,status) VALUES ($1,$2,$3) RETURNING id`
	var id int
	row := db.QueryRow(query ,t.Name,t.Email,t.Status)
	err = row.Scan(&id)
	if err != nil {
		log.Fatal("cant scan id",err.Error())
	}
	t.ID = id
	c.JSON(201, t)
}

func getCustomerHandler(c *gin.Context) {
	db,_ := sql.Open("postgres",os.Getenv("DATABASE_URL"))
	stmt,_ := db.Prepare("SELECT * FROM Customer")
	rows,_ := stmt.Query()

	todos := []Customer{}

	for rows.Next(){

		t := Customer{}
		err1 := rows.Scan(&t.ID,&t.Name,&t.Email ,&t.Status)
		if err1 != nil{
			log.Fatal("error scan", err1.Error())
		}
		todos = append(todos,t)
	}

	
	c.JSON(200, todos)
}

func getCustomerById(c *gin.Context) {
	t := Customer{}
	id := c.Param("id")
	db,_ := sql.Open("postgres",os.Getenv("DATABASE_URL"))
	stmt,err := db.Prepare("SELECT * FROM Customer WHERE id = $1")
	if err != nil{

	}

	row := stmt.QueryRow(id)

	err1 := row.Scan(&t.ID,&t.Name,&t.Email ,&t.Status)
		if err1 != nil{
			log.Fatal("error scan", err.Error())
		}

	c.JSON(http.StatusOK, t)
}

func deleteCustomerById(c *gin.Context){
	t := Customer{}
	id := c.Param("id")
	db,_ := sql.Open("postgres",os.Getenv("DATABASE_URL"))
	stmt := `DELETE FROM Customer WHERE id = $1`
	
	_,err := db.Exec(stmt,id)
	if err != nil{

	}

	//row.Scan(&t.ID,&t.Title, &t.Status)
	c.JSON(200, t)
}

func putCustomerById(c *gin.Context){

}

func checkDataBaseExist(){
	uri := os.Getenv("DATABASE_URL")
	db,err := sql.Open("postgres",uri)
	if err != nil {
		log.Fatal("fatal",err.Error())
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

type Customer struct {
	ID int
	Name string
	Email string
	Status string
}