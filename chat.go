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

	fmt.Println("You can now send messages!")

	// Connect to the Ably Channel with name 'chat'
	channel := client.Channels.Get("chat")

	// Enter the Presence set of the channel
	channel.Presence.Enter(context.Background(), "")

	getHistory(channel)

	subscribe(channel)

	subscribePresence(channel)

	// Start the goroutine to allow for publishing messages
	publishing(channel)
}

func getHistory(channel *ably.RealtimeChannel) {
	// Before subscribing for messages, check the channel's
	// History for any missed messages. By default a channel
	// will keep 2 minutes of history available, but this can
	// be extended to 48 hours
	pages, err := channel.History().Pages(context.Background())
	if err != nil || pages == nil {
		return
	}

	hasHistory := true

	for ; hasHistory; hasHistory = pages.Next(context.Background()) {
		for _, msg := range pages.Items() {
			fmt.Printf("Previous message from %v: '%v'\n", msg.ClientID, msg.Data)
		}
	}
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

func subscribePresence(channel *ably.RealtimeChannel) {
	// Subscribe to presence events (people entering and leaving) on the channel
	_, pErr := channel.Presence.SubscribeAll(context.Background(), func(msg *ably.PresenceMessage) {
		if msg.Action == ably.PresenceActionEnter {
			fmt.Printf("%v has entered the chat\n", msg.ClientID)
		} else if msg.Action == ably.PresenceActionLeave {
			fmt.Printf("%v has left the chat\n", msg.ClientID)
		}
	})
	if pErr != nil {
		err := fmt.Errorf("subscribing to presence in channel: %w", pErr)
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
