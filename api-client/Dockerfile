FROM golang:alpine

ARG API_PORT
ENV API_PORT=${API_PORT}

WORKDIR /app

COPY . .

RUN go build -o main .

EXPOSE ${API_PORT}

CMD ["./main"]