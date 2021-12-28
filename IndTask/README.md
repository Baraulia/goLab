1. run postgres into container, create database:
sudo docker run --name indtask-db -p 5436:5432 -e POSTGRES_USER=baraulia -e POSTGRES_PASSWORD=qwerty -e POSTGRES_DB=baraulia -d postgres

2. create tables in the database using migration (library golang-migrate):
migrate -path ./schema -database "postgres://baraulia:qwerty@localhost:5436/baraulia?sslmode=disable" up


[//]: # (3. install redis)

[//]: # (4. create .env file - example in .env.example)