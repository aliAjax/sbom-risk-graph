FROM golang:1.23 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/riskd ./cmd/riskd
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/riskd /riskd
EXPOSE 8123
ENTRYPOINT ["/riskd","-addr",":8123","-data-dir","/data"]
