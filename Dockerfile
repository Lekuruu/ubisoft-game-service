FROM golang:1.22.6

# Copy the source code
WORKDIR /app
COPY . .

# Download modules
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o ./gs

# Run the compiled binary
ENTRYPOINT ["/app/gs"]