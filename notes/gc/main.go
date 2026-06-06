package main

import (
	"crypto/rand"
	"fmt"
	"runtime"
	"time"
)

const MiB = 1024 * 1024

func main() {
	for i := 1; i <= 10; i++ {
		var chunks [][]byte
		for j := 0; j < 300; j++ {
			c := make([]byte, MiB)
			_, _ = rand.Read(c)
			chunks = append(chunks, c)
		}
		time.Sleep(50 * time.Millisecond)
		printStats(i)
		runtime.KeepAlive(chunks)
	}
}

func printStats(iter int) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("iter %2d  HeapAlloc=%4dMiB  NextGC=%4dMiB  NumGC=%d\n",
		iter,
		m.HeapAlloc/MiB,
		m.NextGC/MiB,
		m.NumGC,
	)
}
