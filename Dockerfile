FROM alpine:3
ARG PLUGIN_MODULE=github.com/traefik/hotelplanner_auth

COPY /plugins/ /plugins-local/src/${PLUGIN_MODULE}

FROM traefik:v3.0
COPY --from=0 /plugins-local /plugins-local