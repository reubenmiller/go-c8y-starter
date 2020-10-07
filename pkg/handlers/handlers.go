package handlers

import (
	"fmt"
	"net/http"

	"github.com/reubenmiller/go-c8y/pkg/c8y"

	"github.com/labstack/echo"
	"github.com/reubenmiller/c8y-microservice-starter/internal/model"
)

// RegisterHandlers registers the http handlers to the given echo server
func RegisterHandlers(e *echo.Echo) {
	e.Add("GET", "/greeter", GreeterHandler)
	e.Add("GET", "/deviceByName", GetDeviceByNameHandler)
}

// GreeterHandler handles the greeter endpoint
func GreeterHandler(c echo.Context) error {
	name := c.QueryParam("name")
	return c.JSON(http.StatusOK, echo.Map{
		"message": fmt.Sprintf("Hello %s", name),
	})
}

// GetDeviceByNameHandler returns a managed object by its name
func GetDeviceByNameHandler(c echo.Context) error {
	cc := c.(*model.RequestContext)
	name := c.QueryParam("name")

	col, _, err := cc.Microservice.Client.Inventory.GetManagedObjects(
		cc.Microservice.WithServiceUser(),
		&c8y.ManagedObjectOptions{
			Query: fmt.Sprintf("name eq '%s'", name),
		},
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error":   "Could not retrieve managed object by name",
			"details": err,
		})
	}

	if len(col.ManagedObjects) == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}

	return c.JSON(http.StatusOK, col.ManagedObjects[0])
}
