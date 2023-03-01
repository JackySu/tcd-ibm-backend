package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

/*
Dto stands for Data Transfer Object
It is used to transfer data between different layers
For example, when we want to return a user to the client,
we don't want to return the password hash to the client.
So we create a UserDto struct to transfer the data.
*/

type Category struct {
	Id   int    `gorm:"primary_key" json:"id"`
	Name string `json:"name"`
	Tags []Tag  `json:"tags"`
}

type CategoryDto struct {
	Id   int      `json:"id"`
	Name string   `json:"name"`
	Tags []TagDto `json:"tags"`
}

func ToCategoryDto(category Category) CategoryDto {
	return CategoryDto{
		Id:   category.Id,
		Name: category.Name,
		Tags: ToTagDtoList(category.Tags),
	}
}

func ToCategoryDtoList(categories []Category) []CategoryDto {
	var categoryDtos []CategoryDto
	for _, category := range categories {
		categoryDtos = append(categoryDtos, ToCategoryDto(category))
	}
	return categoryDtos
}

type ProjectBase struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Tags        []int  `json:"tags"`
}

type UpdateProject struct {
	Title       *string `json:"title,omitempty"`
	Link        *string `json:"link,omitempty"`
	Description *string `json:"description,omitempty"`
	Content     *string `json:"content,omitempty"`
	Tags        *[]int  `json:"tags,omitempty"`
	IsLive      *bool   `json:"is_live,omitempty"`
}

type Project struct {
	gorm.Model
	Email       string `json:"email"`
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Content     string `json:"content"`
	IsLive      bool   `json:"is_live"`
	UserID      uint   `json:"user_id"`
	User        User   `json:"user"` // User that created the project
	Tags        []Tag  `gorm:"many2many:project_tags"`
}

type ProjectDto struct {
	ID          int        `json:"id"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
	Email       string     `json:"email"`
	Title       string     `json:"title"`
	Link        string     `json:"link"`
	Description string     `json:"description"`
	IsLive      bool       `json:"is_live"`
	User        UserDto    `json:"user"`
	Tags        []TagDto   `json:"tags"`
}

type ProjectFullDto struct {
	ProjectDto
	Content string `json:"content"`
}

func ToProjectDto(project Project) ProjectDto {
	return ProjectDto{
		ID:          int(project.ID),
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
		DeletedAt:   project.DeletedAt,
		Email:       project.Email,
		Title:       project.Title,
		Link:        project.Link,
		Description: project.Description,
		IsLive:      project.IsLive,
		User:        ToUserDto(project.User),
		Tags:        ToTagDtoList(project.Tags),
	}
}

func ToProjectFullDto(project Project) ProjectFullDto {
	return ProjectFullDto{
		ProjectDto: ToProjectDto(project),
		Content:    project.Content,
	}
}

func ToProjectDtoList(projects []Project) []ProjectDto {
	var projectDtos []ProjectDto
	for _, project := range projects {
		projectDtos = append(projectDtos, ToProjectDto(project))
	}
	return projectDtos
}

type Tag struct {
	Id         int       `gorm:"primary_key" json:"id"`
	Name       string    `json:"name"`
	NameShort  string    `json:"nameShort"`
	CategoryId int       `json:"categoryId"`
	Projects   []Project `gorm:"many2many:project_tags" json:"projects"`
}

type TagDto struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	NameShort  string `json:"nameShort"`
	CategoryId int    `json:"categoryId"`
}

func ToTagDto(tag Tag) TagDto {
	return TagDto{
		Id:         tag.Id,
		Name:       tag.Name,
		NameShort:  tag.NameShort,
		CategoryId: tag.CategoryId,
	}
}

func ToTagDtoList(tags []Tag) []TagDto {
	var tagDtos []TagDto
	for _, tag := range tags {
		tagDtos = append(tagDtos, ToTagDto(tag))
	}
	return tagDtos
}
