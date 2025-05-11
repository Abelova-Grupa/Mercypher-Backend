package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type person struct {
	ID		int32	`json:"id"`
	Name 	string	`json:"name"`
	Age		int		`json:"age"`
	Married	bool	`json:"married"`
}

// Random data from ChatGPT
var people = []person {
	{ID: 1, Name: "Alice", Age: 30, Married: true},
	{ID: 2, Name: "Bob", Age: 25, Married: false},
	{ID: 3, Name: "Charlie", Age: 35, Married: true},
	{ID: 4, Name: "Diana", Age: 28, Married: false},
	{ID: 5, Name: "Evan", Age: 40, Married: true},
	{ID: 6, Name: "Fiona", Age: 31, Married: true},
	{ID: 7, Name: "George", Age: 26, Married: false},
	{ID: 8, Name: "Hannah", Age: 29, Married: false},
	{ID: 9, Name: "Ian", Age: 34, Married: true},
	{ID: 10, Name: "Jane", Age: 27, Married: false},
	{ID: 11, Name: "Karl", Age: 38, Married: true},
	{ID: 12, Name: "Laura", Age: 24, Married: false},
	{ID: 13, Name: "Mike", Age: 33, Married: true},
	{ID: 14, Name: "Nina", Age: 30, Married: false},
	{ID: 15, Name: "Oscar", Age: 41, Married: true},
	{ID: 16, Name: "Paula", Age: 36, Married: true},
	{ID: 17, Name: "Quentin", Age: 22, Married: false},
	{ID: 18, Name: "Rachel", Age: 32, Married: true},
	{ID: 19, Name: "Steve", Age: 37, Married: true},
	{ID: 20, Name: "Tina", Age: 23, Married: false},
	{ID: 21, Name: "Umar", Age: 29, Married: false},
	{ID: 22, Name: "Violet", Age: 39, Married: true},
	{ID: 23, Name: "Will", Age: 28, Married: false},
	{ID: 24, Name: "Xena", Age: 35, Married: true},
	{ID: 25, Name: "Yusuf", Age: 27, Married: false},
	{ID: 26, Name: "Zara", Age: 34, Married: true},
}

// Handlers

func getPeople(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, people)
}

func addPerson(c *gin.Context) {
	var newPerson person

	if err := c.BindJSON(&newPerson); err != nil {return}

	people = append(people, newPerson)
	c.IndentedJSON(http.StatusCreated, newPerson)
}

// EXAMPLE POST

// curl -X POST http://localhost:8080/people \
//   -H "Content-Type: application/json" \
//   -d '{
//     "id": 31,
//     "name": "Elena",
//     "age": 29,
//     "married": false
//   }'

func main() {

	router := gin.Default()

	router.GET("/people", getPeople)
	router.POST("/people", addPerson)

	router.Run("localhost:8080")
}