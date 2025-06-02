package handlers

import (
	"log"
	"net/http"
	"one1-be-chal/internal/adapters/config"
	"one1-be-chal/internal/adapters/helpers"
	"strings"
	"time"

	"github.com/labstack/echo"
)

func EchoMiddleware() *echo.Echo {
	app := echo.New()
	app.HTTPErrorHandler = func(err error, context echo.Context) {
		if he, ok := err.(*echo.HTTPError); ok {
			if he.Code == http.StatusMethodNotAllowed {
				context.JSON(http.StatusMethodNotAllowed,
					echo.Map{
						"code":    http.StatusMethodNotAllowed,
						"message": "Method Not Allowed",
					},
				)
				return
			}

			if he.Code == http.StatusNotFound {
				context.JSON(http.StatusNotFound,
					echo.Map{
						"code":    http.StatusNotFound,
						"message": "Not Found",
					},
				)
				return
			}

			app.DefaultHTTPErrorHandler(err, context)
		}
	}
	app.Use(LoggerMiddleware)
	return app
}

// Logging middleware that logs HTTP method, path, and execution time.
func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func(start time.Time) {
			log.Println(c.Request().Method, c.Request().URL.Path, time.Since(start))
		}(time.Now())
		return next(c)
	}
}

func JWTMiddleware(config *config.Container) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				return c.JSON(
					http.StatusUnauthorized,
					echo.Map{"error": "Missing or invalid token"},
				)
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")

			claims, err := helpers.ParseJWT(tokenStr, *config)
			if err != nil {
				return c.JSON(
					http.StatusUnauthorized,
					echo.Map{"error": "Invalid token"},
				)
			}
			c.Set("claims", claims)

			return next(c)
		}
	}
}
