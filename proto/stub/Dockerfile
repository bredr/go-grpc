ARG PROTO
FROM ${PROTO} as proto
FROM node:lts-alpine3.10

WORKDIR /app
COPY *.json /app
RUN npm ci

COPY ./src /app/src

COPY --from=proto /go/src/proto /

ONBUILD CMD ["npm", "start"]