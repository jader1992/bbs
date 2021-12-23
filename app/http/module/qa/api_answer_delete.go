package qa

import (
	"bbs/app/http/middleware/auth"
	provider "bbs/app/provider/qa"
	"github.com/jader1992/gocore/framework/gin"
)

func (api *QApi) AnswerDelete(c *gin.Context) {
	qaService := c.MustMake(provider.QaKey).(provider.Service)

	type Param struct {
		ID int64 `json:"id" binding:"required"`
	}

	param := &Param{}
	if err := c.ShouldBind(param); err != nil {
		c.AbortWithError(404, err)
		return
	}

	user := auth.GetAuthUser(c)

	ctx := provider.ContextWithUserID(c, user.ID)
	if err := qaService.DeleteQuestion(ctx, param.ID); err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.ISetOkStatus().IText("操作成功")
}