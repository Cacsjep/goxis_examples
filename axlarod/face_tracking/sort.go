package main

import (
	"time"

	"github.com/Cacsjep/goxis/pkg/axlarod"
)

// Track represents a single tracked object.
// Fields:
//   - ID: Unique identifier for the track.
//   - Box: Current bounding box of the tracked object.
//   - Age: Number of frames the track has been active.
//   - Missed: Number of consecutive frames without a detection.
//   - SortScore: Average IOU score across all detections assigned to this track.
type Track struct {
	ID          int                 // Unique identifier
	Box         axlarod.BoundingBox // Bounding box
	Age         int                 // Active frames
	Missed      int                 // Missed frames
	SortScore   float32             // Average IOU
	CreatedTime time.Time           // Time the track was created
}

// ActiveTime returns the duration since the track was created.
func (t *Track) ActiveTime() time.Duration {
	return time.Since(t.CreatedTime)
}

// SORT represents a simple online and realtime tracking system.
// Fields:
//   - Tracks: A map of active tracks indexed by their unique ID.
//   - NextTrackID: The next unique ID to assign to a new track.
//   - MaxMissed: Maximum frames to keep a track without detections before removal.
//   - MinScore: Minimum detection confidence score required for tracking.
//   - IOUThreshold: Minimum IOU required to match a detection with a track.
type SORT struct {
	Tracks       map[int]*Track // Active tracks
	NextTrackID  int            // Next track ID
	MaxMissed    int            // Max allowed missed frames
	MinScore     float32        // Minimum detection score
	IOUThreshold float32        // IOU threshold for matching
}

// NewSORT initializes a new SORT tracker.
// Args:
//   - maxMissed: Maximum frames a track can be unmatched before being removed.
//   - minScore: Minimum detection score required to consider a detection.
//   - iouThreshold: IOU threshold to determine a match between a detection and a track.
//
// Returns:
//
//	*SORT: A pointer to the initialized SORT instance.
func NewSORT(maxMissed int, minScore, iouThreshold float32) *SORT {
	return &SORT{
		Tracks:       make(map[int]*Track),
		NextTrackID:  1,
		MaxMissed:    maxMissed,
		MinScore:     minScore,
		IOUThreshold: iouThreshold,
	}
}

// IOU computes the Intersection over Union (IOU) between two bounding boxes.
// Args:
//   - box1: The first bounding box.
//   - box2: The second bounding box.
//
// Returns:
//
//	float32: The IOU score, ranging from 0 (no overlap) to 1 (perfect overlap).
//
// Behavior:
// - Returns 0 if the boxes do not overlap.
// - Handles edge cases where areas are zero or overlap is invalid.
func IOU(box1, box2 axlarod.BoundingBox) float32 {
	// Calculate intersection coordinates
	interLeft := max(box1.Left, box2.Left)
	interTop := max(box1.Top, box2.Top)
	interRight := min(box1.Right, box2.Right)
	interBottom := min(box1.Bottom, box2.Bottom)

	// If no overlap, return 0 immediately
	if interLeft >= interRight || interTop >= interBottom {
		return 0
	}

	// Calculate intersection area
	intersection := (interRight - interLeft) * (interBottom - interTop)

	// Calculate areas of both boxes
	area1 := (box1.Right - box1.Left) * (box1.Bottom - box1.Top)
	area2 := (box2.Right - box2.Left) * (box2.Bottom - box2.Top)

	// Union area
	union := area1 + area2 - intersection

	// Handle edge case where union area is zero (shouldn't occur if inputs are valid)
	if union <= 0 {
		return 0
	}

	// Return IOU as intersection over union
	return intersection / union
}

// Update processes the current frame's detections and updates the tracker's state.
//
// This function performs the following steps:
//  1. Matches the current detections to existing tracks using the Intersection over Union (IOU) metric.
//     Tracks are updated with the matched detection, their bounding boxes are adjusted, and their ages are incremented.
//  2. Creates new tracks for unmatched detections. These represent objects that have appeared for the first time.
//  3. Increments the missed count for unmatched tracks and removes stale tracks that exceed the `MaxMissed` threshold.
//
// Args:
//
//	detections []Detection: A slice of Detection objects representing the detected objects in the current frame.
//	  - Each detection has a confidence score, bounding box, and class information.
//
// Returns:
//
//	[]Detection: A slice of Detection objects updated with their assigned track IDs.
//	  - Detections that result in new tracks will have `IsNew` set to true.
//	  - Each detection includes the IOU score (`IOU`) and the matched track ID (`MatchedTrackID`).
//
// Behavior:
//   - Tracks are only matched to detections with a confidence score above the `MinScore` threshold.
//   - Matching is performed by finding the track with the highest IOU score above the `IOUThreshold`.
//   - Tracks that do not receive a detection are considered unmatched. Their `Missed` count is incremented, and they are
//     removed if the count exceeds `MaxMissed`.
func (s *SORT) Update(detections []Detection) []Detection {
	updatedDetections := make([]Detection, 0, len(detections)) // Preallocate for speed

	// Track assignments
	assignedTracks := make(map[int]struct{}, len(s.Tracks))
	assignedDetections := make(map[int]struct{}, len(detections))

	// Step 1: Match detections to existing tracks
	for detIdx, detection := range detections {
		if detection.Score < s.MinScore {
			continue // Skip low-score detections
		}

		var (
			bestTrackID int
			bestIOU     float32
		)
		foundMatch := false

		for trackID, track := range s.Tracks {
			// Skip already assigned tracks
			if _, assigned := assignedTracks[trackID]; assigned {
				continue
			}

			// Calculate IOU
			iou := IOU(detection.Box, track.Box)
			if iou >= s.IOUThreshold && iou > bestIOU {
				bestIOU = iou
				bestTrackID = trackID
				foundMatch = true
			}
		}

		if foundMatch {
			// Update track with matched detection
			track := s.Tracks[bestTrackID]
			track.Box = detection.Box
			track.Age++
			track.Missed = 0

			// Update SortScore as a running average of IOUs
			track.SortScore = (track.SortScore*float32(track.Age-1) + bestIOU) / float32(track.Age)

			// Mark track and detection as assigned
			assignedTracks[bestTrackID] = struct{}{}
			assignedDetections[detIdx] = struct{}{}

			// Assign track ID to detection
			detection.ID = track.ID
			detection.IsNew = false
			detection.Age = track.Age
			detection.TrackingSince = track.ActiveTime()

			updatedDetections = append(updatedDetections, detection)
		}
	}

	// Step 2: Create new tracks for unmatched detections
	for detIdx, detection := range detections {
		if _, assigned := assignedDetections[detIdx]; assigned {
			continue // Skip already matched detections
		}

		// Create a new track
		trackID := s.NextTrackID
		s.NextTrackID++

		s.Tracks[trackID] = &Track{
			ID:          trackID,
			Box:         detection.Box,
			Age:         1,
			Missed:      0,
			SortScore:   0,
			CreatedTime: time.Now(),
		}

		// Assign new track ID to detection
		detection.ID = trackID
		detection.IsNew = true
		detection.TrackingSince = s.Tracks[trackID].ActiveTime()
		updatedDetections = append(updatedDetections, detection)
	}

	// Step 3: Clean up unmatched tracks
	toDelete := make([]int, 0, len(s.Tracks)) // Preallocate potential deletions
	for trackID, track := range s.Tracks {
		if _, assigned := assignedTracks[trackID]; assigned {
			continue // Skip matched tracks
		}

		// Increment missed count
		track.Missed++
		if track.Missed > s.MaxMissed {
			toDelete = append(toDelete, trackID) // Mark track for deletion
		}
	}

	// Delete stale tracks
	for _, trackID := range toDelete {
		delete(s.Tracks, trackID)
	}

	return updatedDetections
}

// GetAverageSortScore calculates and returns the average SortScore for all active tracks.
//
// The SortScore is a metric that evaluates the quality of tracking for each individual track,
// calculated as the running average of Intersection over Union (IOU) scores between a track
// and its assigned detections. Tracks with higher SortScores indicate better tracking performance.
//
// This function only includes tracks that have been active for more than one frame (`Age > 1`)
// to ensure meaningful evaluation, as single-frame tracks may not have sufficient data for an accurate score.
//
// Returns:
//
//	float32 - The average SortScore of all active tracks. If no tracks meet the criteria, the function
//	          returns 0 to indicate that no meaningful score could be calculated.
func (s *SORT) GetAverageSortScore() float32 {
	var totalScore float32
	var trackCount int

	for _, track := range s.Tracks {
		if track.Age > 1 { // Only consider tracks that have been active for more than one frame
			totalScore += track.SortScore
			trackCount++
		}
	}

	if trackCount == 0 {
		return 0
	}
	return totalScore / float32(trackCount)
}

// Inline helper functions for maximum and minimum
func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}
