FROM gcr.io/google_containers/ubuntu-slim:0.2
EXPOSE 8081
ENV PORT 8081
ENV HOST 0.0.0.0
COPY usr/bin/workflow-manager-api /usr/local/bin/workflow-manager-api
CMD ["workflow-manager-api"]
