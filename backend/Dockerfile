FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates tzdata && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY server .
COPY migrations/ ./migrations/
EXPOSE 8083
CMD ["./server"]
