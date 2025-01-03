FROM golang:latest AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go get go.einride.tech/pid

COPY *.go ./
COPY ./apiClient/. ./apiClient/
COPY ./models/. ./models/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /SensiboPidGo

# Deploy the application binary into a lean image
FROM golang:alpine AS build-release-stage

WORKDIR /app

COPY --from=build-stage /SensiboPidGo /app/SensiboPidGo
RUN ls -la

ENTRYPOINT ["/app/SensiboPidGo"]