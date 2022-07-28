# syntax=docker/dockerfile:1

# create a secure and minimal image using the builder pattern
ARG GOOS=darwin
ARG GOARCH=amd64
ARG GO_VERSION=1.18.3

# this image is multi-arch and the build should support a few os/archs
FROM golang:${GO_VERSION}-alpine as builder

# create a non-privileged user to run our app. the user will have permissions to the JSON file
# alternatively, the nobody user can also be used (65534) if no local storage was used
RUN adduser \
  --uid 10003 \
  --disabled-password \
  --gecos "" \
  --no-create-home \
  notroot notroot

WORKDIR /app

# download modules first to cache proceeding steps
COPY go.mod go.sum ./
RUN go mod download

# build app
COPY . .
RUN CGO_ENABLED=0 GOOS="${GOOS}" GOARCH="${GOARCH}" go build -o out/replica-count-manager ./cmd/replica-count-manager

FROM scratch
WORKDIR /
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder --chown=10003:10003 /app/out/replica-count-manager replica-count-manager
USER notroot
ENTRYPOINT ["/replica-count-manager"]
