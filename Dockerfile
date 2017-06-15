FROM golang:1.8.1

COPY ./  ${GOPATH}/src/github.com/qa-dev/jsonwire-grid/

WORKDIR ${GOPATH}/src/github.com/qa-dev/jsonwire-grid/
ENV CONFIG_PATH ./config.json

RUN make

CMD ["jsonwire-grid"]

EXPOSE 4444
