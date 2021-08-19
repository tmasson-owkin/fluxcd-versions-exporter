FROM golang:1.16-alpine AS build

ENV GO111MODULE on
ENV CGO_ENABLED 0

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build

FROM alpine:3.14
WORKDIR /app
RUN apk add curl bash
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
RUN chmod +x ./kubectl
RUN mv ./kubectl /usr/local/bin/kubectl
RUN curl -s https://fluxcd.io/install.sh | bash
COPY --from=build /app/flux_version_exporter .

EXPOSE 8080

CMD [ "/app/flux_version_exporter" ]
