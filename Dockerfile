FROM golang

# Fetch dependencies
RUN go get github.com/tools/godep

# Add project directory to Docker image.
ADD . /go/src/github.axa.com/axa-singapore-meetups/gopee

ENV USER khanh
ENV HTTP_ADDR :8888
ENV HTTP_DRAIN_INTERVAL 1s
ENV COOKIE_SECRET 9e9DMq498fJOA2MB

# Replace this with actual PostgreSQL DSN.
# ENV DSN postgres://khanh@localhost:5432/gopee?sslmode=disable

WORKDIR /go/src/github.axa.com/axa-singapore-meetups/gopee

RUN godep go build

EXPOSE 8888
CMD ./gopee
