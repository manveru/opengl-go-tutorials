include $(GOROOT)/src/Make.$(GOARCH)

lesson=lesson07

all:
	$(GC) $(lesson).go
	$(LD) -o $(lesson) $(lesson).$(O)

TARG=$(lesson)
GOFILES=$(lesson).go
CLEANFILES+=$(lesson)

include $(GOROOT)/src/Make.pkg
