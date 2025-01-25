package main

import (
	"fmt"
	"image/color"
	"strings"
	"sync"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axmdb"
	"github.com/Cacsjep/goxis/pkg/axoverlay"
)

// MdbSceneOverlayApp represents an application that manages scene metadata overlays.
// It contains references to an AcapApplication, an OverlayProvider, and an MDBProvider
// for scene descriptions. It also maintains a list of observations, a channel for
// signaling closure, and a wait group for synchronizing goroutines.
type MdbSceneOverlayApp struct {
	app             *acapapp.AcapApplication
	overlayProvider *axoverlay.OverlayProvider
	mdbProvider     *axmdb.MDBProvider[axmdb.SceneDescription]
	mdbObservation  []axmdb.Observation
	closeChan       chan struct{}
	wg              sync.WaitGroup
}

// newMdbSceneOverlayApp creates a new instance of MdbSceneOverlayApp.
// It initializes the ACAP application instance, message broker provider, and overlay provider.
// It also sets up the necessary cleanup functions to be called on application close.
//
// Returns:
//   - *MdbSceneOverlayApp: A pointer to the newly created MdbSceneOverlayApp instance.
//   - error: An error if there was a problem creating the mdb provider or setting up the overlay.
func newMdbSceneOverlayApp() (*MdbSceneOverlayApp, error) {
	var err error

	// Create a new ACAP application instance
	msoa := &MdbSceneOverlayApp{
		app:            acapapp.NewAcapApplication(),
		mdbObservation: []axmdb.Observation{},
		closeChan:      make(chan struct{}),
	}

	msoa.app.AddCloseCleanFunc(msoa.Close)

	// Create the message broker provider
	msoa.mdbProvider, err = axmdb.NewMDBProvider[axmdb.SceneDescription]("1")
	if err != nil {
		return nil, fmt.Errorf("Failed to create mdb provider: %s", err.Error())
	}

	// Add close clean function to disconnect the provider
	msoa.app.AddCloseCleanFunc(msoa.mdbProvider.Disconnect)

	if err = msoa.SetupOverlay(); err != nil {
		return nil, fmt.Errorf("Failed to create overlay provider: %s", err.Error())
	}

	// Add close clean function to cleanup the overlay provider
	msoa.app.AddCloseCleanFunc(msoa.overlayProvider.Cleanup)

	return msoa, nil
}

// SetupOverlay initializes the overlay provider for the MdbSceneOverlayApp instance.
// It creates a new overlay provider, adds a cleanup function to the app, adds an overlay,
// and performs an initial redraw.
//
// Returns an error if any of these steps fail.
func (msoa *MdbSceneOverlayApp) SetupOverlay() error {
	var err error

	if msoa.overlayProvider, err = axoverlay.NewOverlayProvider(msoa.renderCallback, nil, nil); err != nil {
		return fmt.Errorf("Failed to create overlay provider: %s", err.Error())
	}

	// Add close clean function to cleanup the overlay provider
	msoa.app.AddCloseCleanFunc(msoa.overlayProvider.Cleanup)

	// we are here using no user data, because we can use golang structs to access the MdbSceneOverlayApp instance
	if _, err := msoa.overlayProvider.AddOverlay(axoverlay.NewAnchorCenterRrgbaOverlay(axoverlay.AxOverlayCustomNormalized, nil)); err != nil {
		return fmt.Errorf("Failed to add overlay: %s", err.Error())
	}

	// inital redraw
	if err = msoa.overlayProvider.Redraw(); err != nil {
		return fmt.Errorf("Failed todo a inital redraw overlay: %s", err.Error())
	}

	return nil
}

// renderCallback is a method of MdbSceneOverlayApp that handles the rendering of overlay events.
// It draws transparent background and bounding boxes around observations with a class score above a threshold.
//
// Parameters:
// - renderEvent: A pointer to axoverlay.OverlayRenderEvent which contains the rendering context and stream information.
//
// The method performs the following steps:
// 1. Draws a transparent background using the dimensions of the stream.
// 2. Iterates over the observations in msoa.mdbObservation.
// 3. Skips observations that do not have a class or have a class score below 0.1.
// 4. Normalizes the bounding box coordinates based on the stream dimensions.
// 5. Draws a bounding box around the observation with a label indicating the class type and score.
func (msoa *MdbSceneOverlayApp) renderCallback(renderEvent *axoverlay.OverlayRenderEvent) {
	renderEvent.CairoCtx.DrawTransparent(renderEvent.Stream.Width, renderEvent.Stream.Height)
	for _, obs := range msoa.mdbObservation {

		if obs.Class == nil || obs.Class.Score < 0.1 {
			// we are not interested in observations without class
			continue
		}

		x, y, w, h := BoxNormalize(&obs.BoundingBox, float64(renderEvent.Stream.Width), float64(renderEvent.Stream.Height))

		renderEvent.CairoCtx.DrawBoundingBox(
			x,
			y,
			w,
			h,
			BoxColor(obs.Class.Type),
			fmt.Sprintf("%s %d%%", strings.ToUpper(obs.Class.Type), int(obs.Class.Score*100)),
			axoverlay.ColorWite,
			17,
			"sans",
			0,
		)
	}
}

