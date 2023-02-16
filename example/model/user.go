package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Gender int

const (
	GenderUnknown Gender = iota // 未知
	GenderMale                  // 男性
	GenderFemale                // 女性
)

// Type 用户类型
type Type int

const (
	TypeRobot   Type = 0 // 机器人用户
	TypeGuest   Type = 1 // 游客用户
	TypeGeneral Type = 2 // 普通用户
	TypeSystem  Type = 3 // 系统用户
)

// Status 用户状态
type Status int

const (
	StatusNormal    Status = iota // 正常
	StatusForbidden               // 封禁
)

//go:generate mongo-dao-generator -model-dir=. -model-names=User -dao-dir=../dao/
type User struct {
	ID             primitive.ObjectID `bson:"_id" gen:"autoFill"`
	UID            int32              `bson:"uid" gen:"autoIncr:uid"`         // 用户ID
	Account        string             `bson:"account"`                        // 用户账号
	Password       string             `bson:"password"`                       // 用户密码
	Salt           string             `bson:"salt"`                           // 密码
	Mobile         string             `bson:"mobile"`                         // 用户手机
	Email          string             `bson:"email"`                          // 用户邮箱
	Nickname       string             `bson:"nickname"`                       // 用户昵称
	Signature      string             `bson:"signature"`                      // 用户签名
	Gender         Gender             `bson:"gender"`                         // 用户性别
	Level          int                `bson:"level"`                          // 用户等级
	Experience     int                `bson:"experience"`                     // 用户经验
	Coin           int                `bson:"coin"`                           // 用户金币
	Type           Type               `bson:"type"`                           // 用户类型
	Status         Status             `bson:"status"`                         // 用户状态
	DeviceID       string             `bson:"device_id"`                      // 设备ID
	ThirdPlatforms ThirdPlatforms     `bson:"third_platforms"`                // 第三方平台
	RegisterIP     string             `bson:"register_ip"`                    // 注册IP
	RegisterTime   primitive.DateTime `bson:"register_time" gen:"autoFill"`   // 注册时间
	LastLoginIP    string             `bson:"last_login_ip"`                  // 最近登录IP
	LastLoginTime  primitive.DateTime `bson:"last_login_time" gen:"autoFill"` // 最近登录时间
}

// ThirdPlatforms 第三方平台
type ThirdPlatforms struct {
	Wechat   string `bson:"wechat"`   // 微信登录openid
	Google   string `bson:"google"`   // 谷歌登录userid
	Facebook string `bson:"facebook"` // 脸书登录userid
}
