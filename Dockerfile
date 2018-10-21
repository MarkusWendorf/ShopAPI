# build container
FROM alpine:edge

RUN apk add --no-cache musl-dev
RUN apk add --no-cache go
RUN apk add --no-cache git
WORKDIR /src
COPY go.mod .
COPY go.sum .
# cache dependencies
RUN go mod download
COPY . .
RUN go build

# runtime container
FROM alpine:edge

WORKDIR /root/
COPY --from=0 /src .
ENTRYPOINT /bin/sh setup.sh