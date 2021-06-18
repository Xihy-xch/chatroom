package server

import (
	"log"
	"net/http"

	"github.com/Xihy-xch/tcp-chatroom/internal/logic"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func WebSocketHandleFunc(w http.ResponseWriter, req *http.Request) {
	conn, err := websocket.Accept(w, req, nil)
	if err != nil {
		log.Println("websocket accept error: ", err)
		return
	}

	token := req.FormValue("token")
	nickname := req.FormValue("nickname")

	if l := len(nickname); l < 2 || l > 20 {
		log.Println("nickname illegal: ", nickname)
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("非法昵称, 昵称长度: 4-20"))
		conn.Close(websocket.StatusUnsupportedData, "nickname illegal!")
		return
	}

	if !logic.Broadcaster.CanEnterRoom(nickname) {
		log.Println("昵称已存在: ", nickname)
		wsjson.Write(req.Context(), conn, logic.NewErrorMessage("该昵称已存在! "))
		conn.Close(websocket.StatusUnsupportedData, "nickname exists!")
		return
	}

	user := logic.NewUser(conn, token, nickname, req.RemoteAddr)

	go user.SendMessage(req.Context())

	user.MessageChannel <- logic.NewWelcomeMessage(user)

	msg := logic.NewUserEnterMessage(user)
	logic.Broadcaster.Broadcast(msg)
	logic.Broadcaster.UserEntering(user)
	log.Println("user:", nickname, " joins chat")

	err = user.ReceiveMessage(req.Context())

	logic.Broadcaster.UserLeaving(user)
	msg = logic.NewUserLeaveMessage(user)
	logic.Broadcaster.Broadcast(msg)
	log.Println("user:", nickname, " leaves chat")

	if err != nil {
		log.Println("read from client error:", err)
		conn.Close(websocket.StatusInternalError, "read from client error")
	} else {
		conn.Close(websocket.StatusNormalClosure, "")
	}

}
