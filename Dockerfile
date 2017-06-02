FROM golang:1.8.1

COPY ./  ${GOPATH}/src/jsonwire-grid

WORKDIR ${GOPATH}/src/jsonwire-grid
ENV CONFIG_PATH ./config.json

RUN make

CMD ["service-entrypoint"]

EXPOSE 4444
