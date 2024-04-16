# Build the frontend container
FROM node:20-slim AS frontend

COPY . .
RUN npm install && npm run build

# Build the runtime container
FROM golang:1.21

WORKDIR /usr/src/app

COPY go.mod go.sum .
RUN go mod download

COPY go go
COPY main.go .
COPY --from=frontend dist dist

RUN go build -o /usr/bin/whatsapp-dev && \
    rm -rf go main.go go.mod go.sum dist

EXPOSE 1090/tcp

CMD ["/usr/bin/whatsapp-dev"]