// MdbOnMetaDataWorker is a goroutine that listens for metadata updates and errors from the MDB provider.
// It handles different types of errors by logging them to the system log and processes incoming messages
// to update the overlay with new observations. The function runs in an infinite loop until it receives
// a signal to close via the closeChan channel.
func (msoa *MdbSceneOverlayApp) MdbOnMetaDataWorker() {
	msoa.wg.Add(1)
	defer msoa.wg.Done()
	for {
		select {
		case <-msoa.closeChan:
			return
		case err := <-msoa.mdbProvider.ErrorChan:
			if err == nil {
				msoa.app.Syslog.Errorf("Received nil error from mdb provider")
			}
			switch err.ErrType {
			// Happens on connecting
			case axmdb.MDBProviderErrorTypeConnection:
				msoa.app.Syslog.Critf("Connection error: %s", err.Err)
			// Happens on creating subscriber config
			case axmdb.MDBProviderErrorTypeSubscriberConfigCreate:
				msoa.app.Syslog.Critf("Subscriber config create error: %s", err.Err)
			// Happens on creating subscriber
			case axmdb.MDBProviderErrorTypeSubscriberCreate:
				msoa.app.Syslog.Critf("Subscriber create error: %s", err.Err)
			// Happens when the message is not T axmdb.MessageType
			case axmdb.MDBProviderErrorTypeInvalidMessage:
				msoa.app.Syslog.Critf("Invalid message error: %s", err.Err)
			// Happens when recv mdb message is emtpy
			case axmdb.MDBProviderErrorTypeEmptyPayload:
				msoa.app.Syslog.Errorf("Empty payload error: %s", err.Err)
			// Happens when the message could not parsed from json -> T axmdb.MessageType
			case axmdb.MDBProviderErrorTypeParseMessage:
				msoa.app.Syslog.Errorf("Parse message error: %s", err.Err)
			// Should not happen just to be sure
			default:
				msoa.app.Syslog.Critf("Unknown error: %s", err.Err)
			}
		case msg := <-msoa.mdbProvider.MessageChan:

			msoa.mdbObservation = []axmdb.Observation{}

			for _, mdo := range msg.Frame.Observations {
				msoa.mdbObservation = append(msoa.mdbObservation, mdo)
			}
			if err := msoa.overlayProvider.Redraw(); err != nil {
				msoa.app.Syslog.Errorf("Failed to redraw overlay: %s", err.Error())
			}
		}
	}
}

// Run starts the MdbSceneOverlayApp by launching the metadata worker in a separate goroutine,
// connecting the mdb provider, and running the acap application.
func (msoa *MdbSceneOverlayApp) Run() {
	go msoa.MdbOnMetaDataWorker()

	// Connect the mdb provider
	msoa.mdbProvider.Connect()

	// Run the acap application
	msoa.app.Run()
}

// Close gracefully shuts down the MdbSceneOverlayApp by closing the closeChan
// channel and waiting for all goroutines in the wait group to finish.
func (msoa *MdbSceneOverlayApp) Close() {
	close(msoa.closeChan)
	msoa.wg.Wait()
}

// BoxNormalize scales the normalized coordinates of a bounding box to the given width and height.
//
// Parameters:
//   - bbox: A pointer to an axmdb.Box struct containing the normalized coordinates of the bounding box.
//   - width: The width to scale the bounding box to.
//   - height: The height to scale the bounding box to.
//
// Returns:
//   - rectX: The scaled x-coordinate of the left side of the bounding box.
//   - rectY: The scaled y-coordinate of the top side of the bounding box.
//   - rectWidth: The scaled width of the bounding box.
//   - rectHeight: The scaled height of the bounding box.
func BoxNormalize(bbox *axmdb.Box, width, height float64) (rectX, rectY, rectWidth, rectHeight float64) {
	// Scale the normalized coordinates
	left := bbox.Left * width
	right := bbox.Right * width
	top := bbox.Top * height
	bottom := bbox.Bottom * height
	return left, top, right - left, bottom - top
}

// BoxColor returns the color associated with a given class.
// The function maps specific classes to predefined colors:
// "Human" to green, "Car" to blue, and "Face" to deep orange.
// If the class is not found in the map, it returns red.
//
// Parameters:
//
//	class (string): The class name to get the color for.
//
// Returns:
//
//	color.RGBA: The color associated with the given class.
func BoxColor(class string) color.RGBA {
	color_map := map[string]color.RGBA{
		"Human": axoverlay.ColorMaterialGreen,
		"Car":   axoverlay.ColorMaterialBlue,
		"Face":  axoverlay.ColorMaterialDeepOrange,
	}

	if c, ok := color_map[class]; ok {
		return c
	}
	return axoverlay.ColorMaterialRed
}
