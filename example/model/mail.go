package model

import "go.mongodb.org/mongo-driver/bson/primitive"

//go:generate mongo-dao-generator -model-dir=. -model-names=Mail -dao-dir=../dao/
type Mail struct {
    ID       primitive.ObjectID `bson:"_id" gen:"autoFill"`       // 邮件ID
    Title    string             `bson:"title"`                    // 邮件标题
    Content  string             `bson:"content"`                  // 邮件内容
    Sender   int64              `bson:"sender"`                   // 邮件发送者
    Receiver int64              `bson:"receiver"`                 // 邮件接受者
    Status   int                `bson:"status"`                   // 邮件状态
    SendTime primitive.DateTime `bson:"send_time" gen:"autoFill"` // 发送时间
}