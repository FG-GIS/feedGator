psql postgres://postgres:postgres@localhost:5432/gator
sudo systemctl start postgresql
goose -dir ./sql/schema postgres "$DB_URL" up
