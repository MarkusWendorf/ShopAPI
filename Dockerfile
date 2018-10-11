FROM alpine:edge

RUN apk add --no-cache musl-dev
RUN apk add go
RUN apk add git
ENV GOPATH=/go
COPY . /go/src/shopApi
WORKDIR "/go/src/shopApi"
RUN go get -v ./...
RUN go build


FROM alpine:edge

WORKDIR /root/
COPY --from=0 /go/src/shopApi .
ENTRYPOINT ./shopApi