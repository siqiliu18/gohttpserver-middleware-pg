FROM golang:latest

WORKDIR /app
COPY main ./main
COPY app ./app
COPY middleware ./middleware 

COPY go.mod ./

RUN go mod tidy -e

RUN cd main && go build -o appexec

CMD ["/app/main/appexec"]