package main

import (
	"fmt"

	"github.com/Cacsjep/goxis/pkg/axlarod"
)

var (
	FLOAT_SIZE    = uint(4)
	TENSOR_1_SIZE = 1 * 50 * 4 * FLOAT_SIZE // Locations: [1, 50, 4]
	TENSOR_2_SIZE = 1 * 50 * FLOAT_SIZE     // Classes: [1, 50]
	TENSOR_3_SIZE = 1 * 50 * FLOAT_SIZE     // Scores: [1, 50]
	TENSOR_4_SIZE = 1 * FLOAT_SIZE          // Detections: [1]
	BBOX_SIZE     = 4                       // Each box has 4 coordinates
)

// PredictionResult holds the probabilities of detecting specific objects (e.g., persons, cars)
type PredictionResult struct {
	Detections []Detection
}

// InitializeDetectionModel configures a detection model with the given model file and hardware chip.
// It sets up memory-mapped file configurations for input and output tensors.
// Returns an error if model initialization fails.
func (lea *larodExampleApplication) InitalizeDetectionModel(modelFilePath string, chipString string) error {
	model_defs := axlarod.MemMapConfiguration{
		InputTmpMapFiles: map[int]*axlarod.MemMapFile{
			0: lea.PPModel.Outputs[0].MemMapFile, // Using of ppmodel output as input for detection model
		},
		OutputTmpMapFiles: map[int]*axlarod.MemMapFile{
			0: {Size: TENSOR_1_SIZE}, // BoundingBox locations
			1: {Size: TENSOR_2_SIZE}, // Classes
			2: {Size: TENSOR_3_SIZE}, // Scores
			3: {Size: TENSOR_4_SIZE}, // Number of detections
		},
	}

	if lea.DetectionModel, err = lea.app.Larod.NewInferModel(modelFilePath, chipString, model_defs, nil); err != nil {
		return fmt.Errorf("failed to create detection model: %w", err)
	}
	lea.app.AddModelCleaner(lea.DetectionModel)
	return nil
}

// getDResult retrieves the detection results from the model's output tensors.
// Returns the raw byte output and an error if retrieval fails.
func (lea *larodExampleApplication) getDResult() (*mobileNetFaceResult, error) {
	t1, err := lea.DetectionModel.Outputs[0].GetDataAsFloat32Slice(int(TENSOR_1_SIZE))
	if err != nil {
		return nil, err
	}
	t2, err := lea.DetectionModel.Outputs[1].GetDataAsFloat32Slice(int(TENSOR_2_SIZE))
	if err != nil {
		return nil, err
	}
	t3, err := lea.DetectionModel.Outputs[2].GetDataAsFloat32Slice(int(TENSOR_3_SIZE))
	if err != nil {
		return nil, err
	}
	t4, err := lea.DetectionModel.Outputs[3].GetDataAsFloat32()
	if err != nil {
		return nil, err
	}
	return &mobileNetFaceResult{locations: t1, classes: t2, scores: t3, detections: int(t4)}, nil
}

// Inference executes the model and retrieves the processed results.
// It ensures the model's file pointers are correctly positioned before execution.
// Returns a JobResult containing the inference results or an error if the inference process fails.
func (lea *larodExampleApplication) Inference() (*axlarod.JobResult, error) {

	// Rewind all output files position before each job.
	if err = lea.DetectionModel.RewindAllOutputsMemMapFiles(); err != nil {
		return nil, err
	}

	var result *axlarod.JobResult
	if result, err = lea.app.Larod.ExecuteJob(lea.DetectionModel, func() error {
		return nil // is feeded via memmap
	}, func() (any, error) {
		return lea.getDResult()
	}); err != nil {
		return nil, err
	}

	return result, nil
}

type mobileNetFaceResult struct {
	locations  []float32
	classes    []float32
	scores     []float32
	detections int
}

type Detection struct {
	Score, Class float32
	Box          axlarod.BoundingBox
}

// InferenceOutputRead converts raw model output data into structured prediction results.
// Returns a PredictionResult or an error if data conversion fails.
//
// https://github.com/AxisCommunications/acap-native-sdk-examples/blob/7bff215e7673e4a72630bb89f04c2f7b64cf319c/object-detection/app/object_detection.c#L942
func (lea *larodExampleApplication) InferenceOutputRead(result *mobileNetFaceResult) (*PredictionResult, error) {
	var detections []Detection
	for i := 0; i < result.detections; i++ {
		score := result.scores[i]
		if score >= lea.threshold {
			if BBOX_SIZE*i+3 < len(result.locations) { // Ensure within bounds
				detections = append(detections, Detection{
					Box: axlarod.BoundingBox{
						Top:    result.locations[BBOX_SIZE*i],
						Left:   result.locations[BBOX_SIZE*i+1],
						Bottom: result.locations[BBOX_SIZE*i+2],
						Right:  result.locations[BBOX_SIZE*i+3],
					},
					Score: score,
					Class: result.classes[i],
				})
			}
		}
	}
	lea.detections = detections
	return &PredictionResult{Detections: detections}, nil
}
