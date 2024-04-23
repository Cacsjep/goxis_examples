package main

import (
	"time"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axoverlay"
)

// This example demonstrate how to use overlay provider to draw png images on a stream.
//
// ! Note: Overlay callbacks only invoked when stream is viewed via web ui or rtsp etc..
var (
	app             *acapapp.AcapApplication
	overlayProvider *acapapp.OverlayProvider
	err             error
	image_seq       *ImageSequence
)

// streamSelectCallback can be used to select which streams to render overlays to.
// Note that YCBCR streams are always skipped since these are used for analytics.
// ! Just for demo demonstration
func streamSelectCallback(streamSelectEvent *axoverlay.OverlayStreamSelectEvent) bool {
	return true
}

// adjustmentCallback is called when an overlay needs adjustments.
// This let developers make adjustments to the size and position of their overlays for each stream.
// This callback function is called prior to rendering every time when an overlay
// is rendered on a stream, which is useful if the resolution has been
// updated or rotation has changed.
func adjustmentCallback(adjustmentEvent *axoverlay.OverlayAdjustmentEvent) {
	app := adjustmentEvent.Userdata.(*acapapp.AcapApplication)
	app.Syslog.Infof("Adjust callback for overlay-%d: %dx%d", adjustmentEvent.OverlayId, adjustmentEvent.OverlayWidth, adjustmentEvent.OverlayHeight)
	app.Syslog.Infof("Adjust callback for stream: %dx%d", adjustmentEvent.Stream.Width, adjustmentEvent.Stream.Height)

	*adjustmentEvent.OverlayWidth = adjustmentEvent.Stream.Width
	*adjustmentEvent.OverlayHeight = adjustmentEvent.Stream.Height
}

// renderCallback is called whenever the system redraws an overlay
// This can happen in two cases, Redraw() is called or a new stream is started.
func renderCallback(renderEvent *axoverlay.OverlayRenderEvent) {
	app := renderEvent.Userdata.(*acapapp.AcapApplication)
	app.Syslog.Infof("Render callback for camera: %d", renderEvent.Stream.Camera)
	app.Syslog.Infof("Render callback for overlay-%d: %dx%d", renderEvent.OverlayId, renderEvent.OverlayWidth, renderEvent.OverlayHeight)
	app.Syslog.Infof("Render callback for stream: %dx%d", renderEvent.Stream.Width, renderEvent.Stream.Height)
	renderNextImageInSequenze(app, renderEvent)
}

// renderNextImageInSequenze renders the next image in the sequence
func renderNextImageInSequenze(app *acapapp.AcapApplication, renderEvent *axoverlay.OverlayRenderEvent) {
	surface, err := axoverlay.NewCairoSurfaceFromPNG(image_seq.NextImageFilename())
	if err != nil {
		app.Syslog.Errorf("Failed to load PNG: %v", err)
		return
	}
	defer surface.Destroy()
	centerX := renderEvent.Stream.Width/2 - surface.Width()/2
	stickyBottomY := renderEvent.Stream.Height - surface.Height()
	renderEvent.CairoCtx.PaintSurface(surface, float64(centerX), float64(stickyBottomY))
}

func main() {

	// Initialize a new ACAP application instance.
	// AcapApplication initializes the ACAP application with there name, eventloop, and syslog etc..
	app = acapapp.NewAcapApplication()

	// Initialize image sequence
	image_seq = NewImageSequence(24)

	// Overlayprovider is an highlevel wrapper around AxOvleray to make life easier
	if overlayProvider, err = acapapp.NewOverlayProvider(renderCallback, adjustmentCallback, streamSelectCallback); err != nil {
		panic(err)
	}
	app.AddCloseCleanFunc(overlayProvider.Cleanup)

	// we pass app as userdata to access syslog from app in callbacks
	if _, err = overlayProvider.AddOverlay(acapapp.NewAnchorCenterRrgbaOverlay(axoverlay.AxOverlayTopLeft, app)); err != nil {
		app.Syslog.Crit(err.Error())
	}

	// Draw overlays
	if err = overlayProvider.Redraw(); err != nil {
		app.Syslog.Crit(err.Error())
	}

	// Overlay update - increasing counter and call redraw to invoke a new render call
	ticker := time.NewTicker(time.Millisecond * 300)
	go func() {
		for true {
			<-ticker.C
			if err = overlayProvider.Redraw(); err != nil {
				app.Syslog.Crit(err.Error())
			}
		}
	}()

	// Run gmain loop with signal handler attached.
	// This will block the main thread until the application is stopped.
	// The application can be stopped by sending a signal to the process (e.g. SIGINT).
	// Axoverlay needs a running event loop to handle the overlay callbacks corretly
	app.Run()
}
