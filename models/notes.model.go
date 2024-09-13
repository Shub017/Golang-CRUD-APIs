package models

import (
	"time"

	"github.com/go-playground/validator/v10" // Importing the validator package for struct validation
	"github.com/google/uuid"                 // Importing UUID package for generating unique identifiers
)

// Note represents a note model with various fields
type Note struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"` // Unique identifier for the note, auto-generated UUID
	Title     string    `gorm:"varchar(255);uniqueIndex;not null" json:"title,omitempty"`             // Title of the note, unique and required
	Content   string    `gorm:"not null" json:"content,omitempty"`                                    // Content of the note, required
	Category  string    `gorm:"varchar(100)" json:"category,omitempty"`                               // Category of the note, optional
	Published bool      `gorm:"default:false;not null" json:"published"`                              // Published status of the note, default is false
	CreatedAt time.Time `gorm:"not null" json:"createdAt,omitempty"`                                  // Creation timestamp of the note, required
	UpdatedAt time.Time `gorm:"not null" json:"updatedAt,omitempty"`                                  // Last update timestamp of the note, required
}

// Validator instance for validating structs
var validate = validator.New()

// ErrorResponse represents the structure of validation errors
type ErrorResponse struct {
	Field string `json:"field"`           // Name of the field that caused the validation error
	Tag   string `json:"tag"`             // Validation tag indicating the type of validation error
	Value string `json:"value,omitempty"` // Value that caused the validation error, optional
}

// ValidateStruct validates a given struct based on tags and returns a slice of ErrorResponse
func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse

	// Validate the struct using the validator instance
	err := validate.Struct(payload)
	if err != nil {
		// Process each validation error
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace() // Field name with error
			element.Tag = err.Tag()               // Validation tag (e.g., "required")
			element.Value = err.Param()           // Parameter value associated with the tag (e.g., min length)
			errors = append(errors, &element)     // Add the error to the slice
		}
	}

	return errors // Return the slice of validation errors
}

// CreateNoteSchema represents the schema for creating a new note
type CreateNoteSchema struct {
	Title     string `json:"title" validate:"required"`   // Title of the note, required
	Content   string `json:"content" validate:"required"` // Content of the note, required
	Category  string `json:"category,omitempty"`          // Category of the note, optional
	Published bool   `json:"published,omitempty"`         // Published status of the note, optional
}

// UpdateNoteSchema represents the schema for updating an existing note
type UpdateNoteSchema struct {
	Title     string `json:"title,omitempty"`     // Title of the note, optional
	Content   string `json:"content,omitempty"`   // Content of the note, optional
	Category  string `json:"category,omitempty"`  // Category of the note, optional
	Published *bool  `json:"published,omitempty"` // Published status of the note, optional (pointer to handle null value)
}
