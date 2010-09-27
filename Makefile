all: clean
	make -C winapi               install
	make -C winapi/kernel32      install
	make -C winapi/gdi32         install
	make -C winapi/user32        install
	make -C winapi/advapi32      install
	make -C winapi/comctl32      install
	make -C winapi/comdlg32      install
	make -C winapi/gdiplus       install
	make -C winapi/shell32       install
	make -C winapi/uxtheme       install
	make -C winapi/winspool      install
	make -C drawing              install
	make -C gui                  install
	make -C path                 install
	make -C printing             install
	make -C registry             install
	make -C examples/drawing
	make -C examples/imageviewer

test: clean
	make -C drawing              test
	make -C gui                  test
	make -C path                 test
	make -C printing             test
	make -C registry             test

clean:
	make -C winapi               clean
	make -C winapi/kernel32      clean
	make -C winapi/gdi32         clean
	make -C winapi/user32        clean
	make -C winapi/advapi32      clean
	make -C winapi/comctl32      clean
	make -C winapi/comdlg32      clean
	make -C winapi/gdiplus       clean
	make -C winapi/shell32       clean
	make -C winapi/uxtheme       clean
	make -C winapi/winspool      clean
	make -C drawing              clean
	make -C gui                  clean
	make -C path                 clean
	make -C printing             clean
	make -C registry             clean
	make -C examples/drawing     clean
	make -C examples/imageviewer clean

format:
	gofmt -w .
