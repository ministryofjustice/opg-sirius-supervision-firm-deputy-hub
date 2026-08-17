FROM cypress/included:15.19.0@sha256:418f70c94f6d727580da697faf8f4ca50c40fda78a57590708fa57772237697c

WORKDIR /root

RUN npm install "cypress-failed-log@2.9.2"
RUN npm install axe-core
RUN npm install cypress-axe

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://firm-deputy-hub:8888

COPY cypress.config.js .
COPY cypress cypress
