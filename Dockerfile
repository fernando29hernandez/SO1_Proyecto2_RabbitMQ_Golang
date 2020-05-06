# Start by building the application.
FROM golang:1.13-buster as build

WORKDIR /go/src/app
ADD . /go/src/app
RUN ls
RUN go get -d -v ./...

RUN go build -o /go/bin/app
RUN ls
# Now copy it into our base image.
FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/app /
ADD . /
WORKDIR /
EXPOSE 80
CMD ["/app"]
