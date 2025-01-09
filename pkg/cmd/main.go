package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/MenD32/Tempest/pkg/client"
)

func main() {

	var wg sync.WaitGroup

	const duration time.Duration = 100
	const rps time.Duration = 20

	stop_chan := time.After(duration * time.Second)

	stop_chan_send := make(chan bool)

	defer close(stop_chan_send)

	fmt.Printf("channel size: %d\n", int(rps*duration))

	write_channel := make(chan *client.ParsedResponse, int(rps*duration))

	req := client.Request{}

	fmt.Printf("Starting\n")
	fmt.Printf("%s\n", time.Now().Format("2006-01-02T15:04:05.000Z07:00"))
	fmt.Printf("Sending requests\n")
	ticker := time.NewTicker(time.Second / rps)

	go func() {
		for {
			select {
			case <-ticker.C:
				wg.Add(1)
				go req.Send(write_channel, &wg)
			case <-stop_chan_send:
				return
			}
		}
	}()

	<-stop_chan

	ticker.Stop()
	stop_chan_send <- true

	fmt.Printf("Stopping\n")
	wg.Wait()
	close(write_channel)
	fmt.Printf("Stopped\n")

	file, err := os.OpenFile("/Users/mend/Misc/Tempest/temp/responses.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	fmt.Println("Writing to file")
	fmt.Printf("write_channel length: %d\n", len(write_channel))

	for msg := range write_channel {
		metrics := msg.ParseMetrics()

		output := client.Output{
			Timestamp: msg.Timestamp,
			Metrics:   *metrics,
		}
		outputJson, err := json.Marshal(output)
		if err != nil {
			fmt.Println("Error marshalling metrics:", err)
			continue
		}

		if _, err := file.WriteString(string(outputJson) + "\n"); err != nil {
			fmt.Println("Error writing to file:", err)
		}
	}
}
