CREATE TABLE IF NOT EXISTS riders (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
	owner_id INTEGER NOT NULL,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  number INTEGER NOT NULL,
  team TEXT,
  bike_brand TEXT,
  class TEXT,
  nationality TEXT,
  date_of_birth DATE,
  career_points INTEGER DEFAULT 0,
  status TEXT DEFAULT 'active'
);