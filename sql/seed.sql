-- name: SeedDatabase :exec
do $$
begin
    call seed_db();
end
$$ language plpgsql;
