package middleware

import (
	"net/http"
	"strings"
	"sweng_backend/database"
	"sweng_backend/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var err error

// openssl rand -hex 32
var jwtKey = []byte("78ede33d04003e331827f8a6658fc44c378a55f95907c73941f842da71ca7731")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RegisterHandler(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	if email == "" || password == "" {
		c.JSON(400, gin.H{"status": "email or password is empty"})
		return
	}
	db := database.GetDB()
	if err = db.Where("email = ?", email).First(&model.User{}).Error; err == nil {
		c.JSON(400, gin.H{"status": "email already exists"})
		return
	}
	hashedPassword, _ := HashPassword(password)
	user := model.User{Email: email, Password: string(hashedPassword)}
	db.Create(&user)
	c.JSON(200, gin.H{"status": "success"})
}

func LoginHandler(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	if email == "" || password == "" {
		c.JSON(400, gin.H{"status": "email or password is empty"})
		return
	}
	db := database.GetDB()
	var user model.User
	if err = db.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"status": "email does not exist"})
		return
	}
	if !CheckPasswordHash(password, user.Password) {
		c.JSON(400, gin.H{"status": "wrong password"})
		return
	}
	token, err := ReleaseToken(user)
	if err != nil {
		c.JSON(400, gin.H{"status": "token release failed"})
		return
	}
	c.JSON(200, gin.H{"status": "success", "token": token})
}

func RefreshHandler(c *gin.Context) {
	tokenString := c.PostForm("token")
	if tokenString == "" {
		c.JSON(400, gin.H{"status": "token is empty"})
		return
	}
	token, claims, err := ParseToken(tokenString)
	if err != nil || !token.Valid {
		c.JSON(400, gin.H{"status": "token is invalid"})
		return
	}
	db := database.GetDB()
	var user model.User
	if err = db.Where("id = ?", claims.UserId).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"status": "user does not exist"})
		return
	}
	newToken, err := ReleaseToken(user)
	if err != nil {
		c.JSON(400, gin.H{"status": "token release failed"})
		return
	}
	c.JSON(200, gin.H{"status": "success", "token": newToken})
}

func InfoHandler(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"user": model.ToUserDto(user.(model.User))},
	})
}

func ReleaseToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "tcd-ibm-swEng",
			Subject:   "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return token, claims, err
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// authorization header
		tokenString := ctx.GetHeader("Authorization")

		// validate token format
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "unauthorized",
			})
			ctx.Abort()
			return
		}

		// 'bearer ' is 7 characters long, so remove it from the tokenString
		tokenString = tokenString[7:]

		token, claims, err := ParseToken(tokenString)

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "unauthorized",
			})
			ctx.Abort()
			return
		}

		//token通过验证, 获取claims中的UserID
		userId := claims.UserId
		DB := database.GetDB()
		var user model.User
		DB.First(&user, userId)

		// 验证用户是否存在
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "unauthorized",
			})
			ctx.Abort()
			return
		}

		//用户存在 将user信息写入上下文
		ctx.Set("user", user)

		ctx.Next()
	}
}
