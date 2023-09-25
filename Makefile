run:
		@go run main.go
create:
		@go run ./scripts/seed.go
build:
	docker build -t booking .
start:
	docker run -d -p 8080:8080 --name booking booking