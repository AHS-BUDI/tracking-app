package main

import (
	"database/sql"
	"time"
	"log"
	
	"github.com/gofiber/fiber/v2"
)

type MouseEvent struct {
	ID        int       `json:"id"`
	XPosition int       `json:"x_position"`
	YPosition int       `json:"y_position"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
}

type KeyboardEvent struct {
	ID         int       `json:"id"`
	KeyPressed string    `json:"key_pressed"`
	EventType  string    `json:"event_type"`
	Timestamp  time.Time `json:"timestamp"`
}

func setupAPI(db *sql.DB) {
	app := fiber.New()
	
	// Endpoint untuk mendapatkan data mouse
	app.Get("/api/mouse", func(c *fiber.Ctx) error {
		var events []MouseEvent
		
		rows, err := db.Query("SELECT id, x_position, y_position, event_type, timestamp FROM mouse_events ORDER BY timestamp DESC LIMIT 100")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()
		
		for rows.Next() {
			var event MouseEvent
			if err := rows.Scan(&event.ID, &event.XPosition, &event.YPosition, &event.EventType, &event.Timestamp); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
			events = append(events, event)
		}
		
		return c.JSON(events)
	})
	
	// Endpoint untuk mendapatkan data keyboard
	app.Get("/api/keyboard", func(c *fiber.Ctx) error {
		var events []KeyboardEvent
		
		rows, err := db.Query("SELECT id, key_pressed, event_type, timestamp FROM keyboard_events ORDER BY timestamp DESC LIMIT 100")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()
		
		for rows.Next() {
			var event KeyboardEvent
			if err := rows.Scan(&event.ID, &event.KeyPressed, &event.EventType, &event.Timestamp); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
			events = append(events, event)
		}
		
		return c.JSON(events)
	})
	
	log.Fatal(app.Listen(":3000"))
}