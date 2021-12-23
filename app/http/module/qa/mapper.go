package qa

import (
	"bbs/app/http/module/user"
	"bbs/app/provider/qa"
)

// ConvertAnswersToDTO 将answers转化为带有tree结构的AnswerDTO
func ConvertAnswersToDTO(answers []*qa.Answer) []*AnswerDTO {
	if answers == nil {
		return nil
	}

	ret := make([]*AnswerDTO, 0, len(answers))
	for _, answer := range answers {
		ret = append(ret, ConvertAnswerToDTO(answer))
	}

	return ret
}

// ConvertAnswerToDTO 将answer转化为AnswerDTO
func ConvertAnswerToDTO(answer *qa.Answer) *AnswerDTO {
	if answer == nil {
		return nil
	}

	author := user.ConvertUserToDTO(answer.Author)
	if author == nil {
		author = &user.UserDTO{
			ID: answer.AuthorID,
		}
	}

	return &AnswerDTO{
		ID:        answer.ID,
		Content:   answer.Content,
		CreatedAt: answer.CreateAt,
		UpdatedAt: answer.UpdateAt,
		Author:    author,
	}
}

// ConvertQuestionsToDTO 将questions转换为DTO
func ConvertQuestionsToDTO(questions []*qa.Question) []*QuestionDTO {
	if questions == nil {
		return nil
	}
	ret := make([]*QuestionDTO, 0, len(questions))
	for _, question := range questions {
		ret = append(ret, ConvertQuestionToDTO(question))
	}
	return ret
}

// ConvertQuestionToDTO 将问题 question 转成 QuestionDTO
func ConvertQuestionToDTO(question *qa.Question) *QuestionDTO {
	if question == nil {
		return nil
	}

	author := user.ConvertUserToDTO(question.Author)
	if author == nil {
		author = &user.UserDTO{
			ID: question.AuthorID,
		}
	}

	return &QuestionDTO{
		ID:        question.ID,
		Title:     question.Title,
		Content:   question.Context,
		CreatedAt: question.CreatedAt,
		UpdatedAt: question.UpdatedAt,
		Author:    author,
		Answers:   ConvertAnswersToDTO(question.Answers),
	}
}
