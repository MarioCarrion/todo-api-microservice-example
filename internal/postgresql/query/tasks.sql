-- name: SelectTask :one
SELECT
  id,
  description,
  priority,
  start_date,
  due_date,
  done
FROM
  tasks
WHERE
  id = @id
LIMIT 1;

-- name: InsertTask :one
INSERT INTO tasks (
  description,
  priority,
  start_date,
  due_date
)
VALUES (
  @description,
  @priority,
  @start_date,
  @due_date
)
RETURNING id;

-- name: UpdateTask :one
UPDATE tasks SET
  description = @description,
  priority    = @priority,
  start_date  = @start_date,
  due_date    = @due_date,
  done        = @done
WHERE id = @id
RETURNING id AS res;

-- name: DeleteTask :one
DELETE FROM
  tasks
WHERE
  id = @id
RETURNING id AS res;
