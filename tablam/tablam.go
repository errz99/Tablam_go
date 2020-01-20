package mbox

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const hma string = "<span><tt><b>"
const hmb string = "</b></tt></span>"
const dma string = "<span background=\"white\"><tt>"
const dmb string = "</tt></span>"
const cma string = "<span foreground=\"white\" background=\"#6666dd\"><tt>"
const cmb string = "</tt></span>"

type Tablam struct {
	hasHead      bool
	rbs          []RowBox
	data         [][]string
	datax        [][]string
	position     int
	outPosition  int
	lastPosition int
	headMarkup   [2]string
	dataMarkup   [2]string
	cursorMarkup [2]string
	hsep         int
	max          []int
	changedMax   []int
	separation   int
	sep          string
	aligns       []string
	Grid         *gtk.Grid
}

func NewTablam(data [][]string, hasHead bool, aligns []string) Tablam {
	grid, _ := gtk.GridNew()

	var mbox = Tablam{
		hasHead,
		nil,
		nil,
		nil,
		0,
		-1,
		-1,
		[2]string{hma, hmb},
		[2]string{dma, dmb},
		[2]string{cma, cmb},
		3,
		nil,
		nil,
		1,
		" ",
		nil,
		grid,
	}

	if mbox.hasHead == true {
		mbox.outPosition++
	}

	mbox.position = mbox.outPosition
	mbox.lastPosition = mbox.outPosition
	mbox.sep = strings.Repeat(" ", mbox.separation)
	mbox.max = make([]int, len(data[0]))

	if aligns == nil {
		for range data[0] {
			mbox.aligns = append(mbox.aligns, "right")
		}
	} else {
		mbox.aligns = aligns
	}

	grid.SetHAlign(gtk.ALIGN_CENTER)
	grid.SetBorderWidth(8)
	grid.SetRowSpacing(uint(mbox.hsep))

	for _, d := range data {
		mbox.AddRow(d)
	}

	return mbox
}

func (t *Tablam) SetHeadMarkup(a, b string) {
	t.headMarkup = [2]string{a, b}
}

func (t *Tablam) SetDataMarkup(a, b string) {
	t.dataMarkup = [2]string{a, b}
}

func (t *Tablam) SetCursorMarkup(a, b string) {
	t.cursorMarkup = [2]string{a, b}
}

func (t *Tablam) SetElemAlign(i int, halign string) {
	t.aligns[i] = halign
}

func (t *Tablam) CursorDown() {
	if (t.hasHead == true && len(t.rbs) > 1) || (t.hasHead == false && len(t.rbs) > 0) {
		t.lastPosition = t.position

		if t.position < len(t.rbs)-1 {
			t.position++
		} else {
			t.position = t.outPosition + 1
		}

		t.updateCursor()
	}

	fmt.Println("down", t.position)
}

func (t *Tablam) CursorUp() {
	t.lastPosition = t.position
	t.position--

	if t.position < t.outPosition+1 {
		t.position = len(t.rbs) - 1
	}
	if t.position >= 0 {
		t.updateCursor()
	}
}

func (t *Tablam) CursorIsActive() bool {
	if t.position > t.outPosition {
		return true
	} else {
		return false
	}
}

func (t *Tablam) ClearCursor() {
	if t.position > t.outPosition {
		for i := 0; i < len(t.rbs[0].labels); i++ {
			t.rbs[t.position].labels[i].SetMarkup(
				t.dataMarkup[0] + t.rbs[t.position].datax[i] + t.dataMarkup[1])
		}
		t.position = t.outPosition
	}
}

func (t *Tablam) ActiveData() []string {
	if t.position > t.outPosition {
		return t.rbs[t.position].data
	} else {
		return nil
	}
}

func (t *Tablam) EditActiveRow(edata []string) {
	fmt.Println(edata)
	t.changedMax = []int{}
	edatax := t.newX(edata)

	t.rbs[t.position].data = edata
	t.rbs[t.position].datax = edatax
	t.data[t.position] = edata
	t.datax[t.position] = edatax

	t.updateChanged()
	t.markupActiveRow()
}

func (t *Tablam) DeleteActiveRow() {
	if t.position > t.outPosition && t.position < len(t.rbs) {

		t.rbs = append(t.rbs[:t.position], t.rbs[t.position+1:]...)
		t.data = append(t.data[:t.position], t.data[t.position+1:]...)
		t.datax = append(t.datax[:t.position], t.datax[t.position+1:]...)

		t.Grid.RemoveRow(t.position)

		for i := 0; i < len(t.rbs); i++ {
			t.rbs[i].box.SetName(strconv.Itoa(i))
		}

		if len(t.rbs) == t.outPosition+1 {
			t.position = t.outPosition
		} else if t.position == len(t.rbs) {
			t.position--
		}

		if t.position > t.outPosition {
			t.markupActiveRow()
		}
	}
}

func (t *Tablam) ReverseData() {
	reverse := func(a [][]string) [][]string {
		for i := len(a)/2 - 1; i >= 0; i-- {
			opp := len(a) - 1 - i
			a[i], a[opp] = a[opp], a[i]
		}
		return a
	}

	if t.hasHead {
		tmp := t.data[1:]
		tmpx := t.datax[1:]

		tmp = reverse(tmp)
		tmpx = reverse(tmpx)

		t.data = append(t.data[:1], tmp...)
		t.datax = append(t.datax[:1], tmpx...)

		//reverse(t.data[1:])
		//reverse(t.datax[1:])

	} else {
		t.data = reverse(t.data)
		t.datax = reverse(t.datax)
	}

	for i := t.outPosition + 1; i < len(t.datax); i++ {
		t.rbs[i].data = t.data[i]
		t.rbs[i].datax = t.datax[i]
		for j := 0; j < len(t.rbs[i].labels); j++ {
			t.applyMarkup(i, j, t.rbs[i].datax[j])
		}
	}
}

