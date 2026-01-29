-- name: CreateTodo :one
INSERT INTO
    todos (title, completed)
VALUES
    ($1, $2)
RETURNING
    *;

-- name: GetTodo :one
SELECT
    *
FROM
    todos
WHERE
    id = $1
LIMIT
    1;

-- name: ListTodos :many
SELECT
    *
FROM
    todos
ORDER BY
    id;

-- name: UpdateTodoTitle :one
UPDATE todos
SET
    title = $2,
    updated_at = NOW()
WHERE
    id = $1
RETURNING
    *;

-- name: UpdateTodoCompleted :one
UPDATE todos
SET
    completed = $2,
    updated_at = NOW()
WHERE
    id = $1
RETURNING
    *;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE
    id = $1;