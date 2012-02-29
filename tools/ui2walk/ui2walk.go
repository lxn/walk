// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var forceUpdate *bool = flag.Bool("force", false, "force code generation for up-to-date files")

type UI struct {
	Class         string        `xml:"class"`
	Widget        Widget        `xml:"widget"`
	CustomWidgets CustomWidgets `xml:"customwidgets"`
	TabStops      []string      `xml:"tabstops>tabstop"`
}

type Widget struct {
	Class     string       `xml:"class,attr"`
	Name      string       `xml:"name,attr"`
	Attribute []*Attribute `xml:"attribute"`
	Property  []*Property  `xml:"property"`
	Layout    *Layout      `xml:"layout"`
	Widget    []*Widget    `xml:"widget"`
	ignored   bool
}

type Layout struct {
	Class    string      `xml:"class,attr"`
	Name     string      `xml:"name,attr"`
	Stretch  string      `xml:"stretch,attr"`
	Property []*Property `xml:"property"`
	Item     []*Item     `xml:"item"`
	ignored  bool
}

type Item struct {
	Row    string  `xml:"row,attr"`
	Column string  `xml:"column,attr"`
	Widget *Widget `xml:"widget"`
	Spacer *Spacer `xml:"spacer"`
}

type Spacer struct {
	Name     string      `xml:"name,attr"`
	Property []*Property `xml:"property"`
}

type Attribute struct {
	Name   string `xml:"name,attr"`
	String string `xml:"string"`
}

type Property struct {
	Name   string    `xml:"name,attr"`
	Bool   bool      `xml:"bool"`
	Enum   string    `xml:"enum"`
	Font   *Font     `xml:"font"`
	Number float64   `xml:"number"`
	Rect   Rectangle `xml:"rect"`
	Set    string    `xml:"set"`
	Size   Size      `xml:"size"`
	String string    `xml:"string"`
}

type Font struct {
	Family    string `xml:"family"`
	PointSize int    `xml:"pointsize"`
	Italic    bool   `xml:"italic"`
	Bold      bool   `xml:"bold"`
	Underline bool   `xml:"underline"`
	StrikeOut bool   `xml:"strikeout"`
}

type Rectangle struct {
	X      int `xml:"x"`
	Y      int `xml:"y"`
	Width  int `xml:"width"`
	Height int `xml:"height"`
}

type Size struct {
	Width  int `xml:"width"`
	Height int `xml:"height"`
}

type CustomWidgets struct {
	CustomWidget []*CustomWidget `xml:"customwidget"`
}

