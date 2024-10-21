-- name: GetSnippet :one
select * from snippets
where id = ? limit 1;

-- name: ListSnippets :many
select * from snippets;

-- name: GetLatestSnippets :many
select * from snippets
where (expires is not null and expires > date('now')) or expires is null order by id desc limit 10;

-- name: CreateSnippet :one
insert into snippets (
    title, content, expires, created
) values (?, ?, ?, date('now'))
    returning id;

-- name: DeleteSnippet :exec
delete from snippets
where id = ?;
