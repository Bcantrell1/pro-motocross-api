CREATE TABLE IF NOT EXISTS event_riders (
  event_id INTEGER NOT NULL,
  rider_id INTEGER NOT NULL,
  position INTEGER,
  points INTEGER,
  qualifying_time TEXT,
  dnf BOOLEAN DEFAULT FALSE,
  dq BOOLEAN DEFAULT FALSE,
  PRIMARY KEY (event_id, rider_id),
  FOREIGN KEY (event_id) REFERENCES events (id) ON DELETE CASCADE,
  FOREIGN KEY (rider_id) REFERENCES riders (id) ON DELETE CASCADE
);