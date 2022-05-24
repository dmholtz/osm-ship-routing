FROM golang:1.18-alpine

WORKDIR /app

# Copy go.mod and go.sum for installing dependencies
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy source code and additional resources
COPY . .

# Build server
RUN go build -o osm-server cmd/server/main.go

EXPOSE 8081

CMD ["./osm-server"]