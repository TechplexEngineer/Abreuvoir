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

	fmt.Printf("Waiting for initial sync...\n")
	for client.GetStatus() != frcntgo.ClientInSync {
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Printf("Initial sync complete.\n")

	//time.Sleep(5 * time.Second) //@todo really need to wait for the client to enter the sync state

	isRed, err := client.GetBoolean("/bool")
	if err != nil {
		panic(err)
	}
	fmt.Printf("is Red: %v\n", isRed)

	for {
		time.Sleep(5 * time.Millisecond)
	}
}
