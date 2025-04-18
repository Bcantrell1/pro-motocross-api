package database

import "database/sql"

type Models struct {
	Users  UserModel
	Riders RiderModel
	Events EventModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:  UserModel{DB: db},
		Riders: RiderModel{DB: db},
		Events: EventModel{DB: db},
	}
}
