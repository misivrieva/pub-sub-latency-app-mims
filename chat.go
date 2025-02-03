package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ably/ably-go/ably"
)

// Message represents the initial message sent by an instance.
type Message struct {
	SenderID  string  `json:"sender_id"` // Unique identifier for the sender
	Timestamp float64 `json:"timestamp"` // Timestamp in seconds with microsecond precision
}

// Response represents the response message sent by a receiver.
type Response struct {
	OriginalMessage Message `json:"original_message"` // Original message details
	ReceiverID      string  `json:"receiver_id"`      // Unique identifier for the receiver
	ReceivedAt      float64 `json:"received_at"`      // Timestamp when the message was received
}

// LatencyStats holds latency statistics for a pair of instances.
type LatencyStats struct {
	Min, Max, Avg float64
}

var latencies = make(map[string][]float64)
var mu sync.RWMutex

func main() {
	fmt.Println("Type your username")
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	username = strings.Replace(username, "\n", "", -1)

	// Connect to Ably using the API key and ClientID specified above
	client, err := ably.NewRealtime(
		ably.WithKey("02uj7w.guUaIg:zWMaFId7FlDQysfzRUph2Q2Nu9a1FhrnbRHjpdAK33M"))
	//ably.WithEchoMessages(false), // Uncomment to stop messages you send from being sent back
	ably.WithClientID(username)
	if err != nil {
		panic(err)
	}
	// Set up connection events handler.
	client.Connection.OnAll(func(change ably.ConnectionStateChange) {
		fmt.Printf("Connection event: %s state=%s reason=%s", change.Event, change.Current, change.Reason)
	})
	// Then connect.
	client.Connect()
	// Connect to the Ably Channel with name 'chat'
	channel := client.Channels.Get("chat")

	// Generate a unique ID for this instance
	instanceID := fmt.Sprintf("instance-%d", time.Now().UnixNano())

	// Subscribe to the channel
	// sub, err := channel.SubscribeAll(context.Background(), func(msg *ably.Message) {
	// 	subscribe(channel)
	// })
	sub, err := channel.SubscribeAll(context.Background(), func(msg *ably.Message) {
		subscribe(channel)
	})
	if err != nil {
		log.Fatalf("Error subscribing to channel: %v", err)
	}
	//Keep tracck of state, you need to be attached to send messages
	channel.OnAll(func(change ably.ChannelStateChange) {
		fmt.Printf("Channel event event: %s channel=%s state=%s reason=%s", channel.Name, change.Event, change.Current, change.Reason)
	})
	// Start the goroutine to allow for publishing messages
	publishing(channel)
	// Goroutine to send messages at regular intervals
	go func() {
		for i := 0; i < 3; i++ { // Send 3 messages by default
			msg := Message{
				SenderID:  instanceID,
				Timestamp: float64(time.Now().UnixNano()) / 1e9, // Current time in seconds with microsecond precision
			}

			// Publish the message to the channel
			if err := channel.Publish("message", msg); err != nil {
				log.Printf("Error publishing message: %v", err)
			}

			time.Sleep(5 * time.Second) // Wait 5 seconds between messages
		}
	}()

// Goroutine to listen for messages and responses
go func() {
	//That's too complex and there should be an easier way
	for msg := range sub.Channel.Messages.All().Get(context.Background()) {}
	switch msg.Type {
		case "message":
			// Handle incoming message
			var receivedMsg Message
			if err := json.Unmarshal(msg.Data, &receivedMsg); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}

			// Send a response immediately
			response := Response{
				OriginalMessage: receivedMsg,
				ReceiverID:      instanceID,
				ReceivedAt:      float64(time.Now().UnixNano()) / 1e9,
			}
			if err := sub.Channel.Publish("response", response); err != nil {
				log.Printf("Error publishing response: %v", err)
			}

		case "response":
			// Handle incoming response
			var response Response
			if err := json.Unmarshal(msg.Data, &response); 
			err != nil {
				log.Printf("Error unmarshalling response: %v", err)}		
	}}

			// Calculate latency
			latency := response.ReceivedAt - response.OriginalMessage.Timestamp
			key := fmt.Sprintf("%s -> %s", response.OriginalMessage.SenderID, response.ReceiverID)

			// Store latency in the map
			mu.Lock()
			latencies[key] = append(latencies[key], latency)
			mu.Unlock()
		}
}
}()


func subscribe(channel *ably.RealtimeChannel) {
	// Subscribe to messages sent on the channel
	_, err := channel.SubscribeAll(context.Background(), func(msg *ably.Message) {
		fmt.Printf("Received message from %v: '%v'\n", msg.ClientID, msg.Data)
	})
	if err != nil {
		err := fmt.Errorf("subscribing to channel: %w", err)
		fmt.Println(err)
	}
}

func publishing(channel *ably.RealtimeChannel) {
	reader := bufio.NewReader(os.Stdin)

	for {
		text, _ := reader.ReadString('\n')
		text = strings.ReplaceAll(text, "\n", "")
		// Publish the message typed in to the Ably Channel
		err := channel.Publish(context.Background(), "message", text)
		// await confirmation that message was received by Ably
		if err != nil {
			err := fmt.Errorf("publishing to channel: %w", err)
			fmt.Println(err)
		}
	}
}

// calculateStats calculates the minimum, maximum, and average latency.
func calculateStats(latencies []float64) (min, max, avg float64) {
	if len(latencies) == 0 {
		return 0, 0, 0
	}

	min = latencies[0]
	max = latencies[0]
	sum := 0.0

	for _, lat := range latencies {
		if lat < min {
			min = lat
		}
		if lat > max {
			max = lat
		}
		sum += lat
	}

	avg = sum / float64(len(latencies))
	return min, max, avg
}
