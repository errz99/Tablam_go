package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	//tb "Tablam_go/tablam"
	tb "github.com/errz99/tablam_go"
)

const hma string = "<span foreground=\"black\" background=\"white\" size=\"medium\"><tt><b>"
const hmb string = "</b></tt></span>"
const dma string = "<span foreground=\"blue\" background=\"white\" size=\"medium\"><tt>"
const dmb string = "</tt></span>"
const cma string = "<span foreground=\"black\" background=\"yellow\" size=\"medium\"><tt>"
const cmb string = "</tt></span>"

var shift bool = false

func main() {
	rand.Seed(time.Now().UnixNano())

	mbData := [][]string{
		{"Date", "Name", "URL", "Info"},
		{"20190904", "Vodafone", "www.vodafone.com", "Mi cuenta en la web de Vodafone"},
		{"20191001", "micuenta", "gmail.com", "Cuenta de correo en gmail"},
		{"20190522", "BNK", "www.banco.com", "Pues eso, el banco y tal"},
		{"20181105", "García", "www.zaragoza.es", "Ejemplo con tilde, y alguna cosilla"}}

	gtk.Init(nil)

	mainWin(mbData)

	gtk.Main()
}

func mainWin(mbData [][]string) {
	mwin, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	mwin.SetTitle("mBoxGo Test")
	mwin.SetDefaultSize(600, 400)

	mwinCss := "window { font-size: 15px; }"
	provider, _ := gtk.CssProviderNew()
	provider.LoadFromData(mwinCss)

	context, _ := mwin.GetStyleContext()
	context.AddProvider(provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	mwin.Connect("destroy", func() {
		gtk.MainQuit()
	})

	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	mwin.Add(vbox)

	headText, _ := gtk.LabelNew("Ejemplo")
	headText.SetMarkup("<span foreground=\"green\"><b>Ejemplo</b></span>")
	headText.SetMarginTop(8)
	vbox.Add(headText)

	aligns := []string{"rigth", "left", "center", "left"}
	tab := tb.NewTablam(mbData[0], aligns)
	tab.SetHeadAligns([]string{"left", "left", "center", "left"})

	tab.Box.SetMarginTop(8)
	tab.Box.SetBorderWidth(8)
	//tab.SetCursorMarkup(cma, cmb)
	tab.SetRowSeparation(0)
	tab.SetColumnSeparation(0)
	tab.SetLeftAndRightMargin(2)

	for i := 1; i < len(mbData); i++ {
		tab.AddRow(mbData[i])
	}

	vbox.Add(tab.Box)

	gridContext, _ := tab.Grid.GetStyleContext()
	gridContext.AddProvider(provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	close, _ := gtk.ButtonNewWithLabel("Close")
	vbox.PackEnd(close, false, false, 0)
	close.SetCanFocus(false)

	close.Connect("clicked", func() {
		gtk.MainQuit()
	})

	mwin.Connect("key-press-event", func(_ *gtk.Window, event *gdk.Event) {
		eventKey := gdk.EventKeyNewFromEvent(event)
		kval := eventKey.KeyVal()

		switch kval {
		case gdk.KEY_Shift_L:
			shift = true
			fmt.Println("shift active")

		case gdk.KEY_Up:
			tab.CursorUp()

		case gdk.KEY_Down:
			tab.CursorDown()

		case gdk.KEY_Escape:
			if tab.CursorIsActive() {
				tab.ClearCursor()
			} else {
				gtk.MainQuit()
			}

		case gdk.KEY_Return:
			if tab.ActiveData() != nil {
				fmt.Println(tab.ActiveData())
			} else {
				fmt.Println("no data active")
			}

		case gdk.KEY_Delete:
			if shift == true {
				tab.DeleteAll()
			} else {
				tab.DeleteActiveRow()
			}

		case gdk.KEY_Insert:
			tab.AddRow(modify([]string{"20190101", "Mi veloz router",
				"www.here.com", "Acceso all router de casa"}))

		case gdk.KEY_F12:
			//tab.ReverseData()

		case gdk.KEY_e:
			//if eventKey.state & ModifierType.CONTROL_MASK) {
			toEdit := tab.ActiveData()
			if toEdit != nil {
				edited := modify(toEdit)
				tab.EditActiveRow(edited)
			}
			//}

		default:
		}
	})

	mwin.Connect("key-release-event", func(_ *gtk.Window, event *gdk.Event) {
		eventKey := gdk.EventKeyNewFromEvent(event)
		kval := eventKey.KeyVal()

		switch kval {
		case gdk.KEY_Shift_L:
			shift = false
			fmt.Println("shift inactive")
		default:
		}
	})

	mwin.Connect("scroll-event", func(_ *gtk.Window, event *gdk.Event) {
		fmt.Println("scroll event")
	})

	mwin.Connect("button-press-event", func(_ *gtk.Window, event *gdk.Event) {
		fmt.Println("button press event")

		//auto eb = e.button();

		//if (e.isDoubleClick(eb)) {
		//	writeln("tab double check: get row data");
		//	if (tab.activeData() != []) {
		//		writeln(tab.activeData());
		//	} else {
		//		writeln("no data active");
		//	}
		//} else {
		//	fmt.Println("tab single check: get position")
		//}

		//return true
	})

	mwin.ShowAll()
}

func modify(str []string) []string {
	r := rand.Intn(255)
	y := rand.Intn(2020-2000) + 2000
	m := rand.Intn(13-1) + 1
	d := rand.Intn(29-1) + 1

	date := strconv.Itoa(y) + "/" + strconv.Itoa(m) + "/" + strconv.Itoa(d)
	url := "www." + strconv.Itoa(r) + ".com"
	str[0] = date
	str[2] = url

	return str
}
