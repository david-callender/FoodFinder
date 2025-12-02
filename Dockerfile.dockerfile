# syntax=docker/dockerfile:1

FROM golang:1.25.1
WORKDIR /GopherGrub

COPY ./backend/go.mod ./backend/go.sum ./backend/
COPY ./notifier/go.mod ./notifier/go.sum ./notifier/
COPY ./scraper/go.mod ./scraper/go.sum ./scraper/
COPY ./utils/go.mod ./utils/go.sum ./utils/

# Backend
WORKDIR /GopherGrub/backend
RUN go mod download
COPY ./backend/*.go ./
RUN CGO_ENABLED=0 GOOS=linux GOEXPERIMENT=jsonv2 go build -o /docker-backend


# notifier
WORKDIR /GopherGrub/notifier
RUN go mod download
COPY ./notifier/*.go ./
RUN CGO_ENABLED=0 GOOS=linux GOEXPERIMENT=jsonv2 go build -o /docker-notifier


# Scraper
WORKDIR /GopherGrub/scraper
RUN go mod download
COPY ./scraper/*.go ./
RUN CGO_ENABLED=0 GOOS=linux GOEXPERIMENT=jsonv2 go build -o /docker-scraper


# Utils (if needed as separate module)
WORKDIR /GopherGrub/utils
RUN go mod download
COPY ./utils/*.go ./

# run compiled binary
CMD ["/docker-backend"]