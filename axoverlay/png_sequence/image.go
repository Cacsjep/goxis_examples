package main

import "fmt"

// ImageSequence is a simple struct to keep track of the current image index
type ImageSequence struct {
	currentIndex int
	maxIndex     int
}

// NewImageSequence creates a new ImageSequence with the given maxIndex
func NewImageSequence(maxIndex int) *ImageSequence {
	return &ImageSequence{
		currentIndex: 0,
		maxIndex:     maxIndex,
	}
}

// NextImageFilename returns the next image filename in the sequence
func (seq *ImageSequence) NextImageFilename() string {
	filename := fmt.Sprintf("zinta/%d.png", seq.currentIndex)
	seq.currentIndex++
	if seq.currentIndex > seq.maxIndex {
		seq.currentIndex = 0
	}
	return filename
}
