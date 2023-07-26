-- name: CreateUser :one
insert into users (user_id, client_id, reg_date) values ($1, $2, $3) returning users.id;
-- name: GetUser :one
select * from users where id = $1 limit 1;
-- name: DeleteUser :exec
delete from users where id = $1;
-- name: GetUserByUserId :one
select * from users where user_id = $1 limit 1;