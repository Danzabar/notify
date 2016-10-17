package main

import (
	"flag"
	"github.com/jasonlvhit/gocron"
	"gopkg.in/go-playground/validator.v9"
	"log"
)

var (
	App       *Application
	Validator *validator.Validate
)

// Main execution
func main() {
	migrate := flag.Bool("m", false, "Runs migrations before running server")
	dbDriver := flag.String("driver", "sqlite3", "The database driver notify should use")
	dbCreds := flag.String("creds", "/tmp/main.db", "The database credentials")

	flag.Parse()

	App = NewApp(":8080", *dbDriver, *dbCreds)

	// Run Migrations
	if *migrate {
		log.Println("Running Migrations")
		Migrate()
	}

	// Start the alert task schedule
	gocron.Every(1).Second().Do(SendAlerts)
	gocron.Start()

	// Run The Application
	App.Run()
}

// Runs the auto migrations for specific objects
func Migrate() {
	App.db.AutoMigrate(&Tag{})
	App.db.AutoMigrate(&Notification{})
	App.db.AutoMigrate(&AlertGroup{})
	App.db.AutoMigrate(&Recipient{})
}
