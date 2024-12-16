# Variables
DOCKER_COMPOSE=docker-compose
DOCKER_RUN=docker run -v $(shell pwd)/migrations:/migrations --network=host migrate/migrate
DATABASE_URL="postgres://user:password@localhost:5432/client_go_service?sslmode=disable"

# Command to start the application
up:
	$(DOCKER_COMPOSE) up --build -d
	@echo "Application started successfully."

# Command to run migrations
migrate:
	$(DOCKER_RUN) -path=/migrations -database $(DATABASE_URL) up
	@echo "Migrations applied successfully."

# Command to insert clients from CSV
insert_csv_clients:
	$(DOCKER_COMPOSE) exec app go run cmd/insert_clients/main.go
	@echo "CSV data loaded successfully."

# Command to stop the application
down:
	$(DOCKER_COMPOSE) down
	@echo "Application stopped and containers removed."

# Command to run all steps (up, migrate, load_csv)
setup: up migrate load_csv
	@echo "Application set up successfully."

# Command to clean Docker volumes
clean:
	$(DOCKER_COMPOSE) down -v
	@echo "Volumes cleaned and containers removed."

check:
	@which docker-compose > /dev/null || (echo "docker-compose not found! Install it and try again." && exit 1)
	@which docker > /dev/null || (echo "Docker not found! Install it and try again." && exit 1)
	@echo "Docker and Docker Compose verified."
