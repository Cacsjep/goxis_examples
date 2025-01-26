package main

import (
	"fmt"
	"image/color"

	"github.com/Cacsjep/goxis/pkg/axoverlay"
)

// Initialize the overlay provider
func (lea *larodExampleApplication) InitOverlay() error {
	if lea.overlayProvider, err = axoverlay.NewOverlayProvider(renderCallback, nil, nil); err != nil {
		return err
	}
	lea.app.AddCloseCleanFunc(lea.overlayProvider.Cleanup)
	if _, err = lea.overlayProvider.AddOverlay(axoverlay.NewAnchorCenterRrgbaOverlay(axoverlay.AxOverlayTopLeft, lea)); err != nil {
		return err
	}
	return nil
}

// renderCallback is used to draw bounding boxes from the detections via axoverlay
func renderCallback(renderEvent *axoverlay.OverlayRenderEvent) {
	lea := renderEvent.Userdata.(*larodExampleApplication)
	renderEvent.CairoCtx.DrawTransparent(renderEvent.Stream.Width, renderEvent.Stream.Height)

	// Draw the sort tracker average score
	renderEvent.CairoCtx.DrawText(fmt.Sprintf("Tracking score: %d%%", int(lea.sortTracker.GetAverageSortScore()*100)), 10, 10, 32.0, "serif", axoverlay.ColorBlack)

	for _, obj := range lea.prediction_result.Detections {
		scaled_box := obj.Box.Scale(renderEvent.Stream.Width, renderEvent.Stream.Height)
		cords := scaled_box.ToCords64()
		DrawBoundingBox(
			renderEvent.CairoCtx,
			cords.X,
			cords.Y,
			cords.W,
			cords.H,
			axoverlay.ColorBlack,
			fmt.Sprintf("ID-%d %d%%, %d sec", obj.ID, int(obj.Score*100), int(obj.TrackingSince.Seconds())),
			axoverlay.ColorWite,
			13,
			"sans",
			100,
		)
	}
}

func DrawBoundingBox(ctx *axoverlay.CairoContext, x float64, y float64, width float64, height float64, rectColor color.RGBA, label string, labelColor color.RGBA, labelSize float64, labelFont string, minBoxSizeRenderW int) {
	rectLinewidth := float64(3)
	ctx.DrawBoundingBoxRect(x, y, width, height, rectColor, rectLinewidth, 0.3)
	ctx.DrawBoundingBoxLabel(label, x-(rectLinewidth/2), y-(rectLinewidth/2), 7, labelSize, labelFont, labelColor, rectColor)
}
