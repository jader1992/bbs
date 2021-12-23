package qa

import (
	"bbs/app/http/middleware/auth"
	provider "bbs/app/provider/qa"
	"github.com/jader1992/gocore/framework/gin"
	"github.com/pkg/errors"
)

// AnswerCreate 创建问题
func (api *QApi) AnswerCreate(c *gin.Context) {
	qaService := c.MustMake(provider.QaKey).(provider.Service)

	type Param struct {
		QuestionID int64  `json:"question_id" binding:"required"`
		ParentID   int64  `json:"parent_id"`
		Content    string `json:"content" binding:"required"`
	}

	param := &Param{}
	if err := c.ShouldBind(param); err != nil {
		c.AbortWithError(404, err)
		return
	}

	user := auth.GetAuthUser(c)
	if user == nil {
		c.AbortWithError(500, errors.New("请登录后再操作"))
		return
	}

	answer := &provider.Answer{
		ID:         0,
		QuestionId: param.QuestionID,
		Content:    param.Content,
		ParentID:   param.ParentID,
		AuthorID:   user.ID,
	}

	ctx := provider.ContextWithUserID(c, user.ID)
	if err := qaService.PostAnswer(ctx, answer); err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.ISetOkStatus().IText("操作成功")
}
