FROM golang:latest
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main
EXPOSE 8080
ENV MONGODB_URI mongodb://mongodb_host:27017
ENV MONGODB_DATABASE hotel-reservation
CMD ["./main"]