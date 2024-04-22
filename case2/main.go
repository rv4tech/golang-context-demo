package main

import (
	"context"
	"fmt"
	"time"
)

const (
	ParentContextTimeOut = 15
	SubContextTimeOut    = 5
)

// Parent function that start main context.
// This is where the countdown starts.
func parentGoRoutine() {
	// Parent goroutine runs a command that needs to be time gated
	fmt.Printf("started parentGoRoutine with timeout %v\n", ParentContextTimeOut)
	// Timegate created for this particular goroutine
	ctx, cancel := context.WithTimeout(context.Background(), ParentContextTimeOut*time.Second)
	defer cancel()

	// Start sub goroutine and pass parent context as a parameter
	go subGoRoutine(ctx)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("parentGoRoutine ended: ", ctx.Err())
			return
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func subGoRoutine(parentCtx context.Context) {
	// Add extra time before we start context
	// to emulate subcontext strating late
	time.Sleep(8 * time.Second)

	fmt.Printf("started subGoRoutine with timeout %v\n", SubContextTimeOut)
	// This context inherits parent's context
	// Which means it should comlpete the job in that time window `ParentContextTimeOut`

	// If sub goroutine can not make it in time
	// it gets a `context deadline exceeded` error and shuts down
	// but parent's goroutine continues to go on

	ctx, cancel := context.WithTimeout(parentCtx, SubContextTimeOut*time.Second)
	defer cancel()

	err := slowDBQuery(ctx)
	fmt.Println("subGoRoutine ended: ", err)
}

// Note: no cancel func here, we need only ctx
func slowDBQuery(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func main() {
	go parentGoRoutine()
	for i := 1; i <= ParentContextTimeOut+5; i++ {
		time.Sleep(1 * time.Second)
		fmt.Printf("Time passed: %vs\t\n", i)
	}
}
