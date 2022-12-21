FROM golang:1.18-alpine
WORKDIR /root/eli5-backend
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o ./bin/eli5 ./
EXPOSE 8080
CMD [ "./bin/eli5" ]