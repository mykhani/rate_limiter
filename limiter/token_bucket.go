package limiter

import (
	"fmt"
	"time"
)

type TokenBucket struct {
	done           chan bool
	capacity       int
	available      int
	ticker         *time.Ticker
	refillQuantity int
	refillRateMs   int
}

func (bucket *TokenBucket) addTokens(toAdd int) bool {
	space := bucket.capacity - bucket.available

	if space > 0 {
		if space < toAdd {
			bucket.available += space
			fmt.Printf("Added only %d tokens out of %d\n", space, toAdd)
		} else {
			bucket.available += toAdd
			fmt.Printf("Added %d tokens\n", toAdd)
		}
		return true
	}

	fmt.Printf("Bucket full!!\n")
	return false
}

func (bucket *TokenBucket) GetToken() bool {
	if bucket.available > 0 {
		bucket.available--
		return true
	}

	return false
}

func (filler *TokenBucket) start() {

	filler.ticker = time.NewTicker(time.Duration(filler.refillRateMs) * time.Millisecond)

	// ticker callback listener
	// see how to use contexts to kill this go-routine
	go func() {
		for {
			select {
			case <-filler.done:
				fmt.Println("Stopping the filler goroutine")
				return
			case <-filler.ticker.C:
				added := filler.addTokens(filler.refillQuantity)
				if !added {
					fmt.Println("Failed to add token, bucket full!!")
				}
			}
		}
	}()
}

func (filler *TokenBucket) Close() {
	filler.done <- true
}

func NewTokenBucket(capacity int, refillQuantity int, refillRateMs int) *TokenBucket {

	bucket := &TokenBucket{
		done:           make(chan bool),
		capacity:       capacity,
		available:      0,
		refillQuantity: refillQuantity,
		refillRateMs:   refillRateMs,
	}

	bucket.start()

	return bucket
}
