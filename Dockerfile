ARG WORKINGPATH=/app
ARG ENTRYPATH=/app
ARG CONFIG=development
ARG DEPLOYMENT_IMAGE=scratch

# Stage: BUILD
# Install dependencies first to maximize Docker layer caching.
FROM golang:1.22.2 AS build
ARG WORKINGPATH
ARG CONFIG
ARG BASEHREF
WORKDIR ${WORKINGPATH}

# Install packages.
COPY go.mod go.sum ./
RUN go mod download

# Build from source code.
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /rest ./cmd/rest

# Stage: DEPLOY
FROM $DEPLOYMENT_IMAGE
COPY --from=build /app/internal/config /app/internal/config 
COPY --from=build /rest /rest
ENTRYPOINT ["/rest"]