package user

import (
	"bbs/test"
	"context"
	"github.com/jader1992/gocore/framework/contract"
	"github.com/jader1992/gocore/framework/provider/cache"
	"github.com/jader1992/gocore/framework/provider/config"
	"github.com/jader1992/gocore/framework/provider/log"
	"github.com/jader1992/gocore/framework/provider/orm"
	"github.com/jader1992/gocore/framework/provider/redis"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUserRegisterLogin(t *testing.T)  {
	container := tests.InitBaseContainer()
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

	if err := db.AutoMigrate(&User{});err != nil {
		t.Fatal(err)
	}

	tmp, err := NewUserService(container)
	if err != nil {
		t.Fatal(err)
	}
	us := tmp.(*UserService)
	ctx := context.Background()

	user1 := &User{
		UserName: "guyi",
		Password: "123456",
		Email: "1249200310@qq.com",
	}

	Convey("正常流程", t, func() {

		Convey("注册用户", func() {
			userWithToken, err := us.Register(ctx, user1)
			So(err, ShouldBeNil)
			user1.Token = userWithToken.Token
		})

		Convey("发送邮件", func() {
			err := us.SendRegisterMail(ctx, user1)
			So(err, ShouldBeNil)
		})

		Convey("验证注册信息", func() {
			isExist, err := us.VerifyRegister(ctx, user1.Token)
			So(err, ShouldBeNil)
			So(isExist, ShouldBeTrue)

			// 数据库有数据
			userDB := &User{}
			err = db.Where("username=?", user1.UserName).First(userDB).Error
			So(err, ShouldBeNil)
			So(userDB.ID, ShouldNotBeZeroValue)
		})

		Convey("用户登录", func() {
			userDB, err := us.Login(ctx, user1)
			So(err, ShouldBeNil)
			So(userDB, ShouldNotBeNil)
			user1.Token = userDB.Token
		})

		Convey("用户验证", func() {
			userSession, err := us.VerifyLogin(ctx, user1.Token)
			So(err, ShouldBeNil)
			So(userSession, ShouldNotBeNil)
		})

		Convey("用户登出", func() {
			err := us.Logout(ctx, user1)
			So(err, ShouldBeNil)
			//重新验证为失败
			userSession, err := us.VerifyLogin(ctx, user1.Token)
			So(err, ShouldNotBeNil)
			So(userSession, ShouldBeNil)
		})
	})


}
