package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

// Membuat tabel di PostgreSQL
func createTables(db *sql.DB) {
	mouseTableSQL := `
	CREATE TABLE IF NOT EXISTS mouse_events (
		id SERIAL PRIMARY KEY,
		x_position INTEGER NOT NULL,
		y_position INTEGER NOT NULL,
		event_type VARCHAR(20) NOT NULL,
		timestamp TIMESTAMP NOT NULL
	);`
	
	keyboardTableSQL := `
	CREATE TABLE IF NOT EXISTS keyboard_events (
		id SERIAL PRIMARY KEY,
		key_pressed VARCHAR(50) NOT NULL,
		event_type VARCHAR(20) NOT NULL,
		timestamp TIMESTAMP NOT NULL
	);`
	
	_, err := db.Exec(mouseTableSQL)
	if err != nil {
		log.Fatalf("Error creating mouse_events table: %v", err)
	}
	
	_, err = db.Exec(keyboardTableSQL)
	if err != nil {
		log.Fatalf("Error creating keyboard_events table: %v", err)
	}
	
	fmt.Println("Database tables created successfully")
}

// Tracking pergerakan mouse
func trackMouse(db *sql.DB) {
	stmt, err := db.Prepare("INSERT INTO mouse_events(x_position, y_position, event_type, timestamp) VALUES($1, $2, $3, $4)")
	if err != nil {
		log.Fatalf("Error preparing statement: %v", err)
	}
	defer stmt.Close()
	
	lastX, lastY := robotgo.GetMousePos()
	
	// Setup mouse event listener
	hook.Register(hook.MouseDown, []string{}, func(e hook.Event) {
		x, y := int(e.X), int(e.Y)
		
		var eventType string
		if e.Button == 1 {
			eventType = "left_click"
		} else if e.Button == 2 {
			eventType = "right_click"
		} else {
			eventType = fmt.Sprintf("button_%d_click", e.Button)
		}
		
		_, err = stmt.Exec(x, y, eventType, time.Now())
		if err != nil {
			log.Printf("Error recording mouse click: %v", err)
		}
	})
	
	// Separate goroutine for tracking mouse movement
	go func() {
		for {
			x, y := robotgo.GetMousePos()
			
			// Hanya simpan jika posisi berubah
			if x != lastX || y != lastY {
				_, err = stmt.Exec(x, y, "move", time.Now())
				if err != nil {
					log.Printf("Error recording mouse movement: %v", err)
				}
				
				lastX, lastY = x, y
			}
			
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

// Tracking keyboard menggunakan gohook
func trackKeyboard(db *sql.DB) {
	stmt, err := db.Prepare("INSERT INTO keyboard_events(key_pressed, event_type, timestamp) VALUES($1, $2, $3)")
	if err != nil {
		log.Fatalf("Error preparing statement: %v", err)
	}
	defer stmt.Close()
	
	// Register keyboard event handler
	hook.Register(hook.KeyDown, []string{}, func(e hook.Event) {
		// Get key name from the event
		keyName := hook.RawcodetoKeychar(e.Rawcode)
		if keyName == "" {
			keyName = fmt.Sprintf("Key: %v", e.Rawcode)
		}
		
		_, err = stmt.Exec(keyName, "keydown", time.Now())
		if err != nil {
			log.Printf("Error recording keyboard event: %v", err)
		}
	})
}

// Start tracking function to initialize and start all event hooks
func startTracking(db *sql.DB) {
	// Initialize tables
	createTables(db)
	
	// Set up tracking
	trackMouse(db)
	trackKeyboard(db)
	
	// Start the event hook
	log.Println("Starting event tracking...")
	s := hook.Start()
	<-hook.Process(s)
}