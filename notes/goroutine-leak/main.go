package main

import (
	"context"
	"fmt"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func main() {
	ctx := context.Background()
	dumpGoroutines("before")
	leak2(ctx)
	dumpGoroutines("after")

}

func leak1() {
	ch := make(chan int)
	go func() {
		v := <-ch // 誰も送らない / close しない
		fmt.Println(v)
	}()
}

func leak2(ctx context.Context) {
	go func() {
		for {
			busy := make([]byte, 1024*1024) // 1MiB
			_ = busy
			time.Sleep(time.Second)
		}
	}()
}

func leak3() {
	go func() {
		t := time.NewTicker(time.Second)
		for range t.C {
			tick()
		}
	}()
}

func dumpGoroutines(tag string) {
	fmt.Printf("[%s] goroutines=%d\n", tag, runtime.NumGoroutine())
}
