FROM node:current-alpine AS build-frontend
LABEL maintainer Ascensio System SIA <support@onlyoffice.com>
ARG BACKEND_GATEWAY
ARG PIPEDRIVE_CREATE_MODAL_ID
ARG PIPEDRIVE_EDITOR_MODAL_ID
ENV BACKEND_GATEWAY=$BACKEND_GATEWAY \
    PIPEDRIVE_CREATE_MODAL_ID=$PIPEDRIVE_CREATE_MODAL_ID \
    PIPEDRIVE_EDITOR_MODAL_ID=$PIPEDRIVE_EDITOR_MODAL_ID
WORKDIR /usr/src/app
COPY ./frontend/package*.json ./
RUN npm install
COPY frontend .
RUN npm run build

FROM golang:alpine AS build-gateway
WORKDIR /usr/src/app
COPY backend .
RUN go build services/gateway/main.go

FROM golang:alpine AS build-auth
WORKDIR /usr/src/app
COPY backend .
RUN go build services/auth/main.go

FROM golang:alpine AS build-builder
WORKDIR /usr/src/app
COPY backend .
RUN go build services/builder/main.go

FROM golang:alpine AS build-callback
WORKDIR /usr/src/app
COPY backend .
RUN go build services/callback/main.go

FROM golang:alpine AS build-settings
WORKDIR /usr/src/app
COPY backend .
RUN go build services/settings/main.go

FROM golang:alpine AS gateway
WORKDIR /usr/src/app
COPY --from=build-gateway \
     /usr/src/app/main \
     /usr/src/app/main
EXPOSE 4044
CMD ["./main", "server"]

FROM golang:alpine AS auth
WORKDIR /usr/src/app
COPY --from=build-auth \
     /usr/src/app/main \
     /usr/src/app/main
EXPOSE 5052
CMD ["./main", "server"]

FROM golang:alpine AS builder
WORKDIR /usr/src/app
COPY --from=build-builder \
     /usr/src/app/main \
     /usr/src/app/main
EXPOSE 6260
CMD ["./main", "server"]

FROM golang:alpine AS callback
WORKDIR /usr/src/app
COPY --from=build-callback \
     /usr/src/app/main \
     /usr/src/app/main
EXPOSE 5044
CMD ["./main", "server"]

FROM golang:alpine AS settings
WORKDIR /usr/src/app
COPY --from=build-settings \
     /usr/src/app/main \
     /usr/src/app/main
EXPOSE 5150
CMD ["./main", "server"]

FROM nginx:alpine AS frontend
COPY --from=build-frontend \
    /usr/src/app/build \
    /usr/share/nginx/html
EXPOSE 80
