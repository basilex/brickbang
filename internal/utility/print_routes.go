package utility

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// PrintRoutes prints all registered routes in a formatted table
func PrintRoutes(app *fiber.App) {
	routes := app.GetRoutes()

	fmt.Println("Registered routes:")
	fmt.Println("-------------------------------------------------------------")
	fmt.Printf("%-8s %-40s %-30s\n", "METHOD", "PATH", "NAME")
	fmt.Println("-------------------------------------------------------------")

	for _, r := range routes {
		fmt.Printf("%-8s %-40s %-30s\n", r.Method, r.Path, r.Name)
	}

	fmt.Println("-------------------------------------------------------------")
	fmt.Printf("Total routes: %d\n", len(routes))
}
