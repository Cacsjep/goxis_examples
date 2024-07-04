package main

import (
	"math"
	"time"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axoverlay"
)

type PixelChoord struct {
	X float64
	Y float64
}

const numPoints = 600

var PIXEL_ARRAY []PixelChoord

// Fill the pixel array with a circle of points
func init() {
	cx, cy := 0.5, 0.5
	radius := 0.4

	for i := 0; i < numPoints*2; i++ {
		angle := float64(i) * math.Pi / float64(numPoints)
		r := radius
		if i%2 == 1 {
			r /= 2.0
		}
		x := cx + r*math.Cos(angle)
		y := cy - r*math.Sin(angle)
		PIXEL_ARRAY = append(PIXEL_ARRAY, PixelChoord{X: x, Y: y})
	}
}

// This example demonstrate how to use overlay provider to an array of pixels on a stream.
//
// ! Note: Overlay callbacks only invoked when stream is viewed via web ui or rtsp etc..

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
	renderEvent.CairoCtx.DrawTransparent(renderEvent.Stream.Width, renderEvent.Stream.Height)
	for _, p := range PIXEL_ARRAY {
		renderEvent.CairoCtx.SetOperator(axoverlay.OPERATOR_SOURCE)
		renderEvent.CairoCtx.Rectangle(p.X*float64(renderEvent.Stream.Width), p.Y*float64(renderEvent.Stream.Height), 3, 3) // for visibility we draw a 3x3 pixel
		renderEvent.CairoCtx.SetSourceRGBA(axoverlay.ColorMaterialRed)
		renderEvent.CairoCtx.Fill()
	}
}

func main() {
	// Initialize a new ACAP application instance.
	// AcapApplication initializes the ACAP application with there name, eventloop, and syslog etc..
	app := acapapp.NewAcapApplication()

	// Overlayprovider is an highlevel wrapper around AxOvleray to make life easier
	overlayProvider, err := axoverlay.NewOverlayProvider(renderCallback, adjustmentCallback, streamSelectCallback)
	if err != nil {
		app.Syslog.Crit(err.Error())
	}
	app.AddCloseCleanFunc(overlayProvider.Cleanup)

	// we pass app as userdata to access syslog from app in callbacks
	if _, err = overlayProvider.AddOverlay(axoverlay.NewAnchorCenterRrgbaOverlay(axoverlay.AxOverlayTopLeft, app)); err != nil {
		app.Syslog.Crit(err.Error())
	}

	// Draw overlays
	if err = overlayProvider.Redraw(); err != nil {
		app.Syslog.Crit(err.Error())
	}

	// Overlay update
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
