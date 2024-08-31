CREATE TABLE IF NOT EXISTS reservation (
	id VARCHAR(12) PRIMARY KEY,
	room_id VARCHAR NOT NULL,
	start_time TIMESTAMP NOT NULL,
	end_time TIMESTAMP NOT NULL,
	CONSTRAINT start_time_before_end_time CHECK (start_time < end_time)
);

CREATE INDEX IF NOT EXISTS reservation_room_id_idx ON reservation(room_id);