package logic

import (
	"time"

	"github.com/spf13/cast"
)

type Message struct {
	User    *User     `json:"user,omitempty"`
	Type    int       `json:"type,omitempty"`
	Content string    `json:"content,omitempty"`
	MsgTime time.Time `json:"msg_time"`

	ClientSendTime time.Time `json:"client_send_time"`

	Ats   []string         `json:"ats,omitempty"`
	Users map[string]*User `json:"users,omitempty"`
}

const (
	MsgTypeNormal    = iota + 1 //普通用户消息
	MsgTypeWelcome              // 当前用户欢迎消息
	MsgTypeUserEnter            // 用户进入
	MsgTypeUserLeave            // 用户退出
	MsgTypeError
)

var System = &User{}

func NewMessage(user *User, content, clientTime string) *Message {
	message := &Message{
		User:    user,
		Type:    MsgTypeNormal,
		Content: content,
		MsgTime: time.Now(),
	}
	if clientTime != "" {
		message.ClientSendTime = time.Unix(0, cast.ToInt64(clientTime))
	}
	return message
}

func NewWelcomeMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeWelcome,
		Content: user.NickName + " 您好，欢迎加入聊天室！",
		MsgTime: time.Now(),
	}
}

func NewUserEnterMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeUserEnter,
		Content: user.NickName + " 加入了聊天室",
		MsgTime: time.Now(),
	}
}

func NewUserLeaveMessage(user *User) *Message {
	return &Message{
		User:    user,
		Type:    MsgTypeUserLeave,
		Content: user.NickName + " 离开了聊天室",
		MsgTime: time.Now(),
	}
}

func NewErrorMessage(content string) *Message {
	return &Message{
		User:    System,
		Type:    MsgTypeError,
		Content: content,
		MsgTime: time.Now(),
	}
}
