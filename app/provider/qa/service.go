package qa

import (
	"context"
	"fmt"
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/contract"
	"github.com/jader1992/gocore/framework/provider/orm"
	"github.com/jianfengye/collection"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type QaService struct {
	container framework.Container // 容器
	ormDB     *gorm.DB            // db
	logger    contract.Log        // log
}

// GetQuestions 获取所有问题
func (q *QaService) GetQuestions(ctx context.Context, pager *Pager) ([]*Question, error) {
	questions := make([]*Question, 0, pager.Size)
	total := int64(0)
	if err := q.ormDB.Count(&total).Error; err != nil {
		pager.Total = total
	}

	if err := q.ormDB.WithContext(ctx).Order("created_at desc").Offset(pager.Start).Limit(pager.Size).Find(&questions).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*Question{}, nil
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
func (q *QaService) QuestionsLoadAuthor(ctx context.Context, questions *[]*Question) error {
	if questions == nil {
		return nil
	}

	questionColl := collection.NewObjPointCollection(*questions)
	ids, err := questionColl.Pluck("ID").ToInt64s()
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}

	if err := q.ormDB.WithContext(ctx).Preload("Author").Order("created_at desc").Find(questions, ids).Error; err != nil {
		return err
	}
	return nil
}

// QuestionLoadAnswers 单个问题加载Answers
func (q *QaService) QuestionLoadAnswers(ctx context.Context, question *Question) error {
	if err := q.ormDB.WithContext(ctx).Preload("Answers", func(db *gorm.DB) *gorm.DB {
		return db.Order("answers.created_at desc")
	}).First(question).Error; err != nil {
		return err
	}
	return nil
}

// QuestionsLoadAnswers 批量问题加载Answers
func (q *QaService) QuestionsLoadAnswers(ctx context.Context, questions *[]*Question) error {
	if questions == nil {
		return nil
	}

	questionColl := collection.NewObjPointCollection(*questions)
	ids, err := questionColl.Pluck("ID").ToInt64s()
	if err != nil {
		return err
	}

	if len(ids) == 0 {
		return nil
	}

	if err := q.ormDB.WithContext(ctx).Preload("Answers").Find(questions, ids).Error; err != nil {
		return err
	}
	return nil
}

// PostAnswer 上传某个回答
func (q *QaService) PostAnswer(ctx context.Context, answer *Answer) error {
	if answer.QuestionID == 0 {
		return errors.New("问题不存在")
	}

	// 使用事务
	err := q.ormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		question := &Question{
			ID: answer.QuestionID,
		}

		// 获取问题
		if err := tx.First(question).Error; err != nil {
			return err
		}

		// 增加回答
		if err := tx.Create(answer).Error; err != nil {
			return nil
		}

		// 问题回答数量+1
		question.AnswerNum = question.AnswerNum + 1
		if err := tx.Save(question).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
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
	questionDB := &Question{ID: question.ID}
	if err := q.ormDB.WithContext(ctx).Save(questionDB).Error; err != nil {
		return errors.WithStack(err)
	}

	questionDB.UpdatedAt = time.Now()
	if question.Title != "" {
		questionDB.Title = question.Title
	}
	if question.Context != "" {
		questionDB.Context = question.Context
	}
	if err := q.ormDB.WithContext(ctx).Save(questionDB).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// AnswerLoadAuthor 问题加载Author字段
func (q *QaService) AnswerLoadAuthor(ctx context.Context, question *Answer) error {
	if err := q.ormDB.WithContext(ctx).Preload("Author").First(question).Error; err != nil {
		return err
	}
	return nil
}

// AnswersLoadAuthor 批量加载Author字段
func (q *QaService) AnswersLoadAuthor(ctx context.Context, answers *[]*Answer) error {
	if answers == nil {
		return nil
	}
	answerColl := collection.NewObjPointCollection(*answers)
	ids, err := answerColl.Pluck("ID").ToInt64s()
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}

	// 使用PreLoad的机制，获取Create方法
	if err := q.ormDB.WithContext(ctx).Preload("Author").Order("created_at desc").Find(answers, ids).Error; err != nil {
		return err
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
