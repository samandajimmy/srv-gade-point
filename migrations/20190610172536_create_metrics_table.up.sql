-- Table: metrics
create table metrics (
  id SERIAL PRIMARY KEY NOT NULL,
  module VARCHAR(255) NOT NULL,
  counter INT DEFAULT 0
);
CREATE INDEX index_metrics ON metrics (module);