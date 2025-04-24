package main

import (
	"net/http"
	"strconv"

	"github.com/bcantrell1/pro-motocross-api/internal/database"
	"github.com/gin-gonic/gin"
)

func (app *application) createRider(c *gin.Context) {
	var rider database.Rider

	if err := c.ShouldBindJSON(&rider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := app.models.Riders.Insert(&rider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create the rider."})
		return
	}

	c.JSON(http.StatusCreated, rider)
}

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

func (app *application) getAllRiders(c *gin.Context) {
	riders, err := app.models.Riders.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sever failed to get all riders."})
	}

	c.JSON(http.StatusOK, riders)
}

func (app *application) updateRider(c *gin.Context) {
	id, err := strconv.Atoi("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rider Id."})
		return
	}

	existingRider, err := app.models.Riders.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Get rider Id."})
		return
	}

	if existingRider == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rider was not found."})
	}

	updatedRider := &database.Rider{}

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

func (app *application) deleteRider(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rider Id."})
		return
	}

	if err := app.models.Riders.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the rider."})
	}

	c.JSON(http.StatusNoContent, nil)
}
