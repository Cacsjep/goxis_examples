package main

import (
	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axmdb"
)

// This example demonstrates how to use the axis message broker api
//
// Orginal C Example:https://github.com/AxisCommunications/acap-native-sdk-examples/tree/main/message-broker/consume-scene-metadata
func main() {

	// Initialize a new ACAP application instance.
	// AcapApplication initializes the ACAP application with there name, eventloop, and syslog etc..
	app := acapapp.NewAcapApplication()

	defer app.Close()

	// Run the application using the C-like API
	// usingClikeApi(app)

	// Run the application using the provider
	usingProvider(app)

}

func usingClikeApi(app *acapapp.AcapApplication) {
	con, err := axmdb.MDBConnectionCreate(func(onErr error) {
		app.Syslog.Critf("Connection failed (onError Callback): %s", onErr.Error())
	})

	if err != nil {
		app.Syslog.Critf("Failed to create connection: %s", err.Error())
	}

	app.AddCloseCleanFunc(con.Destroy)

	sub_config, err := axmdb.MDBSubscriberConfigCreate("com.axis.analytics_scene_description.v0.beta", "1", func(msg *axmdb.Message) {
		app.Syslog.Infof("Received message: %s", msg.Payload)
	})

	app.AddCloseCleanFunc(sub_config.Destroy)

	subscriber, err := axmdb.MDBSubscriberCreateAsync(con, sub_config, func(onDone error) {
		if onDone != nil {
			app.Syslog.Critf("Subscriber failed: %s", onDone.Error())
		} else {
			app.Syslog.Infof("Subscriber created")
		}
	})

	app.AddCloseCleanFunc(subscriber.Destroy)

	// Run gmain loop with signal handler attached.
	// This will block the main thread until the application is stopped.
	// The application can be stopped by sending a signal to the process (e.g. SIGINT).
	app.Run()
}

func usingProvider(app *acapapp.AcapApplication) {
	provider, err := axmdb.NewMDBProvider[axmdb.SceneDescription]("1")
	if err != nil {
		app.Syslog.Critf("Failed to create provider: %s", err.Error())
	}
	app.AddCloseCleanFunc(provider.Disconnect)

	app.Syslog.Info("Provider created")

	go func() {
		for {
			select {
			case err := <-provider.ErrorChan:
				switch err.ErrType {
				// Happens on connecting
				case axmdb.MDBProviderErrorTypeConnection:
					app.Syslog.Critf("Connection error: %s", err.Err)
				// Happens on creating subscriber config
				case axmdb.MDBProviderErrorTypeSubscriberConfigCreate:
					app.Syslog.Critf("Subscriber config create error: %s", err.Err)
				// Happens on creating subscriber
				case axmdb.MDBProviderErrorTypeSubscriberCreate:
					app.Syslog.Critf("Subscriber create error: %s", err.Err)
				// Happens when the message is not T axmdb.MessageType
				case axmdb.MDBProviderErrorTypeInvalidMessage:
					app.Syslog.Critf("Invalid message error: %s", err.Err)
				// Happens when recv mdb message is emtpy
				case axmdb.MDBProviderErrorTypeEmptyPayload:
					app.Syslog.Errorf("Empty payload error: %s", err.Err)
				// Happens when the message could not parsed from json -> T axmdb.MessageType
				case axmdb.MDBProviderErrorTypeParseMessage:
					app.Syslog.Errorf("Parse message error: %s", err.Err)
				// Should not happen just to be sure
				default:
					app.Syslog.Critf("Unknown error: %s", err.Err)
				}
			case msg := <-provider.MessageChan:
				app.Syslog.Info(msg.String())
			}
		}
	}()

	// Connect it, any error are sent to the error channel
	provider.Connect()

	// we dont need really here the gmain loop but we have it already so why not use
	// just select{} would also be enough
	app.Run()
}
