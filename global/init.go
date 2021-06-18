package global

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

func init() {
	Init()
}

var RootDir string

var once = new(sync.Once)

func Init() {
	once.Do(func() {
		inferRootDir()
		initConfig()
	})
}

func inferRootDir() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var infer func(d string) string
	infer = func(d string) string {
		if exists(d + "/template") {
			return d
		}

		return infer(filepath.Dir(d))
	}

	RootDir = infer(cwd)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
