FROM golang:1.24

WORKDIR /myapp

ENV TODO_PORT=example
ENV TODO_DBFILE=example
ENV TODO_PASSWORD=example

COPY goapp ./

COPY web/ ./web/

EXPOSE 7540

CMD ["/myapp/goapp"]