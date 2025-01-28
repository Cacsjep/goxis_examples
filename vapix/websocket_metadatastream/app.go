package main

import (
	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/vapix"
)

// This example demonstrates how to use the vapix package to obtain the metadata stream via websocket
// https://help.axis.com/en-us/axis-os-knowledge-base#metadata-via-websocket
func main() {

	// Initialize a new ACAP application instance. You could do it also without it,
	// but we use here the syslog and on close cleanup from acapapp
	app := acapapp.NewAcapApplication()

	// subscribe to all, by passing nil
	wsc := vapix.NewVapixWsMetadataConsumer(nil)

	// subscribe to specific events
	// wsc := vapix.NewVapixWsMetadataConsumer(
	// 	&[]vapix.VapixWsMetadataStreamRequestEventFilter{
	// 		{
	// 			TopicFilter:   "tns1:Device/tnsaxis:IO/VirtualInput",
	// 			ContentFilter: "boolean(//SimpleItem[@Name=\"port\" and @Value=\"1\"])",
	// 		},
	// 	},
	// )

	// connect to the websocket
	conn, err := wsc.Connect()
	if err != nil {
		app.Syslog.Critf("Failed to connect to WebSocket: %s", err.Error())
	}

	// add close cleanup function, signal handler will call this function when the application is closed
	app.AddCloseCleanFunc(func() {
		conn.Close()
	})

	// first message should be
	// 		{"apiVersion":"1.0","method":"events:configure","data":{}} what indicates that our request was accepted
	// otherwise first message is an error like
	//		{"apiVersion":"1.0","method":"events:configure","error":{"code":2104,"message":"Invalid event datasource payload: eventFilterList is empty"}}
	// second message is the actual event data, when our request was accepted
	//		{"apiVersion":"1.0","method":"events:notify","params":{"notification":{"topic":"tns1:Device/tnsaxis:IO/VirtualInput","timestamp":1737794941428,"message":{"source":{"port":"1"},"key":{},"data":{"active":"0"}}}}}

	for {
		resp := &vapix.VapixWsMetadataStreamResponse{}

		if err := conn.ReadJSON(resp); err != nil {
			app.Syslog.Errorf("Error reading message: %s", err.Error())
			// we can't continue without a correct json message, so we try again
			continue
		}

		// Check whether the response contains an error or a valid message.
		// An error in the response indicates an invalid request or configuration issue.
		// Since the application shouldn't proceed without a valid request, a critical error
		// is logged.
		if resp.Error != nil {
			app.Syslog.Critf("Received error: %v", resp.Error)
		}

		if resp.Method == "events:configure" {
			// If we recv a events configure message, it means that our request was accepted
			app.Syslog.Infof("Metadata Stream request was succesfull")
		} else if resp.Method == "events:notify" {
			// If we recv a events notify message, it means that we received the actual event data
			app.Syslog.Infof("Received notify: %v", resp.Params)
		} else {
			// If we recv a unknown message, we log it
			app.Syslog.Warnf("Received unknown: %v", resp)
		}
	}
}
