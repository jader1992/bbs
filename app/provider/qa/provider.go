package qa

import (
	"github.com/jader1992/gocore/framework"
)

type QaProvider struct {
	framework.ServiceProvider

	c framework.Container
}

func (sp *QaProvider) Name() string {
	return QaKey
}

func (sp *QaProvider) Register(c framework.Container) framework.NewInstance {
	return NewQaService
}

func (sp *QaProvider) IsDefer() bool {
	return false
}

func (sp *QaProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}

func (sp *QaProvider) Boot(c framework.Container) error {
	return nil
}
