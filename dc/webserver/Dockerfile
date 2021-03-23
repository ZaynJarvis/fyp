# build environment
FROM node:14-alpine as build-deps
WORKDIR /usr/src/app
COPY package.json yarn.lock ./
RUN yarn
COPY . ./

EXPOSE 3000
CMD node index.js
