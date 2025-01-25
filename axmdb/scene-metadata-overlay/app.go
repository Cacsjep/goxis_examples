package main

// This example demonstrates how to use the axis message broker api and axoverlay to display the scene metadata
// main is the entry point of the application. It initializes a new instance
// of mdbSceneOverlayApp and runs it. If there is an error during initialization,
// the application will panic and terminate.
// The logic of the app is in the msoa.go file.
func main() {
	msoa, err := newMdbSceneOverlayApp()
	if err != nil {
		panic(err)
	}
	msoa.Run()
}
