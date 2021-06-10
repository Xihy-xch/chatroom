package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestGenUserID(t *testing.T) {
	tests := []struct {
		name string
		want int64
	}{
		// TODO: Add test cases.
		{
			name: "测试atomic",
			want: int64(1),
		},
		{
			name: "测试atomic",
			want: int64(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenUserID(); got != tt.want {
				t.Errorf("GenUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkAtomicSpeed(b *testing.B) {
	a := int64(0)
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&a, 1)
		}()
	}
	wg.Wait()
}

func BenchmarkMutex(b *testing.B) {
	a := int64(0)
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mutex.Lock()
			defer mutex.Unlock()
			a++
		}()
	}
	wg.Wait()
}
