include $(GOROOT)/src/Make.inc

lesson=lesson10

all:
	$(GC) $(lesson).go
	$(LD) -o $(lesson) $(lesson).$(O)

TARG=$(lesson)
GOFILES=$(lesson).go
CLEANFILES+=$(lesson)

include $(GOROOT)/src/Make.pkg
