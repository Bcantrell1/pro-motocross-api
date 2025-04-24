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
		v1.GET("/events", app.getAllEvents)
		v1.GET("/events/:id", app.getEvent)

		v1.GET("/riders", app.getAllRiders)
		v1.GET("/riders/:id", app.getRider)

		v1.GET("/events/:id/attendees", app.getAttendeesForEvent)
		v1.GET("/attendees/:id/events", app.getEventsByAttendee)

		v1.POST("/auth/register", app.registerUser)
		v1.POST("/auth/login", app.login)
	}

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	{
		authGroup.POST("/events", app.createEvent)
		authGroup.PUT("/events/:id", app.updateEvent)
		authGroup.DELETE("/events/:id", app.deleteEvent)

		authGroup.POST("/riders", app.createRider)
		authGroup.PUT("/riders/:id", app.updateRider)
		authGroup.DELETE("/riders/:id", app.deleteRider)

		authGroup.POST("/events/:id/attendees/:riderId", app.addAttendeeToEvent)
		authGroup.DELETE("/events/:id/attendees/:riderId", app.deleteAttendeeFromEvent)
	}

	return g
}
