package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type Recipe struct {
	//swagger:ignore
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

var recipes []Recipe

func init() {

	recipes = make([]Recipe, 0)

	file, _ := ioutil.ReadFile("recipes.json")

	_ = json.Unmarshal([]byte(file), &recipes)

}

func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

func IndexHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello world",
	})
}

func ListRecipesHandler(c *gin.Context) {

	c.JSON(http.StatusOK, recipes)

}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Recipe not found"})
	}
	recipes[index] = recipe
	c.JSON(http.StatusOK, recipe)

}

func GetRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			c.JSON(http.StatusOK, recipes[i])
		}
		return
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "Recipe not found"})
}

func main() {
	router := gin.Default()
	router.GET("/", IndexHandler)
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.GET("/recipes/:id", GetRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	http.ListenAndServe(":8080", router)
}
