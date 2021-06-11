package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main1() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://localhost:2021/ws", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close(websocket.StatusInternalError, "内部出错!")

	err = wsjson.Write(ctx, c, "Heloo WebSocket Server")
	if err != nil {
		log.Fatal(err)
	}

	var v interface{}
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("接收到服务端响应: %v\n", v)

	c.Close(websocket.StatusNormalClosure, "")
}
