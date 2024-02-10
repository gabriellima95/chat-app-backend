FROM golang:1.19

ARG AIR_VERSION=1.27.3

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . /app

# RUN go build -o .
# RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin 
# RUN cp ./bin/air /bin/air
# RUN curl -s -L https://github.com/cosmtrek/air/releases/download/v${AIR_VERSION}/air_${AIR_VERSION}_linux_amd64 > air_${AIR_VERSION}_linux_amd64 && \
#     sha256sum -c air_checksums.txt --ignore-missing --strict && \
#     mv air_${AIR_VERSION}_linux_amd64 air && \
#     chmod +x /go/bin/air && rm -rf *.txt

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air
# RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh

EXPOSE 8080

CMD ["air", "serve"]