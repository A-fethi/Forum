package utils

import (
	"fmt"
	"forum/app/config"
)

func ValidatePost(title, content string) error {
	if title == "" {
		config.Logger.Println("Validation failed: Title is empty")
		return fmt.Errorf("Title cannot be empty")
	}
	if content == "" {
		config.Logger.Println("Validation failed: Content is empty")
		return fmt.Errorf("Content cannot be empty")
	}

	if len(title) > 100 {
		config.Logger.Printf("Validation failed: Title length is invalid (title length: %d)", len(title))
		return fmt.Errorf("Title must be between 5 and 100 characters")
	}

	if len(content) > 5000 {
		config.Logger.Printf("Validation failed: Content length is invalid (content length: %d)", len(content))
		return fmt.Errorf("Content must be between 10 and 5000 characters")
	}

	config.Logger.Println("Validation successful for title and content")
	return nil
}
