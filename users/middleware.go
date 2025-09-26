package users

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"github.com/jiquanzhong/realword-gin/common"
)

func stripBearerPrefixFromTokenString(tok string) (string, error) {
	if len(tok) > 7 && tok[0:7] == "Bearer " {
		return tok[7:], nil
	}
	return tok, nil
}

var AuthorizationHeaderExtractor = &request.PostExtractionFilter{
	request.HeaderExtractor{"Authorization"},
	stripBearerPrefixFromTokenString,
}

var MyAuth2Extractor = &request.MultiExtractor{
	AuthorizationHeaderExtractor,
	request.ArgumentExtractor{"access_token"},
}

func UpdateContextUserModels(c *gin.Context, my_userr_id uint) {
	var myUserModel UserModel
	if my_userr_id != 0 {
		db := common.GetDB()
		db.First(&myUserModel, my_userr_id)
	}
	c.Set("my_user_model", myUserModel)
	c.Set("my_user_id", my_userr_id)
}

func AuthMiddleware(auto401 bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		//UpdateContextUserModels(c, 0)
		token, err := request.ParseFromRequest(c.Request, MyAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
			b := []byte(common.NBSecretPassword)
			return b, nil
		})
		if err != nil {
			if auto401 {
				c.AbortWithError(http.StatusUnauthorized, err)
			}
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			my_user_id := uint(claims["id"].(float64))
			UpdateContextUserModels(c, my_user_id)
		}
	}
}
