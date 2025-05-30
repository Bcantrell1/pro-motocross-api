package database

import (
	"context"
	"database/sql"
	"time"
)

type RiderModel struct {
	DB *sql.DB
}

type Rider struct {
	Id           int    `json:"id"`
	OwnerId      int    `json:"ownerId" binding:"required"`
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

func (m *RiderModel) Insert(rider *Rider) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO riders (owner_id, first_name, last_name, number, team, bike_brand, class, nationality, date_of_birth, career_points, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id"

	return m.DB.QueryRowContext(ctx, query, rider.OwnerId, rider.FirstName, rider.LastName, rider.Number, rider.Team, rider.BikeBrand, rider.Class, rider.Nationality, rider.DateOfBirth, rider.CareerPoints, rider.Status).Scan(&rider.Id)
}

func (m *RiderModel) GetAll() ([]*Rider, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM riders"

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	riders := []*Rider{}

	for rows.Next() {
		var rider Rider

		err := rows.Scan(&rider.Id, &rider.OwnerId, &rider.FirstName, &rider.LastName, &rider.Number, &rider.Team, &rider.BikeBrand, &rider.Class, &rider.Nationality, &rider.DateOfBirth, &rider.CareerPoints, &rider.Status)
		if err != nil {
			return nil, err
		}

		riders = append(riders, &rider)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return riders, nil
}

func (m *RiderModel) Get(id int) (*Rider, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM riders WHERE id = $1"

	var rider Rider

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&rider.Id, &rider.OwnerId, &rider.FirstName, &rider.LastName, &rider.Number, &rider.Team, &rider.BikeBrand, &rider.Class, &rider.Nationality, &rider.DateOfBirth, &rider.CareerPoints, &rider.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &rider, nil
}

func (m *RiderModel) Update(rider *Rider) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "UPDATE riders SET owner_id = $1, first_name = $2, last_name = $3, number = $4, team = $5, bike_brand = $6, class = $7, nationality = $8, date_of_birth = $9, career_points = $10, status = $11 WHERE id = $12"

	_, err := m.DB.ExecContext(ctx, query, rider.OwnerId, rider.FirstName, rider.LastName, rider.Number, rider.Team, rider.BikeBrand, rider.Class, rider.Nationality, rider.DateOfBirth, rider.CareerPoints, rider.Status, rider.Id)
	if err != nil {
		return err
	}

	return nil
}

func (m *RiderModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM riders WHERE id = $1"

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
