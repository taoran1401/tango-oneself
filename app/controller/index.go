package controller

import (
	"github.com/gin-gonic/gin"
)

// index controller
type IndexController struct {
}

func NewIndexController() *IndexController {
	return &IndexController{}
}

func (this *IndexController) Index(ctx *gin.Context) {
	/*token, _ := global.JWT.CreateToken()
	c.JSON(200, gin.H{
		"token": token,
	})*/

	/*mapClaims := jwt.MapClaims{
		"userId": "111",
	}
	token, _ := global.JWT.CreateToken(mapClaims)
	fmt.Println(token)*/

	/*token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjU0MjcxfQ.Nd0n6Mfif3YbFAIGYDF1eE-SRSjeUshk2Mb98Nufsng"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoidGVzdCJ9.pet0FltwMs4P3Z1OMfKbRfGZk1WnWUyzCu_BW8w2eec"*/
	/*token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxMTEifQ.zozYKsWmfKJZ-uHsPmes0jVlOlinOIHk90tkMG5ZQX0"
	res, _ := global.JWT.ParseToken(token)
	c.JSON(200, gin.H{
		"token": res,
	})*/
	//userId, _ := ctx.Get("user_info")
	ctx.JSON(200, gin.H{
		"data": "ok",
	})

}
