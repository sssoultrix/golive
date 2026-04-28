FROM golang:1.26-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG CMD_PATH=./cmd/profile-service
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/profile-service ${CMD_PATH}

FROM alpine:latest
COPY --from=build /out/profile-service /profile-service
EXPOSE 8081
ENTRYPOINT ["/profile-service"]