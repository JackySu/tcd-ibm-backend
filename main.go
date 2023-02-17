package main

import (
	"fmt"

	"sweng_backend/database"
	"sweng_backend/middleware"
	"sweng_backend/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetAnnouncement(c *gin.Context) {
	var announcements []model.Announcement
	if err := DB.Find(&announcements).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, announcements)
	}
}

func GetPA(c *gin.Context) {
	var pas []model.PA
	if err := DB.Find(&pas).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, pas)
	}
}

func GetDetail(c *gin.Context) {
	var details []model.Detail
	if err := DB.Find(&details).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, details)
	}
}

func GetProduct(c *gin.Context) {
	var products []model.Product
	if err := DB.Find(&products).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, products)
	}
}

func GetSolution(c *gin.Context) {
	var solutions []model.Solution
	if err := DB.Find(&solutions).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, solutions)
	}
}

func GetVertical(c *gin.Context) {
	var verticals []model.Vertical
	if err := DB.Find(&verticals).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, verticals)
	}
}

func GetType(c *gin.Context) {
	var types []model.Type
	if err := DB.Find(&types).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, types)
	}
}

func main() {

	DB = database.InitDB()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.CORS())

	api := r.Group("/api/v1")

	api.GET("/announcement", GetAnnouncement)
	api.GET("/pa", GetPA)
	api.GET("/detail", GetDetail)
	api.GET("/product", GetProduct)
	api.GET("/solution", GetSolution)
	api.GET("/vertical", GetVertical)
	api.GET("/type", GetType)

	auth := api.Group("/auth")
	auth.POST("/register", middleware.RegisterHandler)
	auth.POST("/login", middleware.LoginHandler)
	auth.POST("/refresh_token", middleware.RefreshHandler)
	auth.GET("/info", middleware.AuthMiddleware(), middleware.InfoHandler)

	/*
		// The 2 snippets below are equivalent
		id := 1
		var c model.Category

		DB.Find(&c, id)                                // SELECT * FROM categories WHERE id = 1;
		DB.Model(&c).Association("Tags").Find(&c.Tags) // fill up the Tags field
		fmt.Println(c)

		c = model.Category{}
		DB.Model(&model.Category{}).Preload("Tags").Find(&c, id) // Preload the Tags field and fill up all fields
		fmt.Println(c)
	*/

	r.Run("0.0.0.0:5297")
}
