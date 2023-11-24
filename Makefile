pg_dump:
	docker exec -i db_gps /bin/bash -c "PGPASSWORD=123 pg_dump --username postgres postgres" > ./dump.sql
pg_reset:
	cat .\dump.sql | docker-compose exec -T db psql -U postgres -d postgres