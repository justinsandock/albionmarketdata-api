FROM golang

WORKDIR /go/src/app

# only listen on all ipv4 addresses, the default (ipv4+6) has issues on Docker
ENV ADA_LISTEN=0.0.0.0:3080
EXPOSE 3080

COPY . .

RUN go-wrapper download 
RUN go-wrapper install

CMD ["go-wrapper", "run"]
