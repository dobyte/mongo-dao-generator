package mail

import "go.mongodb.org/mongo-driver/bson/primitive"

// Status 邮件状态
type Status int

const (
	StatusUnread   Status = 0 // 未读
	StatusRead     Status = 1 // 已读
	StatusReceived Status = 2 // 已领取
)

type Sender struct {
	ID   int64  `bson:"id"`   // 发送者ID，官方发送者ID为0，系统邮件为负数，用户发送的为正数
	Name string `bson:"name"` // 发送者名称
	Icon string `bson:"icon"` // 发送者图标，仅在发送者为用户时存在
}

type Attachment struct {
	PropID  int `bson:"prop_id"`  // 道具ID
	PropNum int `bson:"prop_num"` // 道具数量
}

type Mail struct {
	ID          primitive.ObjectID `bson:"_id" gen:"autoFill"`       // 邮件ID
	Title       string             `bson:"title"`                    // 邮件标题
	Content     string             `bson:"content"`                  // 邮件内容
	Sender      Sender             `bson:"sender"`                   // 邮件发送者
	Receiver    int64              `bson:"receiver"`                 // 邮件接受者
	Attachments []Attachment       `bson:"attachments"`              // 邮件附件
	Status      Status             `bson:"status"`                   // 邮件状态
	SendTime    primitive.DateTime `bson:"send_time" gen:"autoFill"` // 发送时间
}
