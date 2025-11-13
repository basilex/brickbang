package utility

import (
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorRed    = "\033[31m"
	colorGray   = "\033[90m"
)

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return colorGreen
	case "POST":
		return colorYellow
	case "PUT":
		return colorBlue
	case "DELETE":
		return colorRed
	default:
		return colorGray
	}
}

func InspectRoutes(app *fiber.App, useColors bool) {
	routes := app.GetRoutes()

	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Path < routes[j].Path
	})

	fmt.Println("Registered Routes Summary")
	fmt.Println(strings.Repeat("-", 120))
	fmt.Printf("%-8s | %-45s | %-60s\n", "METHOD", "PATH", "HANDLER")
	fmt.Println(strings.Repeat("-", 120))

	for _, r := range routes {
		handlerName := r.Name

		if handlerName == "" && len(r.Handlers) > 0 {
			fn := runtime.FuncForPC(reflect.ValueOf(r.Handlers[len(r.Handlers)-1]).Pointer())
			if fn != nil {
				handlerName = fn.Name()
				handlerName = strings.TrimPrefix(handlerName, "github.com/")
				handlerName = strings.TrimPrefix(handlerName, "fiber/v2.")
			}
		}

		if handlerName == "" {
			handlerName = "<anonymous>"
		}

		if useColors {
			color := colorForMethod(r.Method)
			fmt.Printf("%s%-8s%s | %-45s | %-60s\n", color, r.Method, colorReset, r.Path, handlerName)
		} else {
			fmt.Printf("%-8s | %-45s | %-60s\n", r.Method, r.Path, handlerName)
		}
	}

	fmt.Println(strings.Repeat("-", 120))
	fmt.Printf("Total routes: %d\n", len(routes))
}
