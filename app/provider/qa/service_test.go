package qa

import (
	test "bbs/test"
	"fmt"
	"github.com/jader1992/gocore/framework/contract"
	"github.com/jader1992/gocore/framework/provider/cache"
	"github.com/jader1992/gocore/framework/provider/config"
	"github.com/jader1992/gocore/framework/provider/log"
	"github.com/jader1992/gocore/framework/provider/orm"
	"github.com/jader1992/gocore/framework/provider/redis"
	"testing"
)

func TestQA(t *testing.T) {
	container := test.InitBaseContainer()
	container.Bind(&config.GocoreConfigProvider{})
	container.Bind(&log.GocoreLogServiceProvider{})
	container.Bind(&orm.GormProvider{})
	container.Bind(&redis.GocoreRedisProvider{})
	container.Bind(&cache.GocoreCacheProvider{})

	ormService := container.MustMake(contract.ORMKey).(contract.IORMService)
	db, err := ormService.GetDB(orm.WithConfigPath("database.default"))
	if err != nil {
		t.Fatal(err)
	}

	if err := db.AutoMigrate(&Question{}, &Answer{}); err != nil {
		t.Fatal(err)
	}

	tmp, err := NewQaService(container)
	if err != nil {
		t.Fatal(err)
	}
	qaService := tmp.(Service)
	fmt.Println(qaService)
}
