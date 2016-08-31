package main

import ()

// Global for Application
var App *Application

// Main execution
func main() {
	App := NewApp(":8080")

	// Run The Application
	App.Run()
}
