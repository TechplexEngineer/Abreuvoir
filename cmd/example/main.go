package main

import (
	"fmt"
	"github.com/techplexengineer/frc-networktables-go"
	"time"
)

func main() {

	client, err := frcntgo.NewClient("0.0.0.0", "1735")
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second) //@todo really need to wait for the client to enter the sync state

	isRed, err := client.GetBoolean("/bool")
	if err != nil {
		panic(err)
	}
	fmt.Printf("is Red: %v\n", isRed)

	for {
		time.Sleep(5 * time.Millisecond)
	}
}
