package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jiquanzhong/realword-gin/common"
	"github.com/jiquanzhong/realword-gin/users"
)

func Migrate(db *gorm.DB) {
	users.AutoMigrate()
}

func main() {
	db := common.GetDB()
	Migrate(db)
	defer db.Close()

	r := gin.Default()

	v1 := r.Group("/api")
	users.UsersRegister(v1.Group("/users"))
	users.UserRegister(v1.Group("/user"))
	users.ProfileRegister(v1.Group("/profiles"))

	r.Run(":8002")
}
