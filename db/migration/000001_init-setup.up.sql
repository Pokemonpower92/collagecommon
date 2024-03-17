CREATE TABLE IF NOT EXISTS imagesets(
    id serial PRIMARY KEY,
    name VARCHAR (250) UNIQUE NOT NULL,
    description VARCHAR (250) NOT NULL
);

CREATE TABLE IF NOT EXISTS average_colors(
    id serial PRIMARY KEY,
    imageset_id INTEGER REFERENCES imagesets(id) ON DELETE CASCADE,
    R FLOAT NOT NULL,
    G FLOAT NOT NULL,
    B FLOAT NOT NULL,
    A FLOAT NOT NULL
);
