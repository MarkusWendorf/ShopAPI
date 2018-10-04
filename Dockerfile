FROM alpine:edge

RUN apk add --no-cache mongodb
RUN apk add --no-cache musl-dev
RUN apk add go
RUN apk add git
ENV HOME=/
ENV GOPATH=/go
COPY . /go/src/ShopAPI
CMD /bin/sh /go/src/ShopAPI/setup.sh