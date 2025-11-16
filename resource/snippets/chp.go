package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var db = map[string]string{} // username → bcrypt hash

func main() {
	app := fiber.New()

	app.Post("/register", func(c *fiber.Ctx) error {
		var req AuthRegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		db[req.Username] = string(hash)
		fmt.Println("REGISTER")
		fmt.Println("Username:", req.Username)
		fmt.Println("Password:", req.Password)
		fmt.Println("Stored hash:", string(hash), "len:", len(hash))
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"username": req.Username,
			"hash":     string(hash),
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		var req AuthLoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		hash, ok := db[req.Username]
		if !ok {
			return c.Status(fiber.StatusUnauthorized).SendString("user not found")
		}

		fmt.Println("LOGIN")
		fmt.Println("Username:", req.Username)
		fmt.Println("Password:", req.Password)
		fmt.Println("Stored hash:", hash, "len:", len(hash))

		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
			fmt.Println("Compare error:", err)
			return c.Status(fiber.StatusUnauthorized).SendString("invalid password")
		}

		return c.JSON(fiber.Map{
			"message": "login successful",
		})
	})

	log.Fatal(app.Listen(":8080"))
}
