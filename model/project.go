package model

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
	Id   *int   `gorm:"primary_key" json:"id"`
	Name string `json:"name"`
	Tags []Tag  `json:"tags"`
}

type Project struct {
	Id          *int   `gorm:"primary_key" json:"id"`
	Email       string `json:"email"`
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Date        string `json:"date"`
	IsLive      bool   `json:"is_live"`
	Tags        []Tag  `gorm:"many2many:project_tags" json:"tags"`
}

type Tag struct {
	Id         int       `gorm:"primary_key" json:"id"`
	Name       string    `json:"name"`
	NameShort  string    `json:"nameShort"`
	CategoryId *int      `json:"categoryId"`
	Projects   []Project `gorm:"many2many:project_tags" json:"projects"`
}
