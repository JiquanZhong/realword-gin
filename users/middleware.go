package users

import (
	"github.com/gin-gonic/gin"
	"github.com/jiquanzhong/realword-gin/common"
)

func UpdateContextUserModels(c *gin.Context, my_userr_id uint) {
	var myUserModel UserModel
	if my_userr_id != 0 {
		db := common.GetDB()
		db.First(&myUserModel, my_userr_id)
	}
	c.Set("my_user_model", myUserModel)
	c.Set("my_user_id", my_userr_id)
}
