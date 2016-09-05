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
	App = NewApp(":8080", "main")

	flag.Parse()

	// Run Migrations
	if *migrate {
		log.Println("Running Migrations")
		Migrate()
	}

	// Run The Application
	App.Run()
}

// Runs the auto migrations for specific objects
func Migrate() {
	App.db.AutoMigrate(&Tag{})
	App.db.AutoMigrate(&Notification{})
}
