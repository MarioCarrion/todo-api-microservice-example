CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE priority AS ENUM ('none', 'low', 'medium', 'high');

CREATE TABLE tasks (
  id          UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  description VARCHAR NOT NULL,
  priority    priority DEFAULT 'none' NOT NULL,
  start_date  TIMESTAMP WITHOUT TIME ZONE,
  due_date    TIMESTAMP WITHOUT TIME ZONE,
  done        BOOLEAN NOT NULL DEFAULT FALSE
);

---- create above / drop below ----

DROP TABLE tasks;

DROP TYPE priority;

DROP EXTENSION IF EXISTS "uuid-ossp";
