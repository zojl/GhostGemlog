FROM golang:alpine AS builder
WORKDIR /app
ADD go.mod .
COPY . . 
RUN go build -o app ./main.go

FROM alpine
RUN apk update && \
    apk add --no-cache curl
WORKDIR /app
ARG GEMINI_PORT
ENV GEMINI_PORT 1965
COPY --from=builder /app/app /app/app
HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
    CMD nc -z localhost $PORT || exit 1
CMD ["/app/app"]