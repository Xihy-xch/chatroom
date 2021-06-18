package logic

import (
	"expvar"
	"fmt"
)

func init() {
	expvar.Publish("message_queue", expvar.Func(calcMessageQueueLen))
}

func calcMessageQueueLen() interface{} {
	fmt.Println("===len=:", len(Broadcaster.messageList))
	return len(Broadcaster.messageList)
}

var Broadcaster = &broadcaster{
	users:                 make(map[string]*User),
	enteringList:          make(chan *User),
	leavingList:           make(chan *User),
	messageList:           make(chan *Message, global.MessageQueueLen),
	checkUserChannel:      make(chan string),
	checkUserCanInChannel: make(chan bool),
}

type broadcaster struct {
	users map[string]*User

	enteringList chan *User
	leavingList  chan *User
	messageList  chan *Message

	checkUserChannel      chan string
	checkUserCanInChannel chan bool

	requestUsersChannel chan struct{}
	usersChannel        chan []*User
}

func (b *broadcaster) Start() {
	for {
		select {
		case user := <-b.enteringList:
			b.users[user.NickName] = user

			b.sendUserList()
		case user := <-b.leavingList:
			delete(b.users, user.NickName)
			user.CloseMessageList()

			b.sendUserList()
		case msg := <-b.messageList:
			for _, user := range b.users {
				if user.UID == msg.User.UID {
					continue
				}
				user.MessageChannel <- msg
			}
		case nickname := <-b.checkUserChannel:
			if _, ok := b.users[nickname]; ok {
				b.checkUserCanInChannel <- false
			} else {
				b.checkUserCanInChannel <- true
			}
		}
	}
}

func (b *broadcaster) UserEntering(u *User) {
	b.enteringList <- u
}

func (b *broadcaster) UserLeaving(u *User) {
	b.leavingList <- u
}

func (b *broadcaster) Broadcast(msg *Message) {
	b.messageList <- msg
}


func (b *broadcaster) CanEnterRoom(nickname string) bool {
	b.checkUserChannel <- nickname

	return <-b.checkUserCanInChannel
}


func (b *broadcaster) GetUserList() []*User {
	b.requestUsersChannel <- struct{}{}
	return <-b.usersChannel
}