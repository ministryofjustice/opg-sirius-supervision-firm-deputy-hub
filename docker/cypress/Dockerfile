FROM cypress/included:16.1.1@sha256:e652cda84399aaf99c68d1d0322085dc3faa7ca9ace88d9bb4f0f045d2ce6764

WORKDIR /root

RUN npm install "cypress-failed-log@2.9.2"
RUN npm install axe-core
RUN npm install cypress-axe

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://firm-deputy-hub:8888

COPY cypress.config.js .
COPY cypress cypress
