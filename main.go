package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"database/sql"
	"log"
	"os"
	_"github.com/lib/pq"
)

func main() {
	r := setupRouter()
	r.Run(getPort())
}

func setupRouter() *gin.Engine{
	r := gin.Default()
	r.Use(middleware)
	r.GET("/api/todos", getTodosHandler)
	r.GET("/api/todos/:id", getTodoById)
	r.POST("/api/todos", postTodoHandler)
	r.DELETE("/api/todos/:id", deleteTodoById)
	return r
}

func middleware(c *gin.Context){
	fmt.Println("malaew")
	token := c.GetHeader("Authorization")
	fmt.Println("token:", token)
	if token != "Bearer token 123" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
		c.Abort()
		return
	}
	c.Next()
	fmt.Println("naja")
}

func getPort()string{
	return ":2019"
}

func getTodosHandler(c *gin.Context) {
	db,_ := sql.Open("postgres",os.Getenv("DATABASE_URL"))
	stmt,_ := db.Prepare("SELECT id,title,status FROM todos")
	rows,_ := stmt.Query()

	todos := []Todo{}

	for rows.Next(){

		t := Todo{}
		err1 := rows.Scan(&t.ID,&t.Title, &t.Status)
		if err1 != nil{
			log.Fatal("error scan", err1.Error())
		}
		todos = append(todos,t)
	}

	
	c.JSON(200, todos)
}

func getTodoById(c *gin.Context) {
	t := Todo{}
	id := c.Param("id")
	db,_ := sql.Open("postgres",os.Getenv("DATABASE_URL"))
	stmt,err := db.Prepare("SELECT id,title,status FROM todos WHERE id = $1")
	if err != nil{

	}

	row := stmt.QueryRow(id)

	err1 := row.Scan(&t.ID,&t.Title, &t.Status)
		if err1 != nil{
			log.Fatal("error scan", err.Error())
		}

	c.JSON(http.StatusOK, t)
}

func postTodoHandler(c *gin.Context) {
	t := Todo{}
	
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

	query := `INSERT INTO todos (title,status) VALUES ($1,$2) RETURNING id`
	var id int
	row := db.QueryRow(query ,t.Title,t.Status)
	err = row.Scan(&id)
	if err != nil {
		log.Fatal("cant scan id",err.Error())
	}

	c.JSON(201, t)
}

func deleteTodoById(c *gin.Context){
	t := Todo{}
	id := c.Param("id")
	db,_ := sql.Open("postgres",os.Getenv("DATABASE_URL"))
	stmt := `DELETE FROM todos WHERE id = $1`
	
	_,err := db.Exec(stmt,id)
	if err != nil{

	}

	//row.Scan(&t.ID,&t.Title, &t.Status)
	c.JSON(200, t)
}

type Todo struct {
	ID int
	Title string
	Status string
}