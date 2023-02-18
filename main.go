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
	var announcement model.Announcement
	if err := DB.First(&announcement, c.Param("id")).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, announcement)
	}
}

func GetAnnouncements(c *gin.Context) {
	var announcements []model.Announcement
	if err := DB.Find(&announcements).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, announcements)
	}
}

func GetPA(c *gin.Context) {
	var pa model.PA
	if err := DB.First(&pa, c.Param("id")).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, pa)
	}
}

func GetPAs(c *gin.Context) {
	var pas []model.PA
	if err := DB.Find(&pas).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, pas)
	}
}

func GetDetail(c *gin.Context) {
	var detail model.Detail
	if err := DB.First(&detail, c.Param("id")).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, detail)
	}
}

func GetDetails(c *gin.Context) {
	var details []model.Detail
	if err := DB.Find(&details).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, details)
	}
}

func GetProducts(c *gin.Context) {
	var products []model.Product
	if err := DB.Find(&products).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, products)
	}
}

func GetSolutions(c *gin.Context) {
	var solutions []model.Solution
	if err := DB.Find(&solutions).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, solutions)
	}
}

func GetVerticals(c *gin.Context) {
	var verticals []model.Vertical
	if err := DB.Find(&verticals).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, verticals)
	}
}

func GetTypes(c *gin.Context) {
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

	api.GET("/announcement", GetAnnouncements)
	api.GET("/announcement/:id", GetAnnouncement)

	api.GET("/pa", GetPAs)
	api.GET("/pa/:id", GetPA)

	api.GET("/detail", GetDetails)
	api.GET("/detail/:id", GetDetail)

	api.GET("/product", GetProducts)

	api.GET("/solution", GetSolutions)

	api.GET("/vertical", GetVerticals)

	api.GET("/type", GetTypes)

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
