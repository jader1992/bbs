package qa

import (
	provider "bbs/app/provider/qa"
	"github.com/jader1992/gocore/framework/gin"
)

// QuestionDetail 获取问题详情
// @Summary 获取问题详细
// @Description 获取问题详情，包括所有的答案
// @Accept json
// @Produce json
// @Tags qa
// @questionEditParam id query int true "问题ID"
// @Success 200 {string} Msg "操作成功"
// @Router /question/detail [get]
func (api *QApi) QuestionDetail(c *gin.Context) {
	qaService := c.MustMake(provider.QaKey).(provider.Service)

	id, exist := c.DefaultQueryInt64("id", 0)
	if !exist {
		c.ISetStatus(404).IText("参数错误")
		return
	}

	question, err := qaService.GetQuestion(c, id)
	if err != nil {
		c.ISetStatus(500).IText(err.Error())
		return
	}

	if err := qaService.QuestionLoadAuthor(c, question); err != nil {
		c.ISetStatus(500).IText(err.Error())
		return
	}

	if err := qaService.QuestionLoadAnswers(c, question); err != nil {
		c.ISetStatus(500).IText(err.Error())
		return
	}

	questionDTO := ConvertQuestionToDTO(question)

	c.ISetOkStatus().IJson(questionDTO)
}
