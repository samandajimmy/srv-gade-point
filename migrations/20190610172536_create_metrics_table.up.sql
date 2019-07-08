-- Table: metrics
create table metrics (
  id SERIAL PRIMARY KEY NOT NULL,
  job VARCHAR(75) NOT NULL,
  counter INT DEFAULT 0,
  status VARCHAR(2),
  creation_time TIMESTAMP DEFAULT NULL,
  modification_time TIMESTAMP DEFAULT NULL
);
CREATE INDEX index_metrics ON metrics (job);