package server

import (
	"log"
	"testing"
)

func Test_inferRootDir(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inferRootDir()
			log.Println(rootDir)
		})
	}
}
