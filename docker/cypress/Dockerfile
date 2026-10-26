FROM cypress/included:16.2.0@sha256:c196ed896aec24c7ea8fc1d4185b210493e311339b6c09d0153f9b8dfc6a6058

WORKDIR /root

RUN npm install "cypress-failed-log@2.9.2"
RUN npm install axe-core
RUN npm install cypress-axe

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://firm-deputy-hub:8888

COPY cypress.config.js .
COPY cypress cypress
