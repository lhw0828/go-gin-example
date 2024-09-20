FROM scratch

WORKDIR $GOPATH/src/github.com/lhw0828/go-gin-example
COPY . $GOPATH/src/github.com/lhw0828/go-gin-example

EXPOSE 8000
ENTRYPOINT ["./go-gin-example"]