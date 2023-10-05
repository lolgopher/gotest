# build stage
FROM golang:1.20 as builder

ENV CGO_ENABLED=0
WORKDIR /build
COPY . ./

RUN go get
RUN go mod tidy
RUN go mod download

RUN go test
RUN go build -o app .

RUN strip /build/app

# ---
FROM alpine

COPY --from=builder /build/app .

ENTRYPOINT ["./app"]
