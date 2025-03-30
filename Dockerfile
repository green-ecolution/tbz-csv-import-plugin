FROM node:22-alpine AS builder-ui
WORKDIR /app/build
COPY ./ui/package.json ./ui/yarn.lock ./

RUN yarn --frozen-lockfile
COPY ./ui .

ENV VITE_BASE_URL=api/v1
RUN yarn build

FROM golang:1.23-alpine AS builder

WORKDIR /app/build

COPY ./Makefile ./go.mod ./go.sum .

RUN go mod download
RUN apk add --no-cache make proj-dev build-base

COPY . .
COPY --from=builder-ui /app/build/dist /app/build/ui/dist
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o csv-import .

FROM alpine:3.20 AS runner
RUN adduser -D gorunner
RUN apk add --no-cache proj-dev

USER gorunner
WORKDIR /app

COPY --chown=gorunner:gorunner --from=builder /app/build/csv-import /app

ENTRYPOINT [ "/app/csv-import" ]
