FROM hohice/ginS-builder:1.0 as builder

WORKDIR /go/src/github.com/hohice/gin-web
COPY . .

RUN make swag && make generate && make install

FROM golang:1.10.2-alpine3.7 
#RUN apk add --update ca-certificates && update-ca-certificates
COPY --from=builder /go/bin/* /usr/local/bin/

#CMD [ "ginS","serv" ] 
ENTRYPOINT [ "ginS","serv" ]