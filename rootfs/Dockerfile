FROM alpine:3.3

RUN apk -U add ca-certificates && \
    rm -Rf /var/cache/apk/*
EXPOSE 8081
COPY . /
CMD ["workflow-manager-api"]
