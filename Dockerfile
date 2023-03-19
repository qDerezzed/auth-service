FROM golang:1.19

RUN go version
ENV GOPATH=/

# copy all files
COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

# build go app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/app/main.go

CMD ["./auth-service"]