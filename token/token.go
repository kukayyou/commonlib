package token

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kukayyou/commonlib/myetcd"
	"sync"
	"time"
)

const ISSUSER = "orangetutor"

// CustomClaims 自定义的 metadata在加密后作为 JWT 的第二部分返回给客户端
type CustomClaims struct {
	UserData UserInfo `json:"userInfo"`
	jwt.StandardClaims
}

// CustomClaims 自定义的 metadata在加密后作为 JWT 的第二部分返回给客户端
type ServerClaims struct {
	Server string `json:"server"`
	jwt.StandardClaims
}

// Token jwt服务
var (
	rwlock     sync.RWMutex
	PrivateKey string = "orangetutor"
)

type UserInfo struct {
	UserID   string  `json:"userId"`
}

//检测jwt私钥是否改变
func Init(opt string) {
	go func() {
		for {
			/*myconfig.LoadConfig(file)
			key := myconfig.Config.GetString("private_key")*/
			key := myetcd.GetKey(opt, "jwtkey")
			put(key)
			time.Sleep(time.Second * 10)
		}
	}()
}

//创建token
func CreateUserToken(userInfo UserInfo, expireTime int64) (string, error) {
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

//创建token
func CreateServerToken(server string, expireTime int64) (string, error) {
	claims := ServerClaims{
		server,
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
func CheckUserToken(tokenStr string) (*CustomClaims, error) {
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

//检验token
func CheckServerToken(tokenStr string) (*ServerClaims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &ServerClaims{}, func(token *jwt.Token) (interface{}, error) {
		return get(), nil
	})

	if err != nil {
		return nil, err
	}
	// 解密转换类型并返回
	if claims, ok := t.Claims.(*ServerClaims); ok && t.Valid {
		return claims, nil
	}

	return nil, err
}

//获取私钥
func get() []byte {
	rwlock.RLock()
	defer rwlock.RUnlock()

	return []byte(PrivateKey)
}

//设置私钥
func put(newKey string) {
	rwlock.Lock()
	defer rwlock.Unlock()

	PrivateKey = newKey
}
