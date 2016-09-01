package main

import (
	"flag"
	"log"
)

// Global for Application
var App *Application

// Main execution
func main() {
	migrate := flag.Bool("m", false, "Runs migrations before running server")
	App = NewApp(":8080")

	flag.Parse()

	// Run Migrations
	if *migrate {
		log.Println("Running Migrations")
		App.db.AutoMigrate(&Tag{})
		App.db.AutoMigrate(&Notification{})
	}

	// Run The Application
	App.Run()
}
