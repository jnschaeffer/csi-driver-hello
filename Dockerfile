FROM golang:1.21.0-bookworm AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -v -mod=readonly

FROM alpine:3.18.3 AS runner

COPY --from=build /build/csi-driver-hello /app/csi-driver-hello

CMD "/app/csi-driver-hello"

