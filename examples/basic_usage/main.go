package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/mikeschinkel/go-logutil"
)

// Example demonstrating basic usage of go-logutil for logging utilities
func main() {
	fmt.Println("go-logutil Basic Usage Example")
	fmt.Println("===============================\n")

	// Example 1: LogArgs with simple struct
	fmt.Println("1. LogArgs with Simple Struct")
	fmt.Println("   Converts struct fields to slog attributes")

	type Person struct {
		Name string
		Age  int
	}

	person := Person{
		Name: "Alice",
		Age:  30,
	}

	attrs := logutil.LogArgs(person)
	fmt.Printf("   Struct: %+v\n", person)
	fmt.Printf("   Log attributes count: %d\n", len(attrs))

	// Example 2: LogArgs with JSON tags
	fmt.Println("\n2. LogArgs with JSON Tags")
	fmt.Println("   Respects json tag names and omitempty")

	type Config struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Password string `json:"password,omitempty"`
		Timeout  int    `json:"timeout,omitempty"`
	}

	config := Config{
		Host: "localhost",
		Port: 8080,
		// Password and Timeout are zero values and have omitempty
	}

	configAttrs := logutil.LogArgs(config)
	fmt.Printf("   Config: %+v\n", config)
	fmt.Printf("   Log attributes count: %d (omitted zero values)\n", len(configAttrs))

	// Example 3: LogArgs with time.Time
	fmt.Println("\n3. LogArgs with time.Time")
	fmt.Println("   Formats time.Time as RFC3339")

	type Event struct {
		Name      string
		Timestamp time.Time
	}

	event := Event{
		Name:      "UserLogin",
		Timestamp: time.Now(),
	}

	eventAttrs := logutil.LogArgs(event)
	fmt.Printf("   Event: %s at %s\n", event.Name, event.Timestamp.Format(time.RFC3339))
	fmt.Printf("   Log attributes count: %d\n", len(eventAttrs))

	// Example 4: LogArgs with nested struct
	fmt.Println("\n4. LogArgs with Nested Struct")
	fmt.Println("   Creates nested groups for struct fields")

	type Address struct {
		Street string
		City   string
	}

	type User struct {
		Name    string
		Address Address
	}

	user := User{
		Name: "Bob",
		Address: Address{
			Street: "123 Main St",
			City:   "Springfield",
		},
	}

	userAttrs := logutil.LogArgs(user)
	fmt.Printf("   User: %+v\n", user)
	fmt.Printf("   Log attributes count: %d (includes nested group)\n", len(userAttrs))

	// Example 5: Using LogArgs with slog logger
	fmt.Println("\n5. Using LogArgs with slog Logger")
	fmt.Println("   Demonstrates integration with standard slog")

	logger := slog.Default()

	type LogEvent struct {
		Action string
		UserID int
		Status string
	}

	logEvent := LogEvent{
		Action: "file_upload",
		UserID: 42,
		Status: "success",
	}

	logger.Info("Event occurred", logutil.LogArgs(logEvent)...)
	fmt.Println("   (Logged event with structured attributes)")

	// Example 6: LogArgs with pointer
	fmt.Println("\n6. LogArgs with Pointer")
	fmt.Println("   Handles pointer to struct")

	personPtr := &Person{
		Name: "Charlie",
		Age:  25,
	}

	ptrAttrs := logutil.LogArgs(personPtr)
	fmt.Printf("   Person pointer: %+v\n", personPtr)
	fmt.Printf("   Log attributes count: %d\n", len(ptrAttrs))

	// Example 7: LogArgs with nil
	fmt.Println("\n7. LogArgs with nil")
	fmt.Println("   Returns empty slice for nil input")

	nilAttrs := logutil.LogArgs(nil)
	fmt.Printf("   Nil input result: %d attributes\n", len(nilAttrs))

	fmt.Println("\nUsage Notes:")
	fmt.Println("- LogArgs extracts struct fields as slog attributes")
	fmt.Println("- Respects `json` tags for field names and omitempty")
	fmt.Println("- Handles time.Time, errors, and nested structs specially")
	fmt.Println("- Skips unexported fields automatically")
}
