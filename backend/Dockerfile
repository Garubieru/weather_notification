FROM golang:1.22 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM build as production

WORKDIR /usr/app
COPY --from=build /app/main .
COPY --from=build /app/.env .

EXPOSE 8001
CMD ["./main"]