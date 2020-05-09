package token

import (
	jwt "github.com/dgrijalva/jwt-go"
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
type Token struct {
	rwlock     sync.RWMutex
	privateKey string
	//conf       config.Config
}

type UserInfo struct {
	UserName string `json:"userName"`
	Passwd   string `json:"passwd"`
}

//创建token
func (srv *Token) CreateToken(userInfo UserInfo, expireTime int64) (string, error) {
	claims := CustomClaims{
		userInfo,
		jwt.StandardClaims{
			Issuer:    ISSUSER,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expireTime,
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString(srv.get())
}

//检验token
func (srv *Token) CheckToken(tokenStr string) (*CustomClaims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return srv.get(), nil
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
func (srv *Token) get() string {
	srv.rwlock.RLock()
	defer srv.rwlock.RUnlock()

	return srv.privateKey
}

//设置私钥
func (srv *Token) Put(newKey string) {
	srv.rwlock.Lock()
	defer srv.rwlock.Unlock()

	srv.privateKey = newKey
}

/*
// InitConfig 初始化
func (srv *Token) InitConfig(address string, regisType int, path ...string) {
	if regisType == 0 {
		regisSrc := etcd.NewSource(
			consul.WithAddress(address),
			// consul.WithPrefix("/my/prefix"),
			// consul.StripPrefix(true),
		)
		srv.conf = config.NewConfig()
		err := srv.conf.Load(regisSrc)
		if err != nil {
			mylog.Error("Load regis source error:%s", err.Error())
		}
	}else{
		regisSrc := consul.NewSource(
			consul.WithAddress(address),
			// consul.WithPrefix("/my/prefix"),
			// consul.StripPrefix(true),
		)
		srv.conf = config.NewConfig()
		err := srv.conf.Load(regisSrc)
		if err != nil {
			mylog.Error("Load regis source error:%s", err.Error())
		}
	}

	value := srv.conf.Get(path...).Bytes()

	srv.put(value)
	mylog.Info("JWT privateKey:", string(srv.get()))
	srv.enableAutoUpdate(path...)
}

func (srv *Token) enableAutoUpdate(path ...string) {
	go func() {
		for {
			w, err := srv.conf.Watch(path...)
			if err != nil {
				log.Println(err)
			}
			v, err := w.Next()
			if err != nil {
				log.Println(err)
			}

			value := v.Bytes()
			srv.put(value)
			mylog.Info("New JWT privateKey:", string(srv.get()))
		}
	}()
}

//Decode 解码
func (srv *Token) Decode(tokenStr string) (*CustomClaims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return srv.get(), nil
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

// Encode 将 User 用户信息加密为 JWT 字符串
// expireTime := time.Now().Add(time.Hour * 24 * 3).Unix() 三天后过期
func (srv *Token) Encode(issuer, userName string, expireTime int64) (string, error) {
	claims := CustomClaims{
		userName,
		jwt.StandardClaims{
			Issuer:    issuer,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expireTime,
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString(srv.get())
}
*/
