FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache make nodejs npm git bash \
    && go install github.com/a-h/templ/cmd/templ@latest

ENV PATH=$PATH:/go/bin

COPY ../go.mod ../go.sum ./
RUN go mod download

# COPY ../scripts/wait-for.sh /usr/local/bin/wait-for.sh
# COPY ../scripts/entrypoint.sh /usr/local/bin/entrypoint.sh

COPY ../. .
RUN make build

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/bin/app ./bin/app
COPY --from=builder /app/webapp/dist ./webapp/dist
COPY --from=builder /app/static ./static
COPY --from=builder /app/schemas ./schemas
COPY --from=builder /app/config/config.yml ./config/config.yml
# COPY --from=builder /usr/local/bin/entrypoint.sh ./scripts/entrypoint.sh
# COPY --from=builder /usr/local/bin/wait-for.sh ./scripts/wait-for.sh

COPY --from=builder /app/scripts/entrypoint.sh /usr/local/bin/entrypoint.sh
COPY --from=builder /app/scripts/wait-for.sh /usr/local/bin/wait-for.sh

RUN chmod +x /usr/local/bin/entrypoint.sh /usr/local/bin/wait-for.sh
RUN apk add --no-cache ca-certificates bash netcat-openbsd

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD []
