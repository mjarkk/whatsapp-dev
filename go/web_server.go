package src

import (
	"embed"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// ErrorResponse is the response send by the server when an error occurs
type ErrorResponse struct {
	Error string `json:"error"`
}

// StartWebserverOptions are the options for the webserver
type StartWebserverOptions struct {
	Addr              string
	BasicAuthUsername string
	BasicAuthPassword string
	Rand              *rand.Rand
	Dist              embed.FS
}

// StartWebserver starts the webserver
func StartWebserver(opts StartWebserverOptions) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	authMiddleware := func(c *fiber.Ctx) error {
		return c.Next()
	}

	if opts.BasicAuthUsername != "" || opts.BasicAuthPassword != "" {
		authMiddleware = basicauth.New(basicauth.Config{
			Users: map[string]string{
				opts.BasicAuthUsername: opts.BasicAuthPassword,
			},
		})
	}

	app.Use(compress.New())
	app.Use(cors.New())
	app.Use(logger.New())

	apiRoutes(app.Group("/api", authMiddleware))
	mockRoutes(app.Group(""))

	app.Group("", authMiddleware, filesystem.New(filesystem.Config{
		PathPrefix: "dist",
		Index:      "index.html",
		Root:       http.FS(opts.Dist),
	}))

	fmt.Println("Running Web server at", opts.Addr)
	log.Fatal(app.Listen(opts.Addr))
}
