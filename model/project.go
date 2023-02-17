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
	CategoryId   *int   `gorm:"primary_key" json:"categoryId"`
	CategoryName string `json:"categoryName"`
	Tags         []Tag  `json:"tags"`
}

type Tag struct {
	TagId        int    `gorm:"primary_key" json:"tagId"`
	TagName      string `json:"tagName"`
	TagNameShort string `json:"tagNameShort"`
	CategoryId   *int   `gorm:"foreignKey:CategoryId" json:"categoryId"`
}
