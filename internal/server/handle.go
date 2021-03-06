package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Xihy-xch/tcp-chatroom/internal/logic"
)

func RegisterHandle() {
	inferRootDir()

	go logic.Broadcaster.Start()

	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)
}

var rootDir string

func inferRootDir() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("os.Getwd err: ", err)
	}

	var infer func(d string) string

	infer = func(d string) string {
		if exists(d + "/template") {
			return d
		}

		return infer(filepath.Dir(d))
	}

	rootDir = infer(cwd)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}