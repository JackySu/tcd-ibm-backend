package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"size:63;not null" json:"email"` // `json:"email"` is used to specify the name of the field when it is serialized.
	Password string `gorm:"size:255;not null"`
}

type UserDto struct {
	Email string `json:"email"`
}

func ToUserDto(user User) UserDto {
	return UserDto{
		Email: user.Email,
	}
}
