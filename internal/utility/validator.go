package utility

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"brickbang/internal/exception"
)

// Global validator instance
var validate = validator.New()

// ParseAndValidate parses JSON body and validates struct fields using `validate` tags.
// It returns a structured Fiber error compatible with your UnifiedResponse middleware.
func ParseAndValidate(ctx *fiber.Ctx, dest any) error {
	// Parse incoming JSON
	if err := ctx.BodyParser(dest); err != nil {
		return exception.ErrBadRequest("invalid JSON body")
	}

	// Run validation rules
	if err := validate.Struct(dest); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			var messages []string
			for _, fe := range ve {
				field := strings.ToLower(fe.Field())

				switch fe.Tag() {
				case "required":
					messages = append(messages, fmt.Sprintf("%s is required", field))
				case "email":
					messages = append(messages, fmt.Sprintf("%s must be a valid email", field))
				case "min":
					messages = append(messages, fmt.Sprintf("%s must be at least %s characters", field, fe.Param()))
				case "max":
					messages = append(messages, fmt.Sprintf("%s must be at most %s characters", field, fe.Param()))
				default:
					messages = append(messages, fmt.Sprintf("%s is invalid (%s)", field, fe.Tag()))
				}
			}

			return exception.ErrUnprocessable(strings.Join(messages, "; "))
		}

		// Unknown validation error
		return exception.ErrUnprocessable(err.Error())
	}

	return nil
}
