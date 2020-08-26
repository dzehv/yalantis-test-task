FROM golang:1.15.0-alpine3.12

ENV WORKDIR=/go/src/yal

WORKDIR $WORKDIR

RUN apk add git g++ gcc \
   && rm -rf /var/cache/apk/*

# import the code from the context
COPY . .

# build the executable to `./app`. Mark the build as statically linked
RUN go build -installsuffix 'static' -o exec .

# USER nobody:nobody
WORKDIR $WORKDIR

ENTRYPOINT [ "./exec" ]
EXPOSE 8080/tcp
