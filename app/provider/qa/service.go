package qa

import (
	"context"
	"fmt"
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
	"github.com/jader1992/gocore/framework/provider/orm"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type QaService struct {
	container framework.Container // 容器
	ormDB     *gorm.DB            // db
	logger    contract.Log        // log
}

// GetQuestions 获取所有问题
func (q *QaService) GetQuestions(ctx context.Context, pager *Pager) ([]*Question, error) {
	questions := make([]*Question, 0, pager.Size)
	if err := q.ormDB.WithContext(ctx).Order("created_at desc").Offset(pager.Start).Limit(pager.Size).Find(&questions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return questions, nil
}

// GetQuestion 获取某一个问题
func (q *QaService) GetQuestion(ctx context.Context, questionID int64) (*Question, error) {
	question := &Question{}
	if err := q.ormDB.WithContext(ctx).First(question, questionID).Error; err != nil {
		return nil, err
	}
	return question, nil
}

// PostQuestion 创建一个问题
func (q *QaService) PostQuestion(ctx context.Context, question *Question) error {
	if err := q.ormDB.WithContext(ctx).Create(question).Error; err != nil {
		return err
	}
	return nil
}

// QuestionLoadAuthor 问题加载Author字段
func (q *QaService) QuestionLoadAuthor(ctx context.Context, question *Question) error {
	if err := q.ormDB.WithContext(ctx).Preload("Author").Find(question).Error; err != nil {
		return err
	}
	return nil
}

// QuestionsLoadAuthor 批量加载Author字段
func (q *QaService) QuestionsLoadAuthor(ctx context.Context, questions []*Question) error {
	if err := q.ormDB.WithContext(ctx).Preload("Author").Find(questions).Error; err != nil {
		return err
	}
	return nil
}

// QuestionLoadAnswers 单个问题加载Answers
func (q *QaService) QuestionLoadAnswers(ctx context.Context, question *Question) error {
	if err := q.ormDB.WithContext(ctx).Preload("Answers").First(question).Error; err != nil {
		return err
	}
	return nil
}

// QuestionsLoadAnswers 批量问题加载Answers
func (q *QaService) QuestionsLoadAnswers(ctx context.Context, questions []*Question) error {
	if err := q.ormDB.WithContext(ctx).Preload("Answers").Find(questions).Error; err != nil {
		return err
	}
	return nil
}

// PostAnswer 上传某个回答
func (q *QaService) PostAnswer(ctx context.Context, answer *Answer) error {
	if err := q.ormDB.WithContext(ctx).Preload("Question").First(answer).Error; err != nil {
		return err
	}

	if answer.Question == nil {
		return errors.New("问题不存在")
	}

	if err := q.ormDB.WithContext(ctx).Create(answer).Error; err != nil {
		return err
	}

	answer.Question.AnswerNum = answer.Question.AnswerNum + 1
	if err := q.ormDB.WithContext(ctx).Save(answer.Question).Error; err != nil {
		return err
	}
	return nil
}

// GetAnswer 获取回答
func (q *QaService) GetAnswer(ctx context.Context, answerID int64) (*Answer, error) {
	answer := &Answer{ID: answerID}
	if err := q.ormDB.WithContext(ctx).First(answer).Error; err != nil {
		return nil, err
	}
	return answer, nil
}

// DeleteQuestion 删除问题，同时删除对应的回答
func (q *QaService) DeleteQuestion(ctx context.Context, questionID int64) error {
	question := &Question{ID: questionID}
	if err := q.ormDB.WithContext(ctx).Delete(question).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// DeleteAnswer 删除某个回答
func (q *QaService) DeleteAnswer(ctx context.Context, answerID int64) error {
	answer := &Answer{ID: answerID}
	if err := q.ormDB.WithContext(ctx).Delete(answer).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateQuestion 更新问题
func (q *QaService) UpdateQuestion(ctx context.Context, question *Question) error {
	if err := q.ormDB.WithContext(ctx).Save(question).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func NewQaService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	ormService := container.MustMake(contract.ORMKey).(contract.IORMService)
	logger := container.MustMake(contract.LogKey).(contract.Log)

	db, err := ormService.GetDB(orm.WithConfigPath("database.default"))
	if err != nil {
		logger.Error(context.Background(), "获取gormDB错误", map[string]interface{}{
			"err": fmt.Sprintf("%+v", err),
		})
		return nil, err
	}
	return &QaService{container: container, ormDB: db, logger: logger}, nil
}
