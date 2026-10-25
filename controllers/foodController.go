package controller

import (
	"context"
	"net/http"
	"resturant-management-system/database"
	"resturant-management-system/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second) // if the query takes longer than 100 seconds, cancel it
		foodId := c.Param("food_id")                                                 //read id from url
		var food models.Food                                                         // create empty food struct in memory. to be filled
		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)  //db call
		// find the document where food_id equals foodId”. bson.M is just MongoDB’s map type, like a JS object { food_id: foodId }
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the food item"})
		}
		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func round(num float64) int {}

func toFixed(num float64, precision int) int {}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
