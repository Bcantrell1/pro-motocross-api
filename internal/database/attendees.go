package database

import (
	"context"
	"database/sql"
	"time"
)

type AttendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	Id      int `json:"id"`
	RiderId int `json:"riderId"`
	EventId int `json:"eventId"`
}

func (m *AttendeeModel) Insert(attendee *Attendee) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO attendees (event_id, rider_id) VALUES ($1, $2) RETURNING id"
	err := m.DB.QueryRowContext(ctx, query, attendee.EventId, attendee.RiderId).Scan(&attendee.Id)

	if err != nil {
		return nil, err
	}

	return attendee, nil
}

func (m *AttendeeModel) Delete(riderId, eventId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM attendees WHERE rider_id = $1 AND event_id = $2"
	_, err := m.DB.ExecContext(ctx, query, riderId, eventId)
	if err != nil {
		return err
	}
	return nil
}

func (m *AttendeeModel) GetEventsByAttendee(attendeeId int) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT e.id, e.owner_id, e.name, e.description, e.date, e.location
		FROM events e
		JOIN attendees a ON e.id = a.event_id
		WHERE a.rider_id = $1
	`
	rows, err := m.DB.QueryContext(ctx, query, attendeeId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil

}

func (m *AttendeeModel) GetByEventAndAttendee(eventId, riderId int) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM attendees WHERE event_id = $1 AND rider_id = $2"
	var attendee Attendee

	err := m.DB.QueryRowContext(ctx, query, eventId, riderId).Scan(&attendee.Id, &attendee.RiderId, &attendee.EventId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &attendee, nil
}

func (m AttendeeModel) GetAttendeesByEvent(eventId int) ([]Rider, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
     SELECT r.id, r.first_name, r.last_name
     FROM riders r
     JOIN attendees a ON r.id = a.rider_id
     WHERE a.event_id = $1
 `
	rows, err := m.DB.QueryContext(ctx, query, eventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var riders []Rider
	for rows.Next() {
		var rider Rider
		err := rows.Scan(&rider.Id, &rider.FirstName, &rider.LastName)
		if err != nil {
			return nil, err
		}
		riders = append(riders, rider)
	}
	return riders, nil
}
