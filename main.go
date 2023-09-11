package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mykhani/rate_limiter/limiter"
)

func main() {
	c := make(chan os.Signal)
	quit := false

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("Quit signal received!")
		quit = true
	}()

	bucket := limiter.NewTokenBucket(4, 4, 1000)

	defer func() {
		bucket.Close()
	}()

	for quit != true {
		fmt.Printf("quit: %t, %+v\n", quit, bucket)
		time.Sleep(time.Duration(1) * time.Second)

		// request comes
		// check if token available
		available := bucket.GetToken()
		// handle request or reject it
		if available {
			fmt.Println("Request can be handled")
		} else {
			fmt.Println("Request rate-limited")
		}
	}

	fmt.Println("Interupt received, quiting..!")

}
