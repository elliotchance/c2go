BINARY=c2go

.PHONY: clean install

build:
	go build -o ${BINARY} 

install: build
	go install -x ${BINARY}.go

clean:
	if [ -f *~ ] ; then rm *~ ; fi
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	if [ -f ${GOPATH}/bin/${BINARY} ] ; then rm ${GOPATH}/bin/${BINARY} ; fi

prime:
	./c2go tests/prime.c > prime.go
	go run primie.go


