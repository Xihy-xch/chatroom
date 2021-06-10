package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	enteringUsers       = make(chan *User)
	leavingUsers        = make(chan *User)
	messageBuffer       = make(chan Message, 8)
	idGenerater   int64 = 0
)

type Message struct {
	OwnerID int64
	Content string
}

func main() {
	listener, err := net.Listen("tcp", ":2020")
	if err != nil {
		log.Fatal("监听失败 err: ", err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConn(conn)
	}
}

type User struct {
	ID      int64
	Addr    string
	EnterAt time.Time
	Message chan string
}

func (u *User) String() string {
	return strconv.FormatInt(u.ID, 10)
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	user := &User{
		ID:      GenUserID(),
		Addr:    conn.RemoteAddr().String(),
		EnterAt: time.Now(),
		Message: make(chan string, 8),
	}

	go sendMessage(conn, user.Message)

	message := joinString("Welcome, ", user.String())
	user.Message <- message

	message = joinString("user:`" + user.String() + "`has enter")
	messageBuffer <- Message{
		OwnerID: 0,
		Content: message,
	}
	enteringUsers <- user

	var userActive = make(chan struct{})

	go func() {
		d := 5 * time.Minute
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				conn.Close()
			case <-userActive:
				timer.Reset(d)
			}
		}
	}()

	input := bufio.NewScanner(conn)
	for input.Scan() {
		message = joinString(user.String(), ":", input.Text())
		messageBuffer <- Message{
			OwnerID: user.ID,
			Content: message,
		}

		userActive <- struct{}{}
	}

	if err := input.Err(); err != nil {
		log.Println("读取错误: ", err)
	}

	leavingUsers <- user
	message = joinString("user:`", user.String()+"`has left")
	messageBuffer <- Message{
		OwnerID: user.ID,
		Content: message,
	}
}

func joinString(strs ...string) string {
	var builder strings.Builder

	for _, str := range strs {
		builder.WriteString(str)
	}
	return builder.String()
}

func sendMessage(conn net.Conn, message <-chan string) {
	for msg := range message {
		_, err := fmt.Fprintln(conn, msg)
		if err != nil {
			log.Println("sendMessage Println err: ", err)
			return
		}
	}
}

func GenUserID() int64 {
	return atomic.AddInt64(&idGenerater, 1)
}

func broadcaster() {
	users := make(map[*User]struct{})

	for {
		select {
		case user := <-enteringUsers:
			users[user] = struct{}{}
		case user := <-leavingUsers:
			delete(users, user)
			close(user.Message)
		case msg := <-messageBuffer:
			for user := range users {
				if user.ID == msg.OwnerID {
					continue
				}
				user.Message <- msg.Content
			}
		}
	}
}
