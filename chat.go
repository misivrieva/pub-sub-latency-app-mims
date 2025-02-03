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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "ably-latency",
	Short: "Measure latency of Ably Realtime messages",
	Run:   run,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().String("channel", "latency-test", "Channel name to use")
	rootCmd.PersistentFlags().Int("messages", 3, "Number of initial messages to send")
	rootCmd.PersistentFlags().Int("interval", 5, "Interval between sending messages in seconds")
	rootCmd.PersistentFlags().Int("listen", 30, "Duration to listen for responses in seconds")
}

func initConfig() {
	viper.AutomaticEnv()
	viper.BindPFlag("channel", rootCmd.PersistentFlags().Lookup("channel"))
	viper.BindPFlag("messages", rootCmd.PersistentFlags().Lookup("messages"))
	viper.BindPFlag("interval", rootCmd.PersistentFlags().Lookup("interval"))
	viper.BindPFlag("listen", rootCmd.PersistentFlags().Lookup("listen"))
}

func run(cmd *cobra.Command, args []string) {
	apiKey := viper.GetString("ABLY_API_KEY")
	if apiKey == "" {
		log.Fatal("ABLY_API_KEY environment variable is required")
	}

	channelName := viper.GetString("channel")
	numMessages := viper.GetInt("messages")
	interval := viper.GetInt("interval")

	client, err := ably.NewRealtime(ably.WithKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating Ably client: %v", err)
	}
	defer client.Close()

	channel := client.Channels.Get(channelName)

	// Start publishing messages
	go publishMessages(channel, numMessages, interval)

	// Start listening for messages
	subscribe(channel)

	// Keep the program running
	select {}
}

// Message represents the initial message sent by an instance.
type Message struct {
	SenderID  string  `json:"sender_id"`
	Timestamp float64 `json:"timestamp"`
}

// Response represents the response message sent by a receiver.
type Response struct {
	OriginalMessage Message `json:"original_message"`
	ReceiverID      string  `json:"receiver_id"`
	ReceivedAt      float64 `json:"received_at"`
}

var latencies = make(map[string][]float64)
var mu sync.RWMutex

func main() {
	fmt.Print("Type your username: ")
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	apiKey := os.Getenv("ABLY_API_KEY")
	if apiKey == "" {
		log.Fatal("ABLY_API_KEY environment variable is required")
	}

	client, err := ably.NewRealtime(ably.WithKey(apiKey), ably.WithClientID(username))
	if err != nil {
		log.Fatalf("Error creating Ably client: %v", err)
	}
	defer client.Close()

	channel := client.Channels.Get("chat")

	// Generate a unique ID for this instance
	instanceID := fmt.Sprintf("instance-%d", time.Now().UnixNano())

	// Subscribe to the channel
	_, err = channel.SubscribeAll(context.Background(), func(msg *ably.Message) {
		handleMessage(channel, instanceID, msg)
	})
	if err != nil {
		log.Fatalf("Error subscribing to channel: %v", err)
	}

	// Start publishing messages
	go publishing(channel)

	// Keep the program running
	select {}
}

func handleMessage(channel *ably.RealtimeChannel, instanceID string, msg *ably.Message) {
	switch msg.Name {
	case "message":
		var receivedMsg Message
		data, _ := msg.Data.([]byte)
		if err := json.Unmarshal(data, &receivedMsg); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			return
		}

		// Send a response immediately
		response := Response{
			OriginalMessage: receivedMsg,
			ReceiverID:      instanceID,
			ReceivedAt:      float64(time.Now().UnixNano()) / 1e9,
		}

		responseData, _ := json.Marshal(response)
		if err := channel.Publish(context.Background(), "response", responseData); err != nil {
			log.Printf("Error publishing response: %v", err)
		}

	case "response":
		var response Response
		data, _ := msg.Data.([]byte)
		if err := json.Unmarshal(data, &response); err != nil {
			log.Printf("Error unmarshalling response: %v", err)
			return
		}

		latency := response.ReceivedAt - response.OriginalMessage.Timestamp
		key := fmt.Sprintf("%s -> %s", response.OriginalMessage.SenderID, response.ReceiverID)

		mu.Lock()
		latencies[key] = append(latencies[key], latency)
		mu.Unlock()

		fmt.Printf("Latency for %s: %.6f seconds\n", key, latency)
	}
}

func subscribe(channel *ably.RealtimeChannel) {
	_, err := channel.SubscribeAll(context.Background(), func(msg *ably.Message) {
		fmt.Printf("Received message from %v: '%v'\n", msg.ClientID, msg.Data)
	})
	if err != nil {
		log.Printf("Error subscribing to channel: %v", err)
	}
}

func publishing(channel *ably.RealtimeChannel) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter message: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		err := channel.Publish(context.Background(), "message", text)
		if err != nil {
			log.Printf("Error publishing to channel: %v", err)
		}
	}
}

func publishMessages(channel *ably.RealtimeChannel, numMessages int, interval int) {
	instanceID := fmt.Sprintf("instance-%d", time.Now().UnixNano())

	for i := 0; i < numMessages; i++ {
		msg := Message{
			SenderID:  instanceID,
			Timestamp: float64(time.Now().UnixNano()) / 1e9,
		}

		msgData, _ := json.Marshal(msg)
		if err := channel.Publish(context.Background(), "message", msgData); err != nil {
			log.Printf("Error publishing message: %v", err)
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}
}