func (t *Tablam) AddRow(rdata []string) {
	rb := newRowBox(rdata, t)
	t.changedMax = rb.changedMax
	t.rbs = append(t.rbs, rb)
	t.Grid.Attach(rb.box, 0, len(t.datax), 1, 1)
	t.data = append(t.data, rb.data)
	t.datax = append(t.datax, rb.datax)
	t.updateChanged()
}

func (t *Tablam) updateCursor() {
	if t.lastPosition > t.outPosition {
		for i := 0; i < len(t.rbs[0].labels); i++ {
			t.rbs[t.lastPosition].labels[i].SetMarkup(
				t.dataMarkup[0] + t.rbs[t.lastPosition].datax[i] + t.dataMarkup[1])
		}
	}

	if t.position > t.outPosition {
		for i := 0; i < len(t.rbs[0].labels); i++ {
			t.rbs[t.position].labels[i].SetMarkup(
				t.cursorMarkup[0] + t.rbs[t.position].datax[i] + t.cursorMarkup[1])
		}
	}
}

func (t *Tablam) updateChanged() {
	for _, cm := range t.changedMax {
		for j := range t.rbs {
			grow := t.max[cm] - utf8.RuneCountInString(t.rbs[j].data[cm])
			elemx := t.createX(t.rbs[j].data[cm], cm, grow)

			t.rbs[j].datax[cm] = elemx
			t.applyMarkup(j, cm, elemx)
		}
	}
}

func (t *Tablam) markupActiveRow() {
	for i := 0; i < len(t.rbs[0].labels); i++ {
		t.rbs[t.position].labels[i].SetMarkup(
			t.cursorMarkup[0] + t.rbs[t.position].datax[i] + t.cursorMarkup[1])
	}
}

func (t *Tablam) createX(elem string, i, grow int) string {
	if t.aligns[i] == "left" {
		return t.sep + elem + strings.Repeat(" ", grow) + t.sep

	} else if t.aligns[i] == "rigth" {
		return t.sep + strings.Repeat(" ", grow) + elem + t.sep

	} else if t.aligns[i] == "center" {
		a := grow / 2
		b := grow / 2
		if grow%2 != 0 {
			b++
		}
		return t.sep + strings.Repeat(" ", a) + elem + strings.Repeat(" ", b) + t.sep
	} else {
		return t.sep + elem + strings.Repeat(" ", grow) + t.sep
	}
}

func (t *Tablam) applyMarkup(i, j int, elemx string) {
	if t.hasHead == true && i == 0 {
		t.rbs[i].labels[j].SetMarkup(t.headMarkup[0] + elemx + t.headMarkup[1])
	} else if i == t.position {
		t.rbs[i].labels[j].SetMarkup(t.cursorMarkup[0] + elemx + t.cursorMarkup[1])
	} else {
		t.rbs[i].labels[j].SetMarkup(t.dataMarkup[0] + elemx + t.dataMarkup[1])
	}
}

func (t *Tablam) newX(ndata []string) []string {
	var ndatax []string

	for i := 0; i < len(ndata); i++ {
		nrunes := utf8.RuneCountInString(ndata[i])

		if nrunes > t.max[i] {
			t.max[i] = nrunes
			t.changedMax = append(t.changedMax, i)
		}

		grow := t.max[i] - nrunes
		ndatax = append(ndatax, t.createX(ndata[i], i, grow))
	}

	return ndatax
}

type RowBox struct {
	data       []string
	datax      []string
	labels     []*gtk.Label
	changedMax []int
	box        *gtk.Box
}

func newRowBox(d []string, tab *Tablam) RowBox {
	var rb RowBox
	idx := len(tab.rbs)

	rb.box, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, tab.hsep)
	rb.box.SetName(strconv.Itoa(idx))
	rb.data = d

	//rb.datax = tab.newX(rb.data)

	for i := 0; i < len(rb.data); i++ {
		drunes := utf8.RuneCountInString(rb.data[i])

		if drunes > tab.max[i] {
			tab.max[i] = drunes
			rb.changedMax = append(rb.changedMax, i)
		}

		grow := tab.max[i] - drunes
		rb.datax = append(rb.datax, tab.createX(rb.data[i], i, grow))
	}

	for _, elemx := range rb.datax {
		ebox, _ := gtk.EventBoxNew()
		rb.box.Add(ebox)
		label, _ := gtk.LabelNew(elemx)
		label.SetMarkup(tab.dataMarkup[0] + elemx + tab.dataMarkup[1])
		ebox.Add(label)
		rb.labels = append(rb.labels, label)
	}

	rb.box.Connect("button-press-event", func(_ *gtk.Box, e *gdk.Event) bool {
		//eb := e.Button()
		name, _ := rb.box.GetName()
		namint, _ := strconv.Atoi(name)
		fmt.Println(namint)

		if namint > tab.outPosition {
			//if e.IsDoubleClick(eb) {
			//
			//} else if tab.position != namint {
			tab.lastPosition = tab.position
			tab.position = namint
			tab.updateCursor()
			//}
		} else {
			fmt.Println("button pressed at header")
		}
		return false
	})

	rb.box.ShowAll()
	return rb
}
