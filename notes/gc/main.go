package main

import (
	"crypto/rand"
	"fmt"
	"runtime"
)

const MiB = 1024 * 1024

func main() {
	sizes := []int{
		200 * MiB,
		180 * MiB,
		175 * MiB,
		140 * MiB,
		60 * MiB,
		170 * MiB,
		40 * MiB,
		120 * MiB,
		80 * MiB,
		150 * MiB,
	}

	for iter, size := range sizes {
		// 指定のサイズのランダムなバイト列を生成
		buf := make([]byte, size)
		_, _ = rand.Read(buf)

		// 計測
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		fmt.Printf(
			"iter %d HeapAlloc=%dMiB NextGC=%dMiB NumGC=%d\n",
			iter+1,
			m.HeapAlloc/MiB,
			m.NextGC/MiB,
			m.NumGC,
		)
	}
}
