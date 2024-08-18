package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// use struct tags to conform to the json naming convention
type Page struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (p *Page) save() error {
	filename := "pages/" + p.Title + ".txt"
	return os.WriteFile(filename, []byte(p.Body), 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "pages/" + title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: string(body)}, nil
}

func getPage(c *gin.Context) {
	title := c.Param("title")

	p, err := loadPage(title)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "non-existing page"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func savePage(c *gin.Context) {
	var newPage Page

	if err := c.BindJSON(&newPage); err != nil {
		fmt.Printf("%v", err)
		return
	}

	newPage.save()
	c.JSON(http.StatusCreated, newPage)
}

func main() {
	p1 := &Page{Title: "TestPage", Body: "This is a sample Page."}
	p1.save()
	p2, err := loadPage("TestPage")
	if err != nil {
		fmt.Printf("failed to create the example page: %e", err)
		return
	}
	fmt.Println(p2.Body)

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/view/:title", getPage)
	router.POST("/save", savePage)

	router.Run("localhost:30001")
}
