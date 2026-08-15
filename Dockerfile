FROM golang:1.24 AS build
WORKDIR /src
COPY control-plane/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /atlas-api ./cmd/atlas-api
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /atlas-api /atlas-api
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/atlas-api"]
