include $(GOROOT)/src/Make.inc

TARG=ws
GOFILES=statmsg.go\
	morestore.go\
	ws.go

include $(GOROOT)/src/Make.cmd
