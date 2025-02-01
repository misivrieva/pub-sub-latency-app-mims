package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ably/ably-go/ably"
)

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

	//Keep tracck of state, you need to be attached to send messages
	channel.OnAll(func(change ably.ChannelStateChange) {
		fmt.Printf("Channel event event: %s channel=%s state=%s reason=%s", channel.Name, change.Event, change.Current, change.Reason)
	})

	fmt.Println("You can now send messages!")
	// Start the goroutine to allow for subscribing
	subscribe(channel)
	// Start the goroutine to allow for publishing messages
	publishing(channel)
}

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
