package main

import (
	"flag"
)

// Global for Application
var App *Application

// Main execution
func main() {
	migrate := flag.Bool("m", false, "Runs migrations before running server")
	App := NewApp(":8080")

	// Run Migrations
	if *migrate {
		App.db.AutoMigrate(&Tag{})
		App.db.AutoMigrate(&Notification{})
	}

	// Run The Application
	App.Run()
}
