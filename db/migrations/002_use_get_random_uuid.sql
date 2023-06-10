ALTER TABLE tasks
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

DROP EXTENSION IF EXISTS "uuid-ossp";

---- create above / drop below ----

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE tasks
    ALTER COLUMN id SET DEFAULT uuid_generate_v4();
