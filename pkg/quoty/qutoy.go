package quoty

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"
	"os"
)

// handleSendMessage is used to take care of the AppMentionEvent when the bot is mentioned
func handleSendMessage(channel string, client *slack.Client) error {

	msg := BuildQuote(true, nil)

	// The Channel is available in the event message
	_, _, _, err := client.SendMessage(channel, msg)

	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}

func handleModifyMessage(quoteId *string, channel string, msgTimeStamp string, client *slack.Client) error {

	msg := BuildQuote(false, quoteId)

	// The Channel is available in the event message
	_, _, _, err := client.UpdateMessage(channel, msgTimeStamp, msg)

	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}

// handleEventMessage will take an event and handle it properly based on the type of event
func handleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client) error {
	switch event.Type {
	// First we check if this is an CallbackEvent
	case slackevents.CallbackEvent:

		innerEvent := event.InnerEvent
		// Yet Another Type switch on the actual Data to see if its an AppMentionEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// The application has been mentioned since this Event is a Mention event
			err := handleSendMessage(ev.Channel, client)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported event type")
	}
	return nil
}

func StartQuoty() {

	// Load Env variables from .dot file
	godotenv.Load(".env")

	token := os.Getenv("SLACK_AUTH_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")
	// Create a new client to slack by giving token
	// Set debug to true while developing
	// Also add a ApplicationToken option to the client
	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))
	// go-slack comes with a SocketMode package that we need to use that accepts a Slack client and outputs a Socket mode client instead
	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		// Option to set a custom logger
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	// Create a context that can be used to cancel goroutine
	ctx, cancel := context.WithCancel(context.Background())
	// Make this cancel called properly in a real program , graceful shutdown etc
	defer cancel()

	go func(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
		// Create a for loop that selects either the context cancellation or the events incomming
		for {
			select {
			// inscase context cancel is called exit the goroutine
			case <-ctx.Done():
				log.Println("Shutting down socketmode listener")
				return
			case event := <-socketClient.Events:
				// We have a new Events, let's type switch the event
				// Add more use cases here if you want to listen to other events.
				switch event.Type {
				// handle EventAPI events
				case socketmode.RequestTypeInteractive:
					actionEvent, ok := event.Data.(slack.InteractionCallback)
					if !ok {
						log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
						continue
					}
					// We need to send an Acknowledge to the slack server
					socketClient.Ack(*event.Request)
					// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
					err := handleInteractiveEvent(actionEvent, client)
					if err != nil {
						// Replace with actual err handeling
						log.Fatal(err)
					}

				case socketmode.EventTypeEventsAPI:
					// The Event sent on the channel is not the same as the EventAPI events so we need to type cast it
					eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
						continue
					}
					// We need to send an Acknowledge to the slack server
					socketClient.Ack(*event.Request)
					// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
					err := handleEventMessage(eventsAPIEvent, client)
					if err != nil {
						// Replace with actual err handeling
						log.Fatal(err)
					}
				}

			}
		}
	}(ctx, client, socketClient)

	socketClient.Run()
}

func handleInteractiveEvent(event slack.InteractionCallback, client *slack.Client) interface{} {
	switch event.Type {
	// First we check if this is an CallbackEvent
	case slack.InteractionTypeBlockActions:

		switch event.ActionCallback.BlockActions[0].ActionID {
		case "press_one_more":
			err := handleSendMessage(event.Channel.Name, client)
			if err != nil {
				return errors.New("could not send new message")
			}
		case "press_no_more":
			quoteId := event.ActionCallback.BlockActions[0].Value
			err := handleModifyMessage(&quoteId, event.Channel.ID, event.Message.Msg.Timestamp, client)
			if err != nil {
				return errors.New("could not send new message")
			}
		}
	default:
		return errors.New("unsupported event type")
	}
	return nil
}
