package qa

import (
	"bbs/app/http/middleware/auth"
	provider "bbs/app/provider/qa"
	"github.com/jader1992/gocore/framework/gin"
)

// AnswerDelete 删除回答
// @Summary 删除回答
// @Description 删除回答
// @Accept json
// @Produce json
// @Tags qa
// @questionCreateParam id query int true "删除id"
// @Success 200 string Msg "操作成功"
// @Router /answer/delete [get]
func (api *QApi) AnswerDelete(c *gin.Context) {
	qaService := c.MustMake(provider.QaKey).(provider.Service)

	// 参数校验
	id, exist := c.DefaultQueryInt64("id", 0)
	if !exist {
		c.ISetStatus(400).IText("参数错误")
		return
	}

	user := auth.GetAuthUser(c)

	answer, err := qaService.GetAnswer(c, id)
	if err != nil {
		c.ISetStatus(500).IText(err.Error())
		return
	}

	if answer.AuthorID != user.ID {
		c.ISetStatus(500).IText("没有权限做此操作")
	}

	if err := qaService.DeleteAnswer(c, id); err != nil {
		c.ISetStatus(500).IText(err.Error())
		return
	}

	c.ISetOkStatus().IText("操作成功")
}
