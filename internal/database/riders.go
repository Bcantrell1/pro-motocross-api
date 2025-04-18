package database

import "database/sql"

type RiderModel struct {
	DB *sql.DB
}

type Rider struct {
	Id           int    `json:"id"`
	FirstName    string `json:"firstName" binding:"required,min=3"`
	LastName     string `json:"lastName" binding:"required,min=3"`
	Number       int    `json:"number" binding:"required"`
	Team         string `json:"team"`
	BikeBrand    string `json:"bikeBrand"`
	Class        string `json:"class"`
	Nationality  string `json:"nationality"`
	DateOfBirth  string `json:"dateOfBirth"`
	CareerPoints int    `json:"careerPoints"`
	Status       string `json:"status"`
}
