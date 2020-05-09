package token

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kukayyou/commonlib/myconfig"
	"sync"
	"time"
)

const ISSUSER = "orangetutor"

// CustomClaims 自定义的 metadata在加密后作为 JWT 的第二部分返回给客户端
type CustomClaims struct {
	UserData UserInfo `json:"userInfo"`
	jwt.StandardClaims
}

// Token jwt服务
var (
	rwlock     sync.RWMutex
	PrivateKey string
)

type UserInfo struct {
	UserName string `json:"userName"`
	Passwd   string `json:"passwd"`
}

//检测jwt私钥是否改变
func Init(file string) {
	go func() {
		for {
			myconfig.LoadConfig(file)
			key := myconfig.Config.GetString("private_key")
			put(key)
			time.Sleep(time.Second * 10)
		}
	}()
}

//创建token
func CreateToken(userInfo UserInfo, expireTime int64) (string, error) {
	claims := CustomClaims{
		userInfo,
		jwt.StandardClaims{
			Issuer:    ISSUSER,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expireTime,
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString(get())
}

//检验token
func CheckToken(tokenStr string) (*CustomClaims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return get(), nil
	})

	if err != nil {
		return nil, err
	}
	// 解密转换类型并返回
	if claims, ok := t.Claims.(*CustomClaims); ok && t.Valid {
		return claims, nil
	}

	return nil, err
}

//获取私钥
func get() string {
	rwlock.RLock()
	defer rwlock.RUnlock()

	return PrivateKey
}

//设置私钥
func put(newKey string) {
	rwlock.Lock()
	defer rwlock.Unlock()

	PrivateKey = newKey
}
