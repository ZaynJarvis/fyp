# build environment
FROM node:14-alpine as build-deps
WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn
COPY . ./

EXPOSE 8080

CMD yarn run dev

