FROM cypress/included:15.21.0@sha256:af15ca5f3b851e32360f37e1d446d6e6821fd11928c87a2fd7fad12e8495e63b

WORKDIR /root

RUN npm install "cypress-failed-log@2.9.2"
RUN npm install axe-core
RUN npm install cypress-axe

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://firm-deputy-hub:8888

COPY cypress.config.js .
COPY cypress cypress
