## Build front end assets so json-server can route to them
FROM node:14.17.5-alpine3.14 as build-assets

WORKDIR /app

RUN mkdir -p web/static

COPY web/assets web/assets
COPY webpack.config.js .
COPY package.json .
COPY yarn.lock .

RUN yarn install
RUN yarn build

## build json-server for mock API
FROM node:14.17.5-alpine3.14 as json-server

WORKDIR /app

EXPOSE 3000

COPY json-server .

RUN yarn install

COPY --from=build-assets /app/web/static/ public

ENTRYPOINT ["yarn", "serve"]