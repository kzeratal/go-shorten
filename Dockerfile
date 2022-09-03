FROM golang:latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o /goshorten

ENV REDIS_ENDPOINT=${REDIS_ENDPOINT}
ENV REDIS_PASSWORD=${REDIS_PASSWORD}
EXPOSE 8080

CMD [ "/goshorten" ]