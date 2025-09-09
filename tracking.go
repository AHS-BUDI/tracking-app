package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	
	"github.com/go-vgo/robotgo"
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
		
		// Cek juga mouse click
		if robotgo.GetMouseDown("left") {
			_, err = stmt.Exec(x, y, "left_click", time.Now())
			if err != nil {
				log.Printf("Error recording mouse click: %v", err)
			}
		}
		
		if robotgo.GetMouseDown("right") {
			_, err = stmt.Exec(x, y, "right_click", time.Now())
			if err != nil {
				log.Printf("Error recording mouse click: %v", err)
			}
		}
		
		time.Sleep(100 * time.Millisecond)
	}
}

// Tracking keyboard
func trackKeyboard(db *sql.DB) {
	stmt, err := db.Prepare("INSERT INTO keyboard_events(key_pressed, event_type, timestamp) VALUES($1, $2, $3)")
	if err != nil {
		log.Fatalf("Error preparing statement: %v", err)
	}
	defer stmt.Close()
	
	// Hook untuk menangkap keystroke
	robotgo.EventHook(robotgo.KeyDown, []string{}, func(e robotgo.Event) {
		keyName := e.Key
		_, err = stmt.Exec(keyName, "keydown", time.Now())
		if err != nil {
			log.Printf("Error recording keyboard event: %v", err)
		}
	})
	
	// Perlu melakukan robotgo.EventEnd() jika aplikasi selesai
	s := robotgo.EventStart()
	<-robotgo.EventProcess(s)
}