package main

import (
	"net/http"
	"strconv"

	"github.com/bcantrell1/pro-motocross-api/internal/database"
	"github.com/gin-gonic/gin"
)

// CreateEvent creates a new event
// @Summary Create a new event
// @Description Create a new event with the provided details
// @Tags events
// @Accept json
// @Produce json
// @Param event body database.Event true "Event data"
// @Success 201 {object} database.Event
// @Failure 400 {object} gin.H "Invalid request body"
// @Failure 500 {object} gin.H "Failed to create the event"
// @Router /api/v1/events [post]
func (app *application) createEvent(c *gin.Context) {
	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := app.GetUserFromContext(c)
	event.OwnerId = user.Id

	err := app.models.Events.Insert(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create the event."})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetEvent returns a single event
// @Summary Returns a single event
// @Description Returns a single event
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event Id"
// @Success 200 {object} database.Event
// @Router /api/v1/events/{id} [get]
func (app *application) getEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid events Id."})
	}

	event, err := app.models.Events.Get(id)

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "There is no event found for that Id."})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "The server failed to retrieve an event."})
	}

	c.JSON(http.StatusOK, event)
}

// GetAllEvents returns all events
// @Summary Get all events
// @Description Get a list of all events
// @Tags events
// @Produce json
// @Success 200 {array} database.Event
// @Failure 500 {object} gin.H "Server failed to get all events"
// @Router /api/v1/events [get]
func (app *application) getAllEvents(c *gin.Context) {
	events, err := app.models.Events.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sever failed to get all events."})
	}

	c.JSON(http.StatusOK, events)
}

// UpdateEvent updates an existing event
// @Summary Update an existing event
// @Description Update an existing event with the provided details
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param event body database.Event true "Updated event data"
// @Success 200 {object} database.Event
// @Failure 400 {object} gin.H "Invalid event ID or request body"
// @Failure 403 {object} gin.H "Unauthorized to update the event"
// @Failure 404 {object} gin.H "Event not found"
// @Failure 500 {object} gin.H "Failed to update event"
// @Router /api/v1/events/{id} [put]
func (app *application) updateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event Id."})
		return
	}

	user := app.GetUserFromContext(c)
	existingEvent, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Get event Id."})
		return
	}

	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event was not found."})
	}

	if existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update an event you don't own."})
		return
	}

	updatedEvent := &database.Event{
		Id: id,
	}

	if err := c.ShouldBindJSON(&updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := app.models.Events.Update(updatedEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event!"})
		return
	}

	c.JSON(http.StatusOK, updatedEvent)
}

// DeleteEvent deletes an event
// @Summary Delete an event
// @Description Delete an event by its ID
// @Tags events
// @Param id path int true "Event ID"
// @Success 204 "No Content"
// @Failure 400 {object} gin.H "Invalid event ID"
// @Failure 401 {object} gin.H "Unauthorized to delete the event"
// @Failure 404 {object} gin.H "Event not found"
// @Failure 500 {object} gin.H "Failed to delete the event"
// @Router /api/v1/events/{id} [delete]
func (app *application) deleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event Id."})
		return
	}

	user := app.GetUserFromContext(c)
	existingEvent, err := app.models.Events.Get(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve the requested event."})
		return
	}

	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "This specific event not found."})
		return
	}

	if user.Id != existingEvent.OwnerId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to delete that event!"})
		return
	}

	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the event."})
	}

	c.JSON(http.StatusNoContent, nil)
}

// AddAttendeeToEvent adds a rider to an event
// @Summary Add a rider to an event
// @Description Add a rider as an attendee to an event
// @Tags attendees
// @Param id path int true "Event ID"
// @Param riderId path int true "Rider ID"
// @Success 201 {object} database.Attendee
// @Failure 400 {object} gin.H "Invalid event or rider ID"
// @Failure 401 {object} gin.H "Unauthorized to add attendee"
// @Failure 404 {object} gin.H "Event or rider not found"
// @Failure 409 {object} gin.H "Rider already signed up for this event"
// @Failure 500 {object} gin.H "Failed to add rider to event"
// @Router /api/v1/events/{id}/attendees/{riderId} [post]
func (app *application) addAttendeeToEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event Id."})
		return
	}

	riderId, err := strconv.Atoi(c.Param("riderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rider Id."})
	}

	event, err := app.models.Events.Get(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event."})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found."})
		return
	}

	riderToAdd, err := app.models.Riders.Get(riderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rider."})
		return
	}

	user := app.GetUserFromContext(c)

	if user.Id != event.OwnerId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to add to a rider to an event you don't own."})
		return
	}

	if riderToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rider not found."})
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(event.Id, riderToAdd.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendee."})
		return
	}

	if existingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Rider already signed up for this event!"})
		return
	}

	attendee := database.Attendee{
		EventId: event.Id,
		RiderId: riderToAdd.Id,
	}

	_, err = app.models.Attendees.Insert(&attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add rider to event."})
		return
	}

	c.JSON(http.StatusCreated, attendee)
}

// GetAttendeesForEvent gets all attendees for an event
// @Summary Get attendees for an event
// @Description Get a list of riders attending an event
// @Tags attendees
// @Param id path int true "Event ID"
// @Success 200 {array} database.Rider
// @Failure 400 {object} gin.H "Invalid event ID"
// @Failure 500 {object} gin.H "Failed to retrieve attendees"
// @Router /api/v1/events/{id}/attendees [get]
func (app *application) getAttendeesForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event Id"})
		return
	}

	riders, err := app.models.Attendees.GetAttendeesByEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, riders)
}

// DeleteAttendeeFromEvent removes a rider from an event
// @Summary Remove a rider from an event
// @Description Remove a rider as an attendee from an event
// @Tags attendees
// @Param id path int true "Event ID"
// @Param riderId path int true "Rider ID"
// @Success 204 "No Content"
// @Failure 400 {object} gin.H "Invalid event or rider ID"
// @Failure 401 {object} gin.H "Unauthorized as you don't own the event"
// @Failure 500 {object} gin.H "Failed to delete attendee"
// @Router /api/v1/events/{id}/attendees/{riderId} [delete]
func (app *application) deleteAttendeeFromEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event Id."})
		return
	}

	riderId, err := strconv.Atoi(c.Param("riderId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rider Id."})
		return
	}

	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve the requested event."})
		return
	}

	user := app.GetUserFromContext(c)

	if user.Id != event.OwnerId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized as you don't own the event."})
		return
	}

	err = app.models.Attendees.Delete(riderId, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendee."})
		return
	}

	c.JSON(http.StatusNoContent, nil)

}

// GetEventsByAttendee gets all events for an attendee
// @Summary Get events for an attendee
// @Description Get a list of events that a rider is attending
// @Tags attendees
// @Param id path int true "Attendee ID"
// @Success 200 {array} database.Event
// @Failure 400 {object} gin.H "Invalid attendee ID"
// @Failure 500 {object} gin.H "Failed to retrieve events"
// @Router /api/v1/attendees/{id}/events [get]
func (app *application) getEventsByAttendee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendee Id."})
		return
	}

	events, err := app.models.Events.GetByAttendee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}
