package main

func main() {
	// Conc()
	Nursery()
}

func Scope() {

}

func Nursery() {
	// nursery.RunConcurrently(
	// 	func(ctx context.Context, errCh chan error) {
	// 		time.Sleep(100 * time.Millisecond)
	// 		fmt.Println("job 1 done")
	// 	},
	// 	func(ctx context.Context, errCh chan error) {
	// 		time.Sleep(100 * time.Millisecond)
	// 		fmt.Println("job 2 done")
	// 	},
	// )

	// nursery.RunConcurrently(
	// 	func(ctx context.Context, errCh chan error) {
	// 		fmt.Println("outer job 1")
	// 		nursery.RunConcurrentlyWithContext(ctx,
	// 			func(ctx context.Context, errCh chan error) {
	// 				fmt.Println("inner job 1-1")
	// 			},
	// 			func(ctx context.Context, errCh chan error) {
	// 				fmt.Println("inner job 1-2")
	// 			},
	// 		)
	// 	},
	// 	func(ctx context.Context, errCh chan error) {
	// 		fmt.Println("outer job 2")
	// 	},
	// )

	// err := nursery.RunConcurrentlyWithTimeout(500*time.Millisecond,
	// 	func(ctx context.Context, errCh chan error) {
	// 		select {
	// 		case <-time.After(1 * time.Second):
	// 			fmt.Println("job 1 done")
	// 		case <-ctx.Done():
	// 			fmt.Println("job 1 timed out")
	// 		}
	// 	},
	// 	func(ctx context.Context, errCh chan error) {
	// 		fmt.Println("job 2 done")
	// 	},
	// )
	// fmt.Println("err:", err)

	// nursery.RunUntilFirstCompletion(
	// 	func(ctx context.Context, errCh chan error) {
	// 		select {
	// 		case <-time.After(2 * time.Second):
	// 			fmt.Println("job A done")
	// 		case <-ctx.Done():
	// 			fmt.Println("job A cancelled")
	// 		}
	// 	},
	// 	func(ctx context.Context, errCh chan error) {
	// 		time.Sleep(500 * time.Millisecond)
	// 		fmt.Println("job B done first") // これが先に終わる
	// 	},
	// )

	// nursery.RunUntilFirstCompletionWithTimeout(500*time.Millisecond,
	// 	func(ctx context.Context, errCh chan error) {
	// 		select {
	// 		case <-time.After(2 * time.Second):
	// 			fmt.Println("job A done")
	// 		case <-ctx.Done():
	// 			fmt.Println("job A cancelled")
	// 		}
	// 	},
	// 	func(ctx context.Context, errCh chan error) {
	// 		select {
	// 		case <-time.After(2 * time.Second):
	// 			fmt.Println("job B done")
	// 		case <-ctx.Done():
	// 			fmt.Println("job B cancelled")
	// 		}
	// 	},
	// )
}

// https://zenn.dev/kouichi_itagaki/articles/492a55e26c9d86　参考
func Conc() {
	// var wg conc.WaitGroup
	// wg.Go(func() {
	// 	time.Sleep(1 * time.Second)
	// 	fmt.Print("task 1")
	// })
	// wg.Go(func() {
	// 	time.Sleep(1 * time.Second)
	// 	fmt.Print("task 2")
	// })
	// wg.Go(func() {
	// 	time.Sleep(1 * time.Second)
	// 	fmt.Print("task 3")
	// })
	// wg.Wait()

	// p := pool.New().WithErrors().WithMaxGoroutines(2)

	// for i := range 5 {
	// 	i := i
	// 	p.Go(func() error {
	// 		time.Sleep(500 * time.Millisecond)
	// 		if i == 4 {
	// 			return fmt.Errorf("task %d failed", i)
	// 		}
	// 		fmt.Printf("task %d done\n", i)
	// 		return nil
	// 	})
	// }

	// if err := p.Wait(); err != nil {
	// 	fmt.Println("error:", err)
	// }

	// p := pool.NewWithResults[int]().WithErrors()

	// for i := range 5 {
	// 	i := i
	// 	p.Go(func() (int, error) {
	// 		time.Sleep(2 * time.Second)
	// 		return i * i, nil
	// 	})
	// }

	// results, err := p.Wait()
	// if err != nil {
	// 	fmt.Println("error:", err)
	// 	return
	// }
	// fmt.Println("results:", results)

	// iter.ForEach([]int{1, 2, 3, 4, 5}, func(v *int) {
	// 	time.Sleep(2 * time.Second)
	// 	fmt.Printf("processing %d\n", *v)
	// })

	// s := stream.New()

	// for i := range 5 {
	// 	i := i
	// 	s.Go(func() stream.Callback {
	// 		time.Sleep(time.Second)
	// 		// ここは並列実行
	// 		result := i * i

	// 		// 返すCallbackは投入順に順次実行される
	// 		return func() {
	// 			fmt.Printf("result %d\n", result)
	// 		}
	// 	})
	// }

	// s.Wait()
}
