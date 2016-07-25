FROM alpine:3.3

RUN apk -U add ca-certificates && \
    rm -Rf /var/cache/apk/*
EXPOSE 8081
ENV PORT 8081
ENV HOST 0.0.0.0
COPY usr/bin/workflow-manager-api /usr/local/bin/workflow-manager-api
COPY static_html /static_html
CMD ["workflow-manager-api"]
