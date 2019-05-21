#############################################
# GET/CACHE NPM DEPS
#############################################
FROM node:alpine as npm-dependencies
WORKDIR /app
# get npm modules (cache)
COPY ./react/package.json /app/react/
COPY ./react/package-lock.json /app/react/
RUN set -x \
    && cd /app/react \
    && npm install

#############################################
# BUILD REACT APP
#############################################
FROM node:alpine as frontend
# Copy app and build
COPY ./ /app
COPY --from=npm-dependencies /app/react/node_modules/ /app/react/node_modules/
RUN set -x \
    && cd /app/react/ \
    && npm run build \
    && cp /app/react/build/index.html /app/app/views/react.html \
    && rm -rf /app/public/js \
    && mkdir -p /app/public/js \
    && cp /app/react/build/static/js/* /app/public/js \
    && rm -rf /app/react

#############################################
# BUILD GO APP
#############################################
FROM golang:1.12 as backend
WORKDIR /go/src/app
COPY ./ /go/src/app
RUN curl https://glide.sh/get | sh && glide install
RUN go get github.com/revel/cmd/revel
RUN go get -u github.com/golang/dep/cmd/dep
COPY --from=frontend /app /go/src/app
RUN set -x \
    && revel build app /app prod \
    && cp /go/src/app/docker/entrypoint.sh /app/entrypoint.sh \
    && chmod +x /app/entrypoint.sh \
    && rm -rf /app/src/app/tests

#############################################
# FINAL IMAGE
#############################################
FROM debian:stable-slim

RUN apt-get update \
    && apt-get install -y \
        ca-certificates \
    	wget \
    	curl \
    && rm -rf /var/lib/apt/lists/*

COPY --from=backend /app/ /app/
USER 65534
EXPOSE 9000
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/run.sh"]
