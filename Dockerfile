FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o /app/build/estmate-ease /app/cmd/api/

FROM alpine:latest 
COPY --from=builder /app/build/estmate-ease /app/build/estmate-ease

ENV PORT=$PORT

ENTRYPOINT [ "/app/build/estmate-ease" ]