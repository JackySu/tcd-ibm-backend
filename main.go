package main

import (
	"fmt"
	"strconv"
	"strings"

	"sweng_backend/database"
	"sweng_backend/middleware"
	"sweng_backend/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

var DB *gorm.DB
var DEFAULT_PER_PAGE int = 30
var MAX_PER_PAGE int = 50
var DEFAULT_PAGE int = 1

func GetAllCategoriesWithTags(c *gin.Context) {
	var categoriesWithTags []model.Category
	if err := DB.Preload("Tags").Find(&categoriesWithTags).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, model.ToCategoryDtoList(categoriesWithTags))
	}
}

func CreateProject(c *gin.Context) {

	user, _ := c.Get("user")
	if user == nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "User not found"})
		return
	}

	var projectBase model.ProjectBase

	if err := c.ShouldBindJSON(&projectBase); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": err.Error()})
		return
	}

	var tags = make([]model.Tag, len(projectBase.Tags))
	for i, tagId := range projectBase.Tags {
		if err := DB.First(&tags[i], tagId).Error; err != nil {
			c.AbortWithStatusJSON(400, gin.H{"message": fmt.Sprintf("TagId %d not found", tagId)})
			return
		}
	}

	var project = model.Project{
		Email:       user.(model.User).Email,
		Title:       projectBase.Title,
		Link:        projectBase.Link,
		Description: projectBase.Description,
		Content:     projectBase.Content,
		IsLive:      false,
		User:        user.(model.User),
		Tags:        tags,
	}
	if err := DB.Create(&project).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, model.ToProjectDto(project))
	}
}

func GetFullProjectById(c *gin.Context) {
	id := c.Params.ByName("id")
	var project model.Project
	if err := DB.Preload("Tags").Preload("User").First(&project, id).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, model.ToProjectFullDto(project))
	}
}

func GetProjectsByTagsInCategory(c *gin.Context) {

	tagIdStr := c.DefaultQuery("tags", "")
	keyword := c.DefaultQuery("keyword", "")
	perPageStr := c.DefaultQuery("per_page", "")
	pageStr := c.DefaultQuery("page", "")

	var per_page, page int
	if perPageStr == "" {
		per_page = DEFAULT_PER_PAGE
	} else {
		per_page, _ = strconv.Atoi(perPageStr)
		if per_page > MAX_PER_PAGE {
			per_page = MAX_PER_PAGE
		} else if per_page < 1 {
			per_page = 1
		}
	}
	if pageStr == "" {
		page = DEFAULT_PAGE
	} else {
		page, _ = strconv.Atoi(pageStr)
		if page < 1 {
			page = 1
		}
	}

	if tagIdStr == "" {
		var projects []model.Project

		if err := DB.Model(&model.Project{}).Preload("User").Preload("Tags").Offset((page - 1) * per_page).Limit(per_page).Find(&projects).Error; err != nil {
			c.AbortWithStatus(404)
			fmt.Println(err)
		}

		c.JSON(200, model.ToProjectDtoList(projects))
	}

	tagIdStrSlice := strings.Split(tagIdStr, ",")

	tagIds := make([]int, len(tagIdStrSlice))
	var err error
	for i, s := range tagIdStrSlice {
		tagIds[i], err = strconv.Atoi(s)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"message": "Tag id " + s + " must be integer"})
			return
		}
	}

	var tagInstances []model.Tag
	if err := DB.Model(&model.Tag{}).Where("id IN ?", tagIds).Find(&tagInstances).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}

	var allTagIds []int
	for _, tag := range tagInstances {
		allTagIds = append(allTagIds, tag.Id)
	}

	for _, tagId := range tagIds {
		if !slices.Contains(allTagIds, tagId) {
			c.AbortWithStatusJSON(404, gin.H{"tagId not found": tagId})
			return
		}
	}

	var categoryCount int64
	DB.Model(&model.Category{}).Count(&categoryCount)

	var tagIdsByCategory = make([][]int, categoryCount)
	for i := range tagIdsByCategory {
		for _, tag := range tagInstances {
			if tag.CategoryId == int(i+1) {
				tagIdsByCategory[i] = append(tagIdsByCategory[i], tag.Id)
			}
		}
	}

	// if any tag given in category1 and any tag given in category2 and any tag given in category3 ...
	var conditions []string
	for i, tagList := range tagIdsByCategory {
		if len(tagList) > 0 {
			conditions = append(conditions, fmt.Sprintf("( EXISTS ( SELECT 1 FROM project_tags, tags WHERE projects.id = project_tags.project_id AND tags.Id = project_tags.tag_id AND tags.category_id = %d AND tags.Id IN (%s) ))", i+1, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(tagList)), ","), "[]")))
		}
	}
	// EXISTS ( SELECT 1 ... ) => any_ in SQLAlchemy

	var projects []model.Project

	if err := DB.Model(&model.Project{}).Preload("User").Preload("Tags").Where(strings.Join(conditions, " AND ")).Where("title like ?", "%"+keyword+"%").Offset((page - 1) * per_page).Limit(per_page).Find(&projects).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}

	c.JSON(200, model.ToProjectDtoList(projects))
}

func GetUserProjects(c *gin.Context) {

	user, _ := c.Get("user")
	if user == nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "User not found"})
		return
	}

	var projects []model.Project
	if err := DB.Preload("Tags").Where("user_id = ?", user.(model.User).ID).Find(&projects).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, model.ToProjectDtoList(projects))
	}

}

func GetProjectById(c *gin.Context) {

	user, _ := c.Get("user")
	if user == nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "User not found"})
		return
	}
	id := c.Param("id")
	var project model.Project
	if err := DB.Preload("Tags").Preload("User").Where("user_id = ?", user.(model.User).ID).First(&project, id).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, model.ToProjectDto(project))
	}

}

func UpdateProject(c *gin.Context) {

	user, _ := c.Get("user")
	if user == nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "User not found"})
		return
	}
	id := c.Param("id")
	var updateData model.UpdateProject
	var projectToUpdate model.Project
	if err := DB.Preload("Tags").Preload("User").First(&projectToUpdate, id).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Invalid request"})
		return
	}

	if updateData.IsLive != nil {
		if !middleware.IsAdmin(user.(model.User)) {
			c.AbortWithStatusJSON(400, gin.H{"message": "Only admin can change isLive"})
			return
		}
		projectToUpdate.IsLive = *updateData.IsLive
	}

	if updateData.Tags != nil {
		var tagInstances = make([]model.Tag, len(*updateData.Tags))
		for i, tagId := range *updateData.Tags {
			if err := DB.First(&tagInstances[i], tagId).Error; err != nil {
				c.AbortWithStatusJSON(400, gin.H{"message": fmt.Sprintf("TagId %d not found", tagId)})
				return
			}
		}
		projectToUpdate.Tags = tagInstances
	}

	// []int is not serializable to JSON, convert other key-value pairs to []interface{}
	if err := DB.Model(&projectToUpdate).Omit("tags").Updates(map[string]interface{}{"Title": updateData.Title, "Description": updateData.Description, "Link": updateData.Link, "Content": updateData.Content}).Error; err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	DB.Save(&projectToUpdate)
	c.JSON(200, model.ToProjectDto(projectToUpdate))
}

func DeleteProject(c *gin.Context) {

	user, _ := c.Get("user")
	if user == nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "User not found"})
		return
	}
	id := c.Param("id")
	var projectToDelete model.Project
	if err := DB.Preload("Tags").Preload("User").First(&projectToDelete, id).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}

	if err := DB.Delete(&projectToDelete).Error; err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}

func main() {

	DB = database.InitDB()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(middleware.CORS())

	api := r.Group("/api/v1")

	api.GET("/tags", GetAllCategoriesWithTags)

	api.GET("/projects", GetProjectsByTagsInCategory)
	api.GET("/project/:id", GetFullProjectById)

	auth := api.Group("/user")
	auth.POST("/signup", middleware.RegisterHandler)
	auth.DELETE("/delete/:id", middleware.AuthMiddleware(), middleware.DeleteHandler)
	auth.PUT("/update/:id", middleware.AuthMiddleware(), middleware.UpdateHandler)
	auth.POST("/token", middleware.LoginHandler)

	auth.POST("/refresh_token", middleware.RefreshHandler)
	auth.GET("/info", middleware.AuthMiddleware(), middleware.InfoHandler)

	auth.POST("/project", middleware.AuthMiddleware(), CreateProject)
	auth.GET("/project/:id", middleware.AuthMiddleware(), GetProjectById)
	auth.PUT("/project/:id", middleware.AuthMiddleware(), UpdateProject)
	auth.DELETE("/project/:id", middleware.AuthMiddleware(), DeleteProject)

	auth.GET("/projects", middleware.AuthMiddleware(), GetUserProjects)
	/*
		var c model.Category
		var t1, t2 model.Tag
		t1 = model.Tag{Id: 1, Name: "Tag1", NameShort: "T1"}
		t2 = model.Tag{Id: 2, Name: "Tag2", NameShort: "T2"}
		c = model.Category{Id: 1, Name: "Category1", Tags: []model.Tag{t1, t2}}
		DB.Create(&c)
	*/

	/*
		var p model.Project
		var t1, t2 model.Tag
		t1 = model.Tag{Id: 1, Name: "Tag1", NameShort: "T1"}
		t2 = model.Tag{Id: 2, Name: "Tag2", NameShort: "T2"}
		p = model.Project{Id: 1, Title: "Project1", Link: "SomeLink", Description: "SomeDescription", Content: "Content", Date: time.Now(), IsLive: false, Tags: []model.Tag{t1, t2}}
		DB.Create(&p)
	*/
	/*

	 */

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

	/*
		var p model.Project
		var t1, t2 model.Tag
		t1 = model.Tag{Id: 1, Name: "Tag1", NameShort: "T1", CategoryId: 1}
		t2 = model.Tag{Id: 21, Name: "Tag21", NameShort: "T21", CategoryId: 2}
		p = model.Project{Id: 5, Title: "Project5", Link: "SomeLink", Description: "SomeDescription", Content: "Content", Date: time.Now(), IsLive: false, Tags: []model.Tag{t1, t2}}
		DB.Create(&p)
	*/
	r.Run("0.0.0.0:5297")
}
