package tests

import (
	"github.com/jader1992/gocore/framework"
	"github.com/jader1992/gocore/framework/provider/app"
	"github.com/jader1992/gocore/framework/provider/env"
)

const (
	BasePath = "/Users/jade/Desktop/go/bbs/"
)

func InitBaseContainer() framework.Container {
	// 初始化服务容器
	container := framework.NewGocoreContainer()
	// 绑定App服务提供者
	container.Bind(&app.GocoreAppProvider{BaseFolder: BasePath})
	// 后续初始化需要绑定的服务提供者...
	container.Bind(&env.GocoreEnvProvider{})
	return container
}
