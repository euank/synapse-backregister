FROM alpine
RUN apk add --update ca-certificates
COPY ./bin/synapse-backregister /synapse-backregister

EXPOSE 8000

ENTRYPOINT ["/synapse-backregister"]
