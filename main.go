package main

import (
	"fmt"
	"time"
)

// go build -o testgo main.go
// nohup ./testgo &
func main() {
	for {
		fmt.Println(time.Now())
		time.Sleep(time.Second)
	}
}
