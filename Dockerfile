FROM docker.io/library/golang:1.13.5 as build
WORKDIR /app
COPY . /app/
RUN CGO_ENABLED=0 go build -o /go/bin/hello-world main.go

# ================
# RUNTIME STAGE
# ================
FROM gcr.io/distroless/static
WORKDIR /app

COPY --from=build --chown=nonroot:nonroot /go/bin/hello-world /app/hello-world

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/hello-world"]
