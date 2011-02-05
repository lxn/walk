all: clean
	make -C pkg/walk/winapi          install
	make -C pkg/walk/winapi/kernel32 install
	make -C pkg/walk/winapi/gdi32    install
	make -C pkg/walk/winapi/user32   install
	make -C pkg/walk/winapi/advapi32 install
	make -C pkg/walk/winapi/comctl32 install
	make -C pkg/walk/winapi/ole32    install
	make -C pkg/walk/winapi/oleaut32 install
	make -C pkg/walk/winapi/shdocvw  install
	make -C pkg/walk/winapi/comdlg32 install
	make -C pkg/walk/winapi/gdiplus  install
	make -C pkg/walk/winapi/shell32  install
	make -C pkg/walk/winapi/uxtheme  install
	make -C pkg/walk/winapi/winspool install
	make -C pkg/walk                 install
	make -C cmd/examples/drawing
	make -C cmd/examples/filebrowser
	make -C cmd/examples/imageviewer
	make -C cmd/examples/printing
	make -C cmd/examples/webbrowser

test: clean
	make -C pkg/walk                 test

clean:
	make -C pkg/walk/winapi          clean
	make -C pkg/walk/winapi/kernel32 clean
	make -C pkg/walk/winapi/gdi32    clean
	make -C pkg/walk/winapi/user32   clean
	make -C pkg/walk/winapi/advapi32 clean
	make -C pkg/walk/winapi/comctl32 clean
	make -C pkg/walk/winapi/ole32    clean
	make -C pkg/walk/winapi/oleaut32 clean
	make -C pkg/walk/winapi/shdocvw  clean
	make -C pkg/walk/winapi/comdlg32 clean
	make -C pkg/walk/winapi/gdiplus  clean
	make -C pkg/walk/winapi/shell32  clean
	make -C pkg/walk/winapi/uxtheme  clean
	make -C pkg/walk/winapi/winspool clean
	make -C pkg/walk                 clean
	make -C cmd/examples/drawing     clean
	make -C cmd/examples/filebrowser clean
	make -C cmd/examples/imageviewer clean
	make -C cmd/examples/printing    clean
	make -C cmd/examples/webbrowser  clean

format:
	gofmt -w .
