# Data modelling with cassandra

Schema for `hotel` keyspace.

```cql
CREATE KEYSPACE hotel WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 3};

CREATE TYPE hotel.address (
	street text,
	city text,
	state_or_province text,
	postal_code text,
	country text
);

CREATE TABLE hotel.hotels_by_poi (
	poi_name text,
	hotel_id text,
	name text,
	phone text,
	address frozen<address>,
	PRIMARY KEY ((poi_name), hotel_id)
) WITH comment = 'Q1. Find hotels near given poi'
AND CLUSTERING ORDER BY (hotel_id ASC) ;

CREATE TABLE hotel.hotels (
	id text PRIMARY KEY,
	name text,
	phone text,
	address frozen<address>,
	pois set<text>
) WITH comment = 'Q2. Find information about a hotel';

CREATE TABLE hotel.pois_by_hotel (
	poi_name text,
	hotel_id text,
	description text,
	PRIMARY KEY ((hotel_id), poi_name)
) WITH comment = 'Q3. Find pois near a hotel';

CREATE TABLE hotel.available_rooms_by_hotel_date (
	hotel_id text,
	date date,
	room_number smallint,
	is_available boolean,
	PRIMARY KEY ((hotel_id), date, room_number)
) WITH comment = 'Q4. Find available rooms by hotel date';

CREATE TABLE hotel.amenities_by_room (
	hotel_id text,
	room_number smallint,
	amenity_name text,
	description text,
	PRIMARY KEY ((hotel_id, room_number), amenity_name)
) WITH comment = 'Q5. Find amenities for a room';
```

Schema for `reservation` keyspace.

```cql
CREATE KEYSPACE reservation WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 3};

CREATE TYPE reservation.address (
	street text,
	city text,
	state_or_province text,
	postal_code text,
	country text
);

CREATE TABLE reservation.reservations_by_confirmation (
	confirm_number text,
	hotel_id text,
	start_date date,
	end_date date,
	room_number smallint,
	guest_id uuid,
	PRIMARY KEY (confirm_number)
) WITH comment = 'Q6. Find reservations by confirmation number';

CREATE TABLE reservation.reservations_by_hotel_date (
	hotel_id text,
	start_date date,
	end_date date,
	room_number smallint,
	confirm_number text,
	guest_id uuid,
	PRIMARY KEY ((hotel_id, start_date), room_number)
) WITH comment = 'Q7. Find reservations by hotel and date';

CREATE TABLE reservation.reservations_by_guest (
	guest_last_name text,
	hotel_id text,
	start_date date,
	end_date date,
	room_number smallint,
	confirm_number text,
	guest_id uuid,
	PRIMARY KEY ((guest_last_name), hotel_id)
) WITH comment = 'Q8. Find reservations by guest name';

CREATE TABLE reservation.guests (
	guest_id uuid PRIMARY KEY,
	first_name text,
	last_name text,
	title text,
	emails set<text>,
	phone_numbers list<text>,
	addresses map<text, frozen<address>>,
	confirm_number text
) WITH comment = 'Q9. Find guest by ID.';
```

Insert queries:
```cql
INSERT INTO hotel.hotels_by_poi (poi_name, hotel_id, name, phone, address)
VALUES ('Sunway Resort', '1', 'Sunway Hotel', '012-3456789', {country: 'Malaysia'});

INSERT INTO hotel.hotels_by_poi JSON '{
	"poi_name": "Sunway Playground",
	"hotel_id": "1",
	"name": "Sunway Hotel",
	"phone": "012-3456789",
	"address": {
		"country": "Malaysia"
	}
}';
```

Select queries:

```cql
SELECT *
FROM hotel.hotels_by_poi
WHERE poi_name = 'Sunway Resort';
```

## References:
- https://cassandra.apache.org/doc/latest/data_modeling/data_modeling_schema.html
