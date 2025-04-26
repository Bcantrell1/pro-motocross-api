package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bcantrell1/pro-motocross-api/internal/database"
	"github.com/gin-gonic/gin"
)

// CreateRider creates a new rider
// @Summary Create a new rider
// @Description Create a new rider with the provided details
// @Tags riders
// @Accept json
// @Produce json
// @Param rider body database.Rider true "Rider data"
// @Success 201 {object} database.Rider
// @Failure 400 {object} gin.H "Invalid request body"
// @Failure 500 {object} gin.H "Failed to create the rider"
// @Router /api/v1/riders [post]
func (app *application) createRider(c *gin.Context) {
	var rider database.Rider

	if err := c.ShouldBindJSON(&rider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := app.GetUserFromContext(c)
	rider.OwnerId = user.Id

	err := app.models.Riders.Insert(&rider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create the rider."})
		fmt.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, rider)
}

// GetRider returns a single rider
// @Summary Get a rider by ID
// @Description Get details of a rider by their ID
// @Tags riders
// @Produce json
// @Param id path int true "Rider ID"
// @Success 200 {object} database.Rider
// @Failure 400 {object} gin.H "Invalid rider ID"
// @Failure 404 {object} gin.H "No rider found at that ID"
// @Failure 500 {object} gin.H "Server failed to get the requested rider"
// @Router /api/v1/riders/{id} [get]
func (app *application) getRider(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rider Id."})
	}

	rider, err := app.models.Riders.Get(id)

	if rider == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No rider found at that Id."})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server failed to get the requested rider."})
	}

	c.JSON(http.StatusOK, rider)
}

// GetAllRiders returns all riders
// @Summary Get all riders
// @Description Get a list of all riders
// @Tags riders
// @Produce json
// @Success 200 {array} database.Rider
// @Failure 500 {object} gin.H "Server failed to get all riders"
// @Router /api/v1/riders [get]
func (app *application) getAllRiders(c *gin.Context) {
	riders, err := app.models.Riders.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sever failed to get all riders."})
	}

	c.JSON(http.StatusOK, riders)
}

// UpdateRider updates an existing rider
// @Summary Update a rider
// @Description Update an existing rider with the provided details
// @Tags riders
// @Accept json
// @Produce json
// @Param id path int true "Rider ID"
// @Param rider body database.Rider true "Updated rider data"
// @Success 200 {object} database.Rider
// @Failure 400 {object} gin.H "Invalid rider ID or request body"
// @Failure 404 {object} gin.H "Rider not found"
// @Failure 500 {object} gin.H "Failed to update rider"
// @Router /api/v1/riders/{id} [put]
func (app *application) updateRider(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rider Id."})
		return
	}

	user := app.GetUserFromContext(c)
	existingRider, err := app.models.Riders.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Get rider Id."})
		return
	}

	if existingRider == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rider was not found."})
	}

	if existingRider.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update a rider you don't own."})
		return
	}

	updatedRider := &database.Rider{
		Id: id,
	}

	if err := c.ShouldBindJSON(updatedRider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedRider.Id = id

	if err := app.models.Riders.Update(updatedRider); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update rider!"})
		return
	}

	c.JSON(http.StatusOK, updatedRider)
}

// DeleteRider deletes a rider
// @Summary Delete a rider
// @Description Delete a rider by their ID
// @Tags riders
// @Param id path int true "Rider ID"
// @Success 204 "No Content"
// @Failure 400 {object} gin.H "Invalid rider ID"
// @Failure 500 {object} gin.H "Failed to delete the rider"
// @Router /api/v1/riders/{id} [delete]
func (app *application) deleteRider(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rider Id."})
		return
	}

	user := app.GetUserFromContext(c)
	existingRider, err := app.models.Riders.Get(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve the requested rider."})
		return
	}

	if user.Id != existingRider.Id {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to delete the rider!"})
		return
	}

	if err := app.models.Riders.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the rider."})
	}

	c.JSON(http.StatusNoContent, nil)
}
