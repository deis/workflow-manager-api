FROM alpine:3.3

RUN apk -U add ca-certificates && \
    rm -Rf /var/cache/apk/*
EXPOSE 8081
COPY ./usr/bin/workflow-manager-api /usr/bin/workflow-manager-api
CMD ["workflow-manager-api"]
