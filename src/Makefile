all: clean
	make -C pkg/walk/winapi          install
	make -C pkg/walk                 install
	make -C cmd/tools/ui2walk        install
	make -C cmd/examples/drawing
	make -C cmd/examples/filebrowser
	make -C cmd/examples/imageviewer
	make -C cmd/examples/notifyicon
	make -C cmd/examples/printing
	make -C cmd/examples/tableview
	make -C cmd/examples/webbrowser

test: clean
	make -C pkg/walk                 test

clean:
	make -C pkg/walk/winapi          clean
	make -C pkg/walk                 clean
	make -C cmd/tools/ui2walk        clean
	make -C cmd/examples/drawing     clean
	make -C cmd/examples/filebrowser clean
	make -C cmd/examples/imageviewer clean
	make -C cmd/examples/notifyicon  clean
	make -C cmd/examples/printing    clean
	make -C cmd/examples/tableview   clean
	make -C cmd/examples/webbrowser  clean

format:
	gofmt -w .
