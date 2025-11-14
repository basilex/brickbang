// brickbang/utility/validation.go
package utility

import (
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ValidateBody(ctx *fiber.Ctx, dst interface{}, v *validator.Validate) error {
	if err := ctx.BodyParser(dst); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "bad JSON: "+err.Error())
	}

	if err := v.Struct(dst); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, formatValidationErrors(err))
	}

	return nil
}

func ParseID(ctx *fiber.Ctx, name string) (string, error) {
	id := ctx.Params(name)
	if id == "" {
		return "", fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("%s required", name))
	}

	return id, nil
}

func ParseIntQuery(ctx *fiber.Ctx, key string, def int) (int, error) {
	valStr := ctx.Query(key, strconv.Itoa(def))

	v, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fiber.NewError(
			fiber.StatusBadRequest, fmt.Sprintf("bad %s: %v", key, err),
		)
	}

	return v, nil
}

func formatValidationErrors(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		m := make(map[string]string)

		for _, e := range errs {
			m[e.Field()] = fmt.Sprintf("неверный формат: %s", e.Tag())
		}

		return fmt.Sprint(m)
	}

	return err.Error()
}
