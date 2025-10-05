FROM golang:1.25.1

WORKDIR /app

# Install Air for live reloading
RUN go install github.com/air-verse/air@latest

# Copy go.mod and go.sum first to leverage Docker layer caching
COPY go.* ./

RUN go mod download

# Copy the rest of the sources
COPY . .

EXPOSE 8080

# Use Air for live reloading in development
CMD ["air", "-c", ".air.toml"]