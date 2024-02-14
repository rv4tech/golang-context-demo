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

// Shortcut for time.Sleep()
var wait = func(n time.Duration) {
	time.Sleep(n * time.Second)
}

// Runtime visual representation
var printWithSleep = func() {
	for i := 1; i < ParentContextTimeOut+5; i++ {
		wait(1)
		fmt.Printf("Time passed: %vs\t\n", i)
	}
}

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
			wait(1)
		}
	}
}

func subGoRoutine(parentCtx context.Context) {
	// Add extra time before we start context
	// to emulate subcontext strating late
	wait(8)

	fmt.Printf("started subGoRoutine with timeout %v\n", SubContextTimeOut)
	// This context inherits parent's context
	// Which means it should comlpete the job in that time window `ParentContextTimeOut`
	// If sub goroutine can not make it in time of parent's goroutine
	// It gets a `context deadline exceeded` error and shuts down
	// But main goroutine continues to go on

	ctx, cancel := context.WithTimeout(parentCtx, SubContextTimeOut*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("subGoRoutine ended: ", ctx.Err())
			return
		default:
			wait(1)
		}
	}
}

func main() {
	go parentGoRoutine()
	printWithSleep()
}
