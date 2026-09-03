FROM cypress/included:15.21.1@sha256:8d874a8e5457459c7b8f185f797e46c299f7d0584c9fbda0b0ca23b3e85be749

WORKDIR /root

RUN npm install "cypress-failed-log@2.9.2"
RUN npm install axe-core
RUN npm install cypress-axe

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://firm-deputy-hub:8888

COPY cypress.config.js .
COPY cypress cypress
