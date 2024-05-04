ARG WORKINGPATH=/multi-site-dashboard-go
ARG ENTRYPATH=/multi-site-dashboard-go
ARG PORT=8080
ARG CONFIG=development

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
COPY *.go ./
RUN go build -o /rest ./cmd/rest
EXPOSE ${PORT}
CMD ["/rest"]