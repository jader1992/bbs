package qa

import (
	"bbs/app/http/module/user"
	"bbs/app/provider/qa"
	"github.com/jianfengye/collection"
)

func getAnswerChildren(dto *AnswerDTO, answers []*qa.Answer)  {
	if dto == nil {
		return
	}

	for _, answer := range answers {
		if dto.ID == answer.ParentID {
			if dto.Children == nil {
				dto.Children = []*AnswerDTO{}
			}

			childAnswerDTO := &AnswerDTO{
				ID:        answer.ID,
				Content:   answer.Content,
				AuthorID:  answer.AuthorID,
				CreatedAt: answer.CreateAt,
				UpdatedAt: answer.UpdateAt,
				Author:    user.ConvertUserToDTO(answer.Author),
				Children:  nil,
			}

			getAnswerChildren(childAnswerDTO, answers)

			dto.Children = append(dto.Children, childAnswerDTO)
		}
	}

	if len(dto.Children) > 0 {
		childColl := collection.NewObjPointCollection(dto.Children)
		objs := []*AnswerDTO{}
		childColl.SortByDesc("UpdatedAt").ToObjs(&objs)
		dto.Children = objs
	}

	return
}

func ConvertAnswersToDTO(answers []*qa.Answer) []*AnswerDTO {
	if answers == nil {
		return nil
	}

	answerZero := &AnswerDTO{
		ID: 0,
		Children: nil,
	}

	getAnswerChildren(answerZero, answers)
	return answerZero.Children
}

func ConvertQuestionToDTO(question *qa.Question) *QuestionDTO {
	return &QuestionDTO{
		ID:        question.ID,
		Title:     question.Title,
		Content:   question.Content,
		AuthorID:  question.AuthorID,
		CreatedAt: question.CreatedAt,
		UpdatedAt: question.UpdatedAt,
		Author:    user.ConvertUserToDTO(question.Author),
		Answers:   ConvertAnswersToDTO(question.Answers),
	}
}


