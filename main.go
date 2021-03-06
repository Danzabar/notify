package main

import (
	"flag"
	"github.com/jasonlvhit/gocron"
	"gopkg.in/go-playground/validator.v9"
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
	enableAlert := flag.Bool("a", false, "Enables the alerting schedule")
	runAlert := flag.Bool("r", false, "Runs alerting straight away")
	port := flag.String("port", ":8080", "Port on which the server runs")
	user := flag.String("user", "admin", "Username for basic auth stuffs")
	pass := flag.String("pass", "changeme", "Password for authentication")

	flag.Parse()

	App = NewApp(*port, *dbDriver, *dbCreds, *user, *pass)

	// Run Migrations
	if *migrate {
		App.log.Debug("Running Migrations")
		Migrate()
	}

	// Start the alert task schedule
	if *enableAlert {
		if *runAlert {
			SendAlerts()
		}
		App.log.Debug("Starting Alert scheduler")
		gocron.Every(1).Minute().Do(SendAlerts)
		gocron.Start()
	}

	// Run The Application
	App.Run()
}

// Runs the auto migrations for specific objects
func Migrate() {
	App.db.AutoMigrate(&Tag{})
	App.db.AutoMigrate(&Notification{})
	App.db.AutoMigrate(&AlertGroup{})
}
