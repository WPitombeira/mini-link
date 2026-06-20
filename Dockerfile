FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN go test ./... && go build -trimpath -ldflags="-s -w" -o /out/minilink ./cmd/minilink

FROM scratch
COPY --from=build /out/minilink /minilink
COPY examples/mini-link.yaml /mini-link.yaml
EXPOSE 8080
ENTRYPOINT ["/minilink"]
CMD ["serve", "-addr", ":8080", "-config", "/mini-link.yaml"]
