FROM cypress/included:12.17.4@sha256:102b34b9e4cb9895c44a74293c6931e7282535775045dd8b1c7608667a34c4b6

WORKDIR /root

RUN npm install "cypress-failed-log@2.9.2"
RUN npm install axe-core
RUN npm install cypress-axe

ENV CYPRESS_VIDEO=false
ENV CYPRESS_baseUrl=http://firm-deputy-hub:8888

COPY cypress.config.js .
COPY cypress cypress
