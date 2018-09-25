FROM hohice/gins-builder:1.1 as builder

WORKDIR /go/src/github.com/hohice/gin-web
COPY . .

RUN make swag && make install

FROM golang:1.11-alpine3.7 
#RUN apk add --update ca-certificates && update-ca-certificates
COPY --from=builder /go/bin/* /usr/local/bin/

#CMD [ "serv","magrite","version" ] 
ENTRYPOINT [ "ginS" ]