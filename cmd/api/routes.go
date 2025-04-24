package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	// routes for v1
	v1 := g.Group("/api/v1")
	{
		v1.POST("/events", app.createEvent)
		v1.GET("/events", app.getAllEvents)
		v1.GET("/events/:id", app.getEvent)
		v1.PUT("/events/:id", app.updateEvent)
		v1.DELETE("/events/:id", app.deleteEvent)

		v1.POST("/riders", app.createRider)
		v1.GET("/riders", app.getAllRiders)
		v1.GET("/riders/:id", app.getRider)
		v1.PUT("/riders/:id", app.updateRider)
		v1.DELETE("/riders/:id", app.deleteRider)

		v1.POST("/events/:id/attendees/:riderId", app.addAttendeeToEvent)
		v1.POST("/events/:id/attendees", app.getAttendeesForEvent)

		v1.POST("/auth/register", app.registerUser)
	}

	return g
}
