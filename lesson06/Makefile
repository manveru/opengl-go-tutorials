include $(GOROOT)/src/Make.$(GOARCH)

lesson=lesson06

all:
	$(GC) $(lesson).go
	$(LD) -o $(lesson) $(lesson).$(O)

TARG=$(lesson)
GOFILES=$(lesson).go
CLEANFILES+=$(lesson)

include $(GOROOT)/src/Make.pkg
