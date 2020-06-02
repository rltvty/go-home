package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rltvty/go-home/logwrapper"
	"go.uber.org/zap"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	log := logwrapper.GetInstance()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/endpoints", func(c echo.Context) error {
		content, err := ioutil.ReadFile("./speaker_endpoints.json")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.String(http.StatusOK, string(content))
	})

	speakerGroup := e.Group("/speaker")
	speakerGroup.Use(func (next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			speakerId := c.Param("speakerId")
			speaker, err := getSpeaker(speakerId)
			if err != nil {
				return c.String(http.StatusInternalServerError, "couldn't find speaker")
			}
			c.Set("speaker", speaker)
			return next(c)
		}
	})
	speakerGroup.POST("/:speakerId/endpoint/:endpoint/value/:value", func(c echo.Context) error {
		endpoint := c.Param("endpoint")
		value := c.Param("value")
		speaker := c.Get("speaker")
		return c.String(http.StatusAccepted, fmt.Sprintf("%s %s set to %s", speaker, endpoint, value))
	})

	// Start server
	err := e.Start(":8000")
	if err != nil {
		log.Fatal("Error starting API server", zap.Any("error", err))
	}
}

func getSpeaker(speakerId string) (interface{}, error)  {
	if len(speakerId) > 2 {
		return fmt.Sprintf("Speaker: %s", speakerId), nil
	}
	return nil, fmt.Errorf("couldn't find speaker")
}

