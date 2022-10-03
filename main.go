package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// https://stackoverflow.com/questions/39776481/how-to-wait-until-buffered-channel-semaphore-is-empty
// https://stackoverflow.com/a/19892995
func kc() {

	i := 0
	ch := make(chan int)
	wg := &sync.WaitGroup{}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	defer wg.Wait()
	defer close(ch)

	go processPch(ch)
	for i < 3000000000 {
		time.Sleep(10 * time.Millisecond)
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			fmt.Printf("Closing due to signal\n")
			i = 3000000000
		default:
			wg.Add(1)
			fmt.Printf("i %+v\n", i)
			go aggMsg(i, ch, wg)
			i = i + 1
			// time.Sleep(1 * time.Second)
		}
	}
	// go processPch(ch)
	fmt.Printf("Waiting\n")
	wg.Wait()
	fmt.Printf("Closing\n")
}

func aggMsg(i int, ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(2 * time.Millisecond)
	fmt.Printf("Sending %+v\n", i*10)
	ch <- i * 10

}
func processPch(ch <-chan int) {
	for i := range ch {
		// i := <-ch
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("Received processPch %+v\n", i)

	}
}

func main() {
	fmt.Println("Hello World")
	kc()
}
