package main

import (
	"github.com/techplexengineer/frc-networktables-go"
	"time"
)

func main() {
	client, err := frcntgo.InitClient()
	if err != nil {
		panic(err)
	}

	time.Sleep(5)

	_ = client.GetBoolean("/FMSInfo/IsRedAlliance")
	//fmt.Printf("is Red: %v\n", isRed)

	for {
		time.Sleep(5)
	}
}