type CustomWidget struct {
	Class   string `xml:"class"`
	Extends string `xml:"extends"`
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func parseUI(reader io.Reader) (*UI, error) {
	ui := &UI{}

	if err := xml.NewDecoder(reader).Decode(ui); err != nil {
		return nil, err
	}

	return ui, nil
}

func writeAttribute(buf *bytes.Buffer, attr *Attribute, qualifiedReceiver string) (err error) {
	switch attr.Name {
	case "title":
		buf.WriteString(fmt.Sprintf(
			"if err := %s.SetTitle(`%s`); err != nil {\nreturn err\n}\n",
			qualifiedReceiver, attr.String))

	default:
		fmt.Printf("Ignoring unsupported attribute: '%s'\n", attr.Name)
		return nil
	}

	return nil
}

func writeAttributes(buf *bytes.Buffer, attrs []*Attribute, qualifiedReceiver string) error {
	for _, attr := range attrs {
		if err := writeAttribute(buf, attr, qualifiedReceiver); err != nil {
			return err
		}
	}

	return nil
}

func writeProperty(buf *bytes.Buffer, prop *Property, qualifiedReceiver string, widget *Widget) (err error) {
	if prop.Name == "windowTitle" && widget.Class == "QWidget" {
		return
	}

	switch prop.Name {
	case "decimals":
		buf.WriteString(fmt.Sprintf("if err := %s.SetDecimals(%d); err != nil {\nreturn err\n}\n", qualifiedReceiver, int(prop.Number)))

	case "echoMode":
		switch prop.Enum {
		case "QLineEdit::Normal":
			// nop

		case "QLineEdit::Password":
			buf.WriteString(fmt.Sprintf("%s.SetPasswordMode(true)\n", qualifiedReceiver))

		default:
			fmt.Printf("Ignoring unsupported echoMode: '%s'\n", prop.Enum)
			return nil
		}

	case "enabled":
		buf.WriteString(fmt.Sprintf("%s.SetEnabled(%t)\n", qualifiedReceiver, prop.Bool))

	case "font":
		f := prop.Font
		family := f.Family
		if family == "" {
			family = "MS Shell Dlg 2"
		}
		pointSize := f.PointSize
		if pointSize == 0 {
			pointSize = 8
		}
		buf.WriteString(fmt.Sprintf("if font, err = walk.NewFont(\"%s\", %d, ",
			family, pointSize))
		included := []bool{f.Bold, f.Italic, f.StrikeOut, f.Underline}
		flags := []string{"walk.FontBold", "walk.FontItalic", "walk.FontStrikeOut", "walk.FontUnderline"}
		var includedFlags []string
		for i := 0; i < len(included); i++ {
			if included[i] {
				includedFlags = append(includedFlags, flags[i])
			}
		}
		if len(includedFlags) == 0 {
			buf.WriteString("0")
		} else {
			buf.WriteString(strings.Join(includedFlags, "|"))
		}
		buf.WriteString(`); err != nil {
			return err
			}
			`)
		buf.WriteString(fmt.Sprintf("%s.SetFont(font)\n", qualifiedReceiver))

	case "geometry":
		if qualifiedReceiver == "w" {
			// Only set client size for top level
			buf.WriteString(fmt.Sprintf(
				`if err := %s.SetClientSize(walk.Size{%d, %d}); err != nil {
			return err
			}
			`,
				qualifiedReceiver, prop.Rect.Width, prop.Rect.Height))
		} else {
			buf.WriteString(fmt.Sprintf(
				`if err := %s.SetBounds(walk.Rectangle{%d, %d, %d, %d}); err != nil {
			return err
			}
			`,
				qualifiedReceiver, prop.Rect.X, prop.Rect.Y, prop.Rect.Width, prop.Rect.Height))
		}

	case "maximumSize", "minimumSize":
		// We do these two guys in writeProperties, because we want to map them
		// to a single method call, if both are present.

	case "maxLength":
		buf.WriteString(fmt.Sprintf("%s.SetMaxLength(%d)\n", qualifiedReceiver, int(prop.Number)))

	case "readOnly":
		buf.WriteString(fmt.Sprintf("%s.SetReadOnly(%t)\n", qualifiedReceiver, prop.Bool))

	case "text":
		buf.WriteString(fmt.Sprintf(
			"if err := %s.SetText(`%s`); err != nil {\nreturn err\n}\n",
			qualifiedReceiver, prop.String))

	case "title", "windowTitle":
		buf.WriteString(fmt.Sprintf(
			"if err := %s.SetTitle(`%s`); err != nil {\nreturn err\n}\n",
			qualifiedReceiver, prop.String))

	case "orientation":
		var orientation string
		switch prop.Enum {
		case "Qt::Horizontal":
			orientation = "walk.Horizontal"

		case "Qt::Vertical":
			orientation = "walk.Vertical"

		default:
			return errors.New(fmt.Sprintf("unknown orientation: '%s'", prop.Enum))
		}

		buf.WriteString(fmt.Sprintf(
			`if err := %s.SetOrientation(%s); err != nil {
			return err
			}
			`,
			qualifiedReceiver, orientation))

	default:
		fmt.Printf("Ignoring unsupported property: '%s'\n", prop.Name)
		return nil
	}

	return
}

func writeProperties(buf *bytes.Buffer, props []*Property, qualifiedReceiver string, widget *Widget) error {
	var minSize, maxSize Size
	var hasMinOrMaxSize bool

	for _, prop := range props {
		if err := writeProperty(buf, prop, qualifiedReceiver, widget); err != nil {
			return err
		}

		if prop.Name == "minimumSize" {
			minSize = prop.Size
			hasMinOrMaxSize = true
		}
		if prop.Name == "maximumSize" {
			maxSize = prop.Size
			hasMinOrMaxSize = true
		}
	}

	if hasMinOrMaxSize {
		buf.WriteString(fmt.Sprintf(
			`if err := %s.SetMinMaxSize(walk.Size{%d, %d}, walk.Size{%d, %d}); err != nil {
			return err
			}
			`,
			qualifiedReceiver, minSize.Width, minSize.Height, maxSize.Width, maxSize.Height))
	}

	return nil
}

func writeItemInitializations(buf *bytes.Buffer, items []*Item, parent *Widget, qualifiedParent string, layout string) error {
	for _, item := range items {
		var itemName string

		if item.Spacer != nil {
			itemName = item.Spacer.Name
			name2Prop := make(map[string]*Property)

			for _, prop := range item.Spacer.Property {
				name2Prop[prop.Name] = prop
			}

			orientation := name2Prop["orientation"]
			sizeType := name2Prop["sizeType"]
			sizeHint := name2Prop["sizeHint"]

			var orientStr string
			var fixedStr string
			var secondParamStr string

			if orientation.Enum == "Qt::Horizontal" {
				orientStr = "H"

				if sizeType != nil && sizeType.Enum == "QSizePolicy::Fixed" {
					fixedStr = "Fixed"
					secondParamStr = fmt.Sprintf(", %d", sizeHint.Size.Width)
				}
			} else {
				orientStr = "V"

				if sizeType != nil && sizeType.Enum == "QSizePolicy::Fixed" {
					fixedStr = "Fixed"
					secondParamStr = fmt.Sprintf(", %d", sizeHint.Size.Height)
				}
			}

			if layout == "" {
				buf.WriteString(fmt.Sprintf(
					`
					// anonymous spacer
					if _, err := walk.New%sSpacer%s(%s%s); err != nil {
					return err
					}
					`,
					orientStr, fixedStr, qualifiedParent, secondParamStr))
			} else {
				buf.WriteString(fmt.Sprintf(
					`
					// %s
					%s, err := walk.New%sSpacer%s(%s%s)
					if err != nil {
					return err
					}
					`,
					itemName, itemName, orientStr, fixedStr, qualifiedParent, secondParamStr))
			}
		}

		if item.Widget != nil && !item.Widget.ignored {
			itemName = fmt.Sprintf("w.ui.%s", item.Widget.Name)
			if err := writeWidgetInitialization(buf, item.Widget, parent, qualifiedParent); err != nil {
				return err
			}
		}

		if layout != "" && itemName != "" && item.Row != "" && item.Column != "" {
			buf.WriteString(fmt.Sprintf(
				`				if err := %s.SetLocation(%s, %s, %s); err != nil {
				return err
				}
				`,
				layout, itemName, item.Row, item.Column))
		}
	}

	return nil
}

func writeLayoutInitialization(buf *bytes.Buffer, layout *Layout, parent *Widget, qualifiedParent string) error {
	var typ string
	switch layout.Class {
	case "QGridLayout":
		typ = "GridLayout"

	case "QHBoxLayout":
		typ = "HBoxLayout"

	case "QVBoxLayout":
		typ = "VBoxLayout"

	default:
		return errors.New(fmt.Sprintf("unsupported layout type: '%s'", layout.Class))
	}

	buf.WriteString(fmt.Sprintf("%s := walk.New%s()\n",
		layout.Name, typ))

	buf.WriteString(fmt.Sprintf(
		`if err := %s.SetLayout(%s); err != nil {
		return err
		}
		`,
		qualifiedParent, layout.Name))

	spacing := 6
	margL, margT, margR, margB := 9, 9, 9, 9

	for _, prop := range layout.Property {
		switch prop.Name {
		case "spacing":
			spacing = int(prop.Number)

		case "leftMargin":
			margL = int(prop.Number)

		case "topMargin":
			margT = int(prop.Number)

		case "rightMargin":
			margR = int(prop.Number)

		case "bottomMargin":
			margB = int(prop.Number)

		case "margin":
			m := int(prop.Number)
			margL, margT, margR, margB = m, m, m, m
		}
	}

	if margL != 0 || margT != 0 || margR != 0 || margB != 0 {
		buf.WriteString(fmt.Sprintf(
			`if err := %s.SetMargins(walk.Margins{%d, %d, %d, %d}); err != nil {
			return err
			}
			`,
			layout.Name, margL, margT, margR, margB))
	}

	if spacing != 0 {
		buf.WriteString(fmt.Sprintf(
			`if err := %s.SetSpacing(%d); err != nil {
			return err
			}
			`,
			layout.Name, spacing))
	}

	var layoutName string
	if typ == "GridLayout" {
		layoutName = layout.Name
	}

	if err := writeItemInitializations(buf, layout.Item, parent, qualifiedParent, layoutName); err != nil {
		return err
	}

	return nil
}

func writeWidgetInitialization(buf *bytes.Buffer, widget *Widget, parent *Widget, qualifiedParent string) error {
	receiver := fmt.Sprintf("w.ui.%s", widget.Name)

	var typ string
	var custom bool
	switch widget.Class {
	case "QCheckBox":
		typ = "CheckBox"

	case "QComboBox":
		typ = "ComboBox"

	case "QDateEdit":
		typ = "DateEdit"

	case "QDoubleSpinBox", "QSpinBox":
		typ = "NumberEdit"

	case "QFrame":
		typ = "Composite"

	case "QGroupBox":
		typ = "GroupBox"

	case "QLabel":
		typ = "Label"

	case "QLineEdit":
		typ = "LineEdit"

	case "QPlainTextEdit", "QTextEdit":
		typ = "TextEdit"

	case "QProgressBar":
		typ = "ProgressBar"

	case "QPushButton", "QToolButton":
		typ = "PushButton"

	case "QRadioButton":
		typ = "RadioButton"

	case "QSplitter":
		typ = "Splitter"

	case "QTabWidget":
		typ = "TabWidget"

	case "QTableView", "QTableWidget":
		typ = "TableView"

	case "QTreeView", "QTreeWidget":
		typ = "TreeView"

	case "QWebView":
		typ = "WebView"

	case "QWidget":
		if parent != nil && parent.Class == "QTabWidget" {
			typ = "TabPage"
		} else {
			typ = "Composite"
		}

	default:
		// FIXME: We assume this is a custom widget in the same package.
		// We also require a func NewFoo(parent) (*Foo, error).
		typ = widget.Class
		custom = true
	}

	if custom {
		buf.WriteString(fmt.Sprintf(
			`
			// %s
			if %s, err = New%s(%s); err != nil {
			return err
			}
			`,
			widget.Name, receiver, typ, qualifiedParent))
	} else {
		if typ == "TabPage" {
			buf.WriteString(fmt.Sprintf(
				`
				// %s
				if %s, err = walk.NewTabPage(); err != nil {
				return err
				}
				`,
				widget.Name, receiver))
		} else {
			buf.WriteString(fmt.Sprintf(
				`
				// %s
				if %s, err = walk.New%s(%s); err != nil {
				return err
				}
				`,
				widget.Name, receiver, typ, qualifiedParent))
		}
	}

	buf.WriteString(fmt.Sprintf("%s.SetName(\"%s\")\n",
		receiver, widget.Name))

	if err := writeAttributes(buf, widget.Attribute, receiver); err != nil {
		return err
	}

	if err := writeProperties(buf, widget.Property, receiver, widget); err != nil {
		return err
	}

	if widget.Layout != nil && !widget.Layout.ignored {
		if err := writeLayoutInitialization(buf, widget.Layout, widget, receiver); err != nil {
			return err
		}
	}

	if typ == "TabPage" {
		buf.WriteString(fmt.Sprintf(
			`if err := %s.Pages().Add(%s); err != nil {
			return err
			}
			`,
			qualifiedParent, receiver))
	}

	return writeWidgetInitializations(buf, widget.Widget, widget, receiver)
}

func writeWidgetInitializations(buf *bytes.Buffer, widgets []*Widget, parent *Widget, qualifiedParent string) error {
	for _, widget := range widgets {
		if widget.ignored {
			continue
		}

		if err := writeWidgetInitialization(buf, widget, parent, qualifiedParent); err != nil {
			return err
		}
	}

	return nil
}

func writeWidgetDecl(buf *bytes.Buffer, widget *Widget, parent *Widget) error {
	var typ string
	switch widget.Class {
	case "QCheckBox":
		typ = "walk.CheckBox"

	case "QComboBox":
		typ = "walk.ComboBox"

	case "QDateEdit":
		typ = "walk.DateEdit"

	case "QDoubleSpinBox", "QSpinBox":
		typ = "walk.NumberEdit"

	case "QFrame":
		typ = "walk.Composite"

	case "QGroupBox":
		typ = "walk.GroupBox"

	case "QLabel":
		typ = "walk.Label"

	case "QLineEdit":
		typ = "walk.LineEdit"

	case "QPlainTextEdit", "QTextEdit":
		typ = "walk.TextEdit"

	case "QProgressBar":
		typ = "walk.ProgressBar"

	case "QPushButton", "QToolButton":
		typ = "walk.PushButton"

	case "QRadioButton":
		typ = "walk.RadioButton"

	case "QSplitter":
		typ = "walk.Splitter"

	case "QTabWidget":
		typ = "walk.TabWidget"

	case "QTableView", "QTableWidget":
		typ = "walk.TableView"

	case "QTreeView", "QTreeWidget":
		typ = "walk.TreeView"

	case "QWebView":
		typ = "walk.WebView"

	case "QWidget":
		if parent != nil && parent.Class == "QTabWidget" {
			typ = "walk.TabPage"
		} else {
			typ = "walk.Composite"
		}

	default:
		// FIXME: For now, we assume this is a custom widget in the same package
		typ = widget.Class
	}

	buf.WriteString(fmt.Sprintf("%s *%s\n", widget.Name, typ))

	if widget.Layout != nil {
		return writeItemDecls(buf, widget.Layout.Item, widget)
	}

	return writeWidgetDecls(buf, widget.Widget, widget)
}

func writeWidgetDecls(buf *bytes.Buffer, widgets []*Widget, parent *Widget) error {
	for _, widget := range widgets {
		if err := writeWidgetDecl(buf, widget, parent); err != nil {
			return err
		}
	}

	return nil
}

func writeItemDecls(buf *bytes.Buffer, items []*Item, parent *Widget) error {
	for _, item := range items {
		if item.Widget == nil {
			continue
		}

		if err := writeWidgetDecl(buf, item.Widget, parent); err != nil {
			return err
		}
	}

	return nil
}

func generateCode(buf *bytes.Buffer, ui *UI) error {
	// Comment, package decl, imports
	buf.WriteString(
		`// THIS FILE WAS GENERATED BY A TOOL, DO NOT EDIT!

		package main
		
		import (
			"github.com/lxn/walk"
		)
		
		`)

	// Embed the corresponding Walk type.
	var embeddedType string
	switch ui.Widget.Class {
	case "QMainWindow":
		embeddedType = "MainWindow"

	case "QDialog":
		embeddedType = "Dialog"

	case "QWidget":
		embeddedType = "Composite"

	default:
		return errors.New(fmt.Sprintf("Top level '%s' currently not supported.", ui.Widget.Class))
	}

	// Struct containing all descendant widgets.
	buf.WriteString(fmt.Sprintf("type %s%sUI struct {\n", strings.ToLower(ui.Class[:1]), ui.Class[1:]))

	// Descendant widget decls
	if ui.Widget.Widget != nil {
		if err := writeWidgetDecls(buf, ui.Widget.Widget, &ui.Widget); err != nil {
			return err
		}
	}

	if ui.Widget.Layout != nil {
		if err := writeItemDecls(buf, ui.Widget.Layout.Item, &ui.Widget); err != nil {
			return err
		}
	}

	// end struct
	buf.WriteString("}\n\n")

	// init func
	var qualifiedParent string
	switch embeddedType {
	case "MainWindow":
		buf.WriteString(fmt.Sprintf(
			`func (w *%s) init() (err error) {
			if w.MainWindow, err = walk.NewMainWindow()`,
			ui.Widget.Name))
		qualifiedParent = "w.ClientArea()"

	case "Dialog":
		buf.WriteString(fmt.Sprintf(
			`func (w *%s) init(owner walk.RootWidget) (err error) {
			if w.Dialog, err = walk.NewDialog(owner)`,
			ui.Widget.Name))
		qualifiedParent = "w"

	case "Composite":
		buf.WriteString(fmt.Sprintf(
			`func (w *%s) init(parent walk.Container) (err error) {
			if w.Composite, err = walk.NewComposite(parent)`,
			ui.Widget.Name))
		qualifiedParent = "w"
	}

	buf.WriteString(fmt.Sprintf(`; err != nil {
			return err
			}
			
			succeeded := false
			defer func(){
				if !succeeded {
					w.Dispose()
				}
			}()
			
			var font *walk.Font
			if font == nil {
				font = nil
			}
			
			w.SetName("%s")
			`,
		ui.Widget.Name))

	if err := writeProperties(buf, ui.Widget.Property, "w", &ui.Widget); err != nil {
		return err
	}

	if ui.Widget.Widget != nil {
		if err := writeWidgetInitializations(buf, ui.Widget.Widget, &ui.Widget, qualifiedParent); err != nil {
			return err
		}
	}

	if ui.Widget.Layout != nil {
		if err := writeLayoutInitialization(buf, ui.Widget.Layout, &ui.Widget, qualifiedParent); err != nil {
			return err
		}
	}

	buf.WriteString("\n// Tab order\n")
	for i := len(ui.TabStops) - 1; i >= 0; i-- {
		buf.WriteString(fmt.Sprintf(`if err = w.ui.%s.BringToTop(); err != nil {
			return err
		}
		`,
			ui.TabStops[i]))
	}

	// end func
	buf.WriteString(`
		succeeded = true
		
		return nil
		}`)

	return nil
}

func processFile(uiFilePath string) error {
	goFilePath := uiFilePath[:len(uiFilePath)-3] + "_ui.go"

	uiFileInfo, err := os.Stat(uiFilePath)
	if err != nil {
		return err
	}

	goFileInfo, err := os.Stat(goFilePath)
	if !*forceUpdate && err == nil && !uiFileInfo.ModTime().After(goFileInfo.ModTime()) {
		// The go file should be up-to-date
		return nil
	}

	fmt.Printf("Processing '%s'\n", uiFilePath)
	defer fmt.Println("")

	uiFile, err := os.Open(uiFilePath)
	if err != nil {
		return err
	}
	defer uiFile.Close()

	reader := bufio.NewReader(uiFile)

	ui, err := parseUI(reader)
	if err != nil {
		return err
	}

	goFile, err := os.Create(goFilePath)
	if err != nil {
		return err
	}
	defer goFile.Close()

	buf := bytes.NewBuffer(nil)

	if err := generateCode(buf, ui); err != nil {
		return err
	}

	if _, err := io.Copy(goFile, buf); err != nil {
		return err
	}
	if err := goFile.Close(); err != nil {
		return err
	}

	dirPath := os.Getenv("GOBIN")
	if dirPath == "" {
		dirPath = filepath.Join(runtime.GOROOT(), "bin")
	}

	gofmtPath := filepath.Join(dirPath, "gofmt.exe")

	gofmt, err := os.StartProcess(gofmtPath, []string{gofmtPath, "-w", goFilePath}, &os.ProcAttr{Files: []*os.File{nil, nil, os.Stderr}})
	if err != nil {
		return err
	}
	defer gofmt.Release()

	return nil
}

func processDirectory(dirPath string) error {
	dir, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		fullPath := path.Join(dirPath, name)

		fi, err := os.Stat(fullPath)
		if err != nil {
			return err
		}

		if fi.IsDir() {
			if err := processDirectory(fullPath); err != nil {
				return err
			}
		} else if !fi.IsDir() && strings.HasSuffix(name, ".ui") {
			if err := processFile(fullPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	flag.Parse()

	cwd, err := os.Getwd()
	logFatal(err)

	logFatal(processDirectory(cwd))
}
