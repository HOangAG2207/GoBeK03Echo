FROM golang:alpine AS build

RUN mkdir -p /opt/app

WORKDIR /opt/app

COPY . .

RUN go build -o GoBeK03-Echo cmd/api/main.go

FROM alpine

WORKDIR /app

COPY --from=build /opt/app/GoBeK03-Echo /app/GoBeK03-Echo
COPY --from=build /opt/app/docs /app/docs

CMD ["/app/GoBeK03-Echo"]