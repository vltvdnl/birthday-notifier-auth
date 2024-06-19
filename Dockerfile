FROM golang:1.21

WORKDIR /sso

COPY go.* ./
RUN go mod download

COPY . . 
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /docker-bn-sso-app ./cmd/sso

EXPOSE 44044

CMD ["/docker-bn-sso-app"]