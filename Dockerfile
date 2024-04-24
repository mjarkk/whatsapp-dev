# Build the frontend container
FROM node:21-slim AS frontend

COPY . /app
WORKDIR /app

RUN npm i && npm run build

# Build the runtime container
FROM golang:1.21

WORKDIR /usr/src/app

COPY go.mod go.sum .
RUN go mod download

COPY go go
COPY main.go .
COPY --from=frontend /app/dist dist

RUN go build -o /usr/bin/app && \
    rm -rf go main.go go.mod go.sum dist

EXPOSE 1090/tcp

CMD ["/usr/bin/app"]
