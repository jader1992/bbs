package qa

import (
	"bbs/app/provider/qa"
	"github.com/jader1992/gocore/framework/gin"
)

type QApi struct {

}

func RegisterRoutes(r *gin.Engine) error {
	api := &QApi{}
	if !r.IsBind(qa.QaKey) {
		r.Bind(&qa.QaProvider{})
	}

	// 问题列表
	r.GET("/question/list", api.QuestionList)
	// 问题详情
	r.GET("/question/detail", api.QuestionDetail)
	// 创建问题
	r.POST("/question/create", api.QuestionCreate)
	// 删除问题
	r.POST("/question/delete", api.QuestionDelete)
	// 更新问题
	r.POST("/question/edit", api.QuestionEdit)
	// 创建回答
	r.POST("/answer/create", api.AnswerCreate)
	// 删除回答
	r.POST("/answer/delete", api.AnswerDelete)

	return nil
}


