include $(GOROOT)/src/Make.inc

TARG=walk/crutches
GOFILES=\
	crutches.go\
	crutch.go\

OFILES=\
	crutch.$O\
	queue.$O\
	wildcall.$O\

HFILES=\
	crutch.h\

GORUNTSRC=$(GOROOT)/src/pkg/runtime

## Fragment based on $(GOROOT)/src/pkg/runtime/Makefile
CFLAGS_windows=-D__WINDOWS__
CFLAGS=-I. -I$(GORUNTSRC) -I$(GORUNTSRC)/$(GOOS) -I$(GORUNTSRC)/$(GOARCH) -I$(GORUNTSRC)/$(GOOS)/$(GOARCH) -wF $(CFLAGS_$(GOARCH)) $(CFLAGS_$(GOOS))

include $(GOROOT)/src/Make.pkg
