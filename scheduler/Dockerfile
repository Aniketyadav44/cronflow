FROM golang:1.24

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o scheduler .

EXPOSE 8000

CMD [ "./scheduler" ]