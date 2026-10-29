package controller

import (
	"context"
	"log"
	"net/http"
	"resturant-management-system/database"
	"resturant-management-system/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := menuCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing the menu item"})
		}
		var allMenus []bson.M
		if err = result.All(ctx, &allMenus); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenus)
	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		menuId := c.Param("menu_id")
		var menu models.Menu

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "menu not found"})
			return
		}

		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(menu)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		if menu.Start_Date != nil && menu.End_Date != nil && !menu.Start_Date.Before(*menu.End_Date) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "start_date must be before end_date"})
			return
		}

		now := time.Now().Truncate(time.Second)
		menu.Created_at = now
		menu.Updated_at = now
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result, insertErr := menuCollection.InsertOne(ctx, menu)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "menu item was not created"})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		menuId := c.Param("menu_id")
		var menu models.Menu

		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if menu.Start_Date != nil && menu.End_Date != nil && !menu.Start_Date.Before(*menu.End_Date) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "start_date must be before end_date"})
			return
		}

		var updateObj bson.D

		if menu.Name != nil {
			updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
		}
		if menu.Category != nil {
			updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Category})
		}
		if menu.Start_Date != nil {
			updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.Start_Date})
		}
		if menu.End_Date != nil {
			updateObj = append(updateObj, bson.E{Key: "end_date", Value: menu.End_Date})
		}

		if len(updateObj) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
			return
		}

		menu.Updated_at = time.Now().Truncate(time.Second)
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.Updated_at})

		result, err := menuCollection.UpdateOne(
			ctx,
			bson.M{"menu_id": menuId},
			bson.D{{Key: "$set", Value: updateObj}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "menu update failed"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "menu not found"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
