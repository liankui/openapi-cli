FROM scratch

ADD cmd/openapi-cli/openapi-cli openapi-cli

CMD ["./openapi-cli v"]
