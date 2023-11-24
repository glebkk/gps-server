pg_dump:
	docker exec -i db_gps /bin/bash -c "PGPASSWORD=123 pg_dump --username postgres postgres" > ./dump.sql