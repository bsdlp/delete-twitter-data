FROM golang:alpine AS builder
WORKDIR /src/
ADD . /src/
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags='-extldflags=-static' -o /deleter 

FROM gcr.io/distroless/static
COPY --from=builder /deleter /bin/deleter

ENTRYPOINT [ "/bin/deleter" ]