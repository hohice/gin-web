FROM 172.16.1.99/transwarp/walm-builder:1.0 as builder

WORKDIR /go/src/transwarp/github.com/hohice/gin-web
COPY . .

RUN make swag && make generate && make install

FROM 172.16.1.99/gold/helm:tos18-latest
#RUN apk add --update ca-certificates && update-ca-certificates
COPY --from=builder /go/bin/* /usr/local/bin/

CMD [ "ginS","serv" ] 
#ENTRYPOINT [ "ginS","serv" ]