FROM node:26.12.1-alpine3.24@sha256:40cd011b55415bc894e45b02fdd5a5dfde056921b7d23a64c74f80497af48797 AS json-server

WORKDIR /app

EXPOSE 3000

COPY json-server .

RUN npm ci --ignore-scripts

ENTRYPOINT ["npm", "run", "serve"]
