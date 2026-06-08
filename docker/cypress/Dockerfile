FROM cypress/included:15.16.0@sha256:057f6c6bbb661342ea712f94f5b870823823e3617354db7a78a01f2aabbfc2ad

WORKDIR /root

RUN npm install "cypress-failed-log@2.9.2"
RUN npm install axe-core
RUN npm install cypress-axe

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://firm-deputy-hub:8888

COPY cypress.config.js .
COPY cypress cypress
