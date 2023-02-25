package model

import "time"

/*
Dto stands for Data Transfer Object
It is used to transfer data between different layers
For example, when we want to return a user to the client,
we don't want to return the password hash to the client.
So we create a UserDto struct to transfer the data.
*/

type Announcement struct {
	Aid              int    `gorm:"primary_key" json:"aid"`
	Title            string `gorm:"size:255" json:"title"`
	TitleLink        string `gorm:"size:255" json:"titleLink"`
	Date             string `gorm:"size:63" json:"date"`
	AnnouncementType string `gorm:"size:63" json:"announcementType"`
	Desc             string `gorm:"size:255" json:"desc"`
}

type Detail struct {
	PPid        int    `gorm:"primary_key" json:"ppid"`
	Description string `gorm:"size:255" json:"description"`
	Url         string `gorm:"size:255" json:"url"`
	Type        string `gorm:"size:63" json:"type"`
}

type PA struct {
	PPid        int    `gorm:"primary_key" json:"ppid"`
	Heading     string `json:"Heading"`
	Summary     string `json:"Summary"`
	Product     string `json:"Product"`
	Solutions   string `json:"Solutions"`
	Vertical    string `json:"Vertical"`
	Image1Url   string `json:"Image1Url"`
	ProductType string `json:"ProductType"`
	DetailPage  string `json:"DetailPage"`
	IsLive      bool   `json:"islive"`
	IsNew       bool   `json:"isnew"`
	MetaDesc    string `json:"metaDesc"`
	MetaKeyword string `json:"metaKeyword"`
}

type Product struct {
	Pid   string `gorm:"primary_key" json:"pid"`
	PName string `gorm:"size:255" json:"pname"`
	PLink string `json:"plink"`
}

type Solution struct {
	Sid   string `gorm:"primary_key" json:"sid"`
	SName string `gorm:"size:255" json:"sname"`
}

type Vertical struct {
	Vid   string `gorm:"primary_key" json:"vid"`
	VName string `gorm:"size:255" json:"vname"`
}

type Type struct {
	Tid      string `gorm:"primary_key" json:"tid"`
	TypeName string `gorm:"size:255" json:"typename"`
}

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

func (c *Category) ToDto() CategoryDto {
	return CategoryDto{
		Id:   c.Id,
		Name: c.Name,
		Tags: ToTagDto(c.Tags),
	}
}

func ToCategoryDto(categories []Category) []CategoryDto {
	var categoryDtos []CategoryDto
	for _, category := range categories {
		categoryDtos = append(categoryDtos, category.ToDto())
	}
	return categoryDtos
}

type Project struct {
	Id          int       `gorm:"primary_key" json:"id"`
	Email       string    `json:"email"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Date        time.Time `json:"date"`
	IsLive      bool      `json:"is_live"`
	Tags        []Tag     `gorm:"many2many:project_tags" json:"tags"`
}

type ProjectDto struct {
	Id          int       `json:"id"`
	Email       string    `json:"email"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Date        time.Time `json:"date"`
	IsLive      bool      `json:"is_live"`
	Tags        []TagDto  `json:"tags"`
}

func (p *Project) ToDto() ProjectDto {
	return ProjectDto{
		Id:          p.Id,
		Email:       p.Email,
		Title:       p.Title,
		Link:        p.Link,
		Description: p.Description,
		Date:        p.Date,
		IsLive:      p.IsLive,
		Tags:        ToTagDto(p.Tags),
	}
}

func ToProjectDto(projects []Project) []ProjectDto {
	var projectDtos []ProjectDto
	for _, project := range projects {
		projectDtos = append(projectDtos, project.ToDto())
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

func (t *Tag) ToDto() TagDto {
	return TagDto{
		Id:         t.Id,
		Name:       t.Name,
		NameShort:  t.NameShort,
		CategoryId: t.CategoryId,
	}
}

func ToTagDto(tags []Tag) []TagDto {
	var tagDtos []TagDto
	for _, tag := range tags {
		tagDtos = append(tagDtos, tag.ToDto())
	}
	return tagDtos
}
