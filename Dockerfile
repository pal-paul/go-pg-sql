FROM golang:1.21-bullseye as build

RUN mkdir -p /pg-sql/
COPY . /pg-sql/
WORKDIR /pg-sql

ENV GO111MODULE=on
RUN make install
RUN make build

# Now copy it into our base image.
FROM gcr.io/distroless/base
COPY --from=build /pg-sql/cmd/main /cmd

ENTRYPOINT ["/cmd"]