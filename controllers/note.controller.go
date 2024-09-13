package controllers

import (
	"strconv"
	"strings"
	"time"

	initializers "golang-fiber/initializers"
	"golang-fiber/models"

	"github.com/gofiber/fiber/v2" // Fiber framework for building the API
	"gorm.io/gorm"                // ORM for interacting with the database
)

// CreateNoteHandler handles the creation of a new note.
// It parses the request body, validates the input, and creates a new note in the database.
func CreateNoteHandler(c *fiber.Ctx) error {
	// Define a variable to hold the parsed request body
	var payload *models.CreateNoteSchema

	// Parse the request body into the payload struct
	if err := c.BodyParser(&payload); err != nil {
		// If parsing fails, return a bad request error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Validate the payload
	errors := models.ValidateStruct(payload)
	if errors != nil {
		// If validation fails, return a bad request error with validation errors
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	// Create a new note with the provided data
	now := time.Now()
	newNote := models.Note{
		Title:     payload.Title,
		Content:   payload.Content,
		Category:  payload.Category,
		Published: payload.Published,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Attempt to save the new note to the database
	result := initializers.DB.Create(&newNote)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
			// Handle the case where the note title already exists
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "Title already exists, please use another title"})
		}
		// Handle other database errors
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error.Error()})
	}

	// Return a success response with the created note
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"note": newNote}})
}

// FindNotes retrieves a list of notes with pagination support.
// It fetches notes from the database based on the provided page and limit query parameters.
func FindNotes(c *fiber.Ctx) error {
	// Extract page and limit query parameters
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var notes []models.Note
	// Fetch notes from the database with pagination
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&notes)
	if results.Error != nil {
		// Handle database errors
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	// Return a success response with the list of notes
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(notes), "notes": notes})
}

// UpdateNote updates an existing note based on the provided noteId.
// It parses the request body, performs the update, and returns the updated note.
func UpdateNote(c *fiber.Ctx) error {
	// Extract the noteId from the URL parameters
	noteId := c.Params("noteId")

	var payload *models.UpdateNoteSchema

	// Parse the request body into the payload struct
	if err := c.BodyParser(&payload); err != nil {
		// If parsing fails, return a bad request error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var note models.Note
	// Find the existing note by its ID
	result := initializers.DB.First(&note, "id = ?", noteId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Handle the case where the note does not exist
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No note with that ID exists"})
		}
		// Handle other database errors
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Prepare the updates based on the provided payload
	updates := make(map[string]interface{})
	if payload.Title != "" {
		updates["title"] = payload.Title
	}
	if payload.Category != "" {
		updates["category"] = payload.Category
	}
	if payload.Content != "" {
		updates["content"] = payload.Content
	}
	if payload.Published != nil {
		updates["published"] = payload.Published
	}

	// Update the note with the new data
	updates["updated_at"] = time.Now()
	initializers.DB.Model(&note).Updates(updates)

	// Return a success response with the updated note
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"note": note}})
}

// FindNoteById retrieves a single note by its ID.
// It returns the note if found or an error if not.
func FindNoteById(c *fiber.Ctx) error {
	// Extract the noteId from the URL parameters
	noteId := c.Params("noteId")

	var note models.Note
	// Find the note by its ID
	result := initializers.DB.First(&note, "id = ?", noteId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Handle the case where the note does not exist
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No note with that ID exists"})
		}
		// Handle other database errors
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Return a success response with the note data
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"note": note}})
}

// DeleteNote deletes a note by its ID.
// It returns a success response if the note is deleted or an error if not found.
func DeleteNote(c *fiber.Ctx) error {
	// Extract the noteId from the URL parameters
	noteId := c.Params("noteId")

	// Attempt to delete the note by its ID
	result := initializers.DB.Delete(&models.Note{}, "id = ?", noteId)

	if result.RowsAffected == 0 {
		// Handle the case where the note does not exist
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No note with that ID exists"})
	} else if result.Error != nil {
		// Handle other database errors
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	// Return a no content response to indicate successful deletion
	return c.SendStatus(fiber.StatusNoContent)
}
