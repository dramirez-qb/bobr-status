FROM golang:alpine as builder
RUN mkdir /build
COPY . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o status .

FROM scratch
COPY --from=builder /build/status /app/
ENTRYPOINT [ "/app/status" ]
