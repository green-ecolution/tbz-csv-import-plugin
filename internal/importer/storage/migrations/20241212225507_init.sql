-- +goose Up
CREATE TABLE IF NOT EXISTS imports (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  user_id VARCHAR(255),
  raw_csv TEXT
);

CREATE TABLE IF NOT EXISTS trees (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  tree_number VARCHAR(255),
  species VARCHAR(255) NOT NULL DEFAULT '',
  area VARCHAR(255),
  planting_year INT,
  street VARCHAR(255),
  latitude FLOAT,
  longitude FLOAT
);

CREATE TABLE IF NOT EXISTS tree_import (
  tree_id INT,
  import_id INT,
  PRIMARY KEY (tree_id, import_id),
  FOREIGN KEY (tree_id) REFERENCES trees(id),
  FOREIGN KEY (import_id) REFERENCES imports(id)
);

-- +goose Down
DROP TABLE IF EXISTS tree_import;
DROP TABLE IF EXISTS trees;
DROP TABLE IF EXISTS imports;
