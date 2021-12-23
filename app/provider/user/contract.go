package user

import (
    "context"
    "time"
)

const UserKey = "user"

// User 代表一个用户，注意这里的用户信息字段在不同接口和参数可能为空
type User struct {
    ID int64 // 代表用户id, 只有注册成功之后才有这个id，唯一表示一个用户
    UserName string
    Password string
    Email string
    Avatar string
    Token string

    CreatedAt time.Time
}

// Service 用户相关的服务
type Service interface {
    // Register 注册用户,注意这里只是将用户注册, 并没有激活, 需要调用
    // 参数：user必填，username，password, email
    // 返回值： user 带上token
    Register(ctx context.Context, user *User) (*User, error)
    // SendRegisterMail 发送注册的邮件
    // 参数：user必填： username, password, email, token
    SendRegisterMail(ctx context.Context, user *User) error
    // VerifyRegister 注册用户，验证注册信息, 返回验证是否成功
    VerifyRegister(ctx context.Context, token string) (bool, error)

    // Login 登录相关，使用用户名密码登录，获取完成User信息
    Login(ctx context.Context, user *User) (*User, error)
    // Logout 登出
    Logout(ctx context.Context, user *User) error
    // VerifyLogin 登录验证
    VerifyLogin(ctx context.Context, token string) (*User, error)
}
