GO=/Users/vlad/go/bin/6g
LINKER=/Users/vlad/go/bin/6l

all: ws

ws: ws.6
	$(LINKER) -o $@ $<

ws.6: statmsg.6 nullstore.6 ws.go
	$(GO) ws.go

statmsg.6: statmsg.go
	$(GO) statmsg.go

nullstore.6: statmsg.6 nullstore.go
	$(GO) nullstore.go

clean:
	rm ws *.6

.PHONY: all clean
