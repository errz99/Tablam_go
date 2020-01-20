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

type MBox struct {
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

func NewMBox(data [][]string, hasHead bool, aligns []string) MBox {
	grid, _ := gtk.GridNew()

	var mbox = MBox{
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

func (mb *MBox) SetHeadMarkup(a, b string) {
	mb.headMarkup = [2]string{a, b}
}

func (mb *MBox) SetDataMarkup(a, b string) {
	mb.dataMarkup = [2]string{a, b}
}

func (mb *MBox) SetCursorMarkup(a, b string) {
	mb.cursorMarkup = [2]string{a, b}
}

func (mb *MBox) SetElemAlign(i int, halign string) {
	mb.aligns[i] = halign
}

func (mb *MBox) CursorDown() {
	if (mb.hasHead == true && len(mb.rbs) > 1) || (mb.hasHead == false && len(mb.rbs) > 0) {
		mb.lastPosition = mb.position

		if mb.position < len(mb.rbs)-1 {
			mb.position++
		} else {
			mb.position = mb.outPosition + 1
		}

		mb.updateCursor()
	}

	fmt.Println("down", mb.position)
}

func (mb *MBox) CursorUp() {
	mb.lastPosition = mb.position
	mb.position--

	if mb.position < mb.outPosition+1 {
		mb.position = len(mb.rbs) - 1
	}
	if mb.position >= 0 {
		mb.updateCursor()
	}
}

func (mb *MBox) CursorIsActive() bool {
	if mb.position > mb.outPosition {
		return true
	} else {
		return false
	}
}

func (mb *MBox) ClearCursor() {
	if mb.position > mb.outPosition {
		for i := 0; i < len(mb.rbs[0].labels); i++ {
			mb.rbs[mb.position].labels[i].SetMarkup(
				mb.dataMarkup[0] + mb.rbs[mb.position].datax[i] + mb.dataMarkup[1])
		}
		mb.position = mb.outPosition
	}
}

func (mb *MBox) ActiveData() []string {
	if mb.position > mb.outPosition {
		return mb.rbs[mb.position].data
	} else {
		return nil
	}
}

func (mb *MBox) EditActiveRow(edata []string) {
	fmt.Println(edata)
	mb.changedMax = []int{}
	edatax := mb.newX(edata)

	mb.rbs[mb.position].data = edata
	mb.rbs[mb.position].datax = edatax
	mb.data[mb.position] = edata
	mb.datax[mb.position] = edatax

	mb.updateChanged()
	mb.markupActiveRow()
}

func (mb *MBox) DeleteActiveRow() {
	if mb.position > mb.outPosition && mb.position < len(mb.rbs) {

		mb.rbs = append(mb.rbs[:mb.position], mb.rbs[mb.position+1:]...)
		mb.data = append(mb.data[:mb.position], mb.data[mb.position+1:]...)
		mb.datax = append(mb.datax[:mb.position], mb.datax[mb.position+1:]...)

		mb.Grid.RemoveRow(mb.position)

		for i := 0; i < len(mb.rbs); i++ {
			mb.rbs[i].box.SetName(strconv.Itoa(i))
		}

		if len(mb.rbs) == mb.outPosition+1 {
			mb.position = mb.outPosition
		} else if mb.position == len(mb.rbs) {
			mb.position--
		}

		if mb.position > mb.outPosition {
			mb.markupActiveRow()
		}
	}
}

func (mb *MBox) ReverseData() {
	reverse := func(a [][]string) [][]string {
		for i := len(a)/2 - 1; i >= 0; i-- {
			opp := len(a) - 1 - i
			a[i], a[opp] = a[opp], a[i]
		}
		return a
	}

	if mb.hasHead {
		tmp := mb.data[1:]
		tmpx := mb.datax[1:]

		tmp = reverse(tmp)
		tmpx = reverse(tmpx)

		mb.data = append(mb.data[:1], tmp...)
		mb.datax = append(mb.datax[:1], tmpx...)

		//reverse(mb.data[1:])
		//reverse(mb.datax[1:])

	} else {
		mb.data = reverse(mb.data)
		mb.datax = reverse(mb.datax)
	}

	for i := mb.outPosition + 1; i < len(mb.datax); i++ {
		mb.rbs[i].data = mb.data[i]
		mb.rbs[i].datax = mb.datax[i]
		for j := 0; j < len(mb.rbs[i].labels); j++ {
			mb.applyMarkup(i, j, mb.rbs[i].datax[j])
		}
	}
}

func (mb *MBox) AddRow(rdata []string) {
	rb := newRowBox(rdata, mb)
	mb.changedMax = rb.changedMax
	mb.rbs = append(mb.rbs, rb)
	mb.Grid.Attach(rb.box, 0, len(mb.datax), 1, 1)
	mb.data = append(mb.data, rb.data)
	mb.datax = append(mb.datax, rb.datax)
	mb.updateChanged()
}

func (mb *MBox) updateCursor() {
	if mb.lastPosition > mb.outPosition {
		for i := 0; i < len(mb.rbs[0].labels); i++ {
			mb.rbs[mb.lastPosition].labels[i].SetMarkup(
				mb.dataMarkup[0] + mb.rbs[mb.lastPosition].datax[i] + mb.dataMarkup[1])
		}
	}

	if mb.position > mb.outPosition {
		for i := 0; i < len(mb.rbs[0].labels); i++ {
			mb.rbs[mb.position].labels[i].SetMarkup(
				mb.cursorMarkup[0] + mb.rbs[mb.position].datax[i] + mb.cursorMarkup[1])
		}
	}
}

func (mb *MBox) updateChanged() {
	for _, cm := range mb.changedMax {
		for j := range mb.rbs {
			grow := mb.max[cm] - utf8.RuneCountInString(mb.rbs[j].data[cm])
			elemx := mb.createX(mb.rbs[j].data[cm], cm, grow)

			mb.rbs[j].datax[cm] = elemx
			mb.applyMarkup(j, cm, elemx)
		}
	}
}

func (mb *MBox) markupActiveRow() {
	for i := 0; i < len(mb.rbs[0].labels); i++ {
		mb.rbs[mb.position].labels[i].SetMarkup(
			mb.cursorMarkup[0] + mb.rbs[mb.position].datax[i] + mb.cursorMarkup[1])
	}
}

func (mb *MBox) createX(elem string, i, grow int) string {
	if mb.aligns[i] == "left" {
		return mb.sep + elem + strings.Repeat(" ", grow) + mb.sep

	} else if mb.aligns[i] == "rigth" {
		return mb.sep + strings.Repeat(" ", grow) + elem + mb.sep

	} else if mb.aligns[i] == "center" {
		a := grow / 2
		b := grow / 2
		if grow%2 != 0 {
			b++
		}
		return mb.sep + strings.Repeat(" ", a) + elem + strings.Repeat(" ", b) + mb.sep
	} else {
		return mb.sep + elem + strings.Repeat(" ", grow) + mb.sep
	}
}

func (mb *MBox) applyMarkup(i, j int, elemx string) {
	if mb.hasHead == true && i == 0 {
		mb.rbs[i].labels[j].SetMarkup(mb.headMarkup[0] + elemx + mb.headMarkup[1])
	} else if i == mb.position {
		mb.rbs[i].labels[j].SetMarkup(mb.cursorMarkup[0] + elemx + mb.cursorMarkup[1])
	} else {
		mb.rbs[i].labels[j].SetMarkup(mb.dataMarkup[0] + elemx + mb.dataMarkup[1])
	}
}

func (mb *MBox) newX(ndata []string) []string {
	var ndatax []string

	for i := 0; i < len(ndata); i++ {
		nrunes := utf8.RuneCountInString(ndata[i])

		if nrunes > mb.max[i] {
			mb.max[i] = nrunes
			mb.changedMax = append(mb.changedMax, i)
		}

		grow := mb.max[i] - nrunes
		ndatax = append(ndatax, mb.createX(ndata[i], i, grow))
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

func newRowBox(d []string, mb *MBox) RowBox {
	var rb RowBox
	idx := len(mb.rbs)

	rb.box, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, mb.hsep)
	rb.box.SetName(strconv.Itoa(idx))
	rb.data = d

	//rb.datax = mb.newX(rb.data)

	for i := 0; i < len(rb.data); i++ {
		drunes := utf8.RuneCountInString(rb.data[i])

		if drunes > mb.max[i] {
			mb.max[i] = drunes
			rb.changedMax = append(rb.changedMax, i)
		}

		grow := mb.max[i] - drunes
		rb.datax = append(rb.datax, mb.createX(rb.data[i], i, grow))
	}

	for _, elemx := range rb.datax {
		ebox, _ := gtk.EventBoxNew()
		rb.box.Add(ebox)
		label, _ := gtk.LabelNew(elemx)
		label.SetMarkup(mb.dataMarkup[0] + elemx + mb.dataMarkup[1])
		ebox.Add(label)
		rb.labels = append(rb.labels, label)
	}

	rb.box.Connect("button-press-event", func(_ *gtk.Box, e *gdk.Event) bool {
		//eb := e.Button()
		name, _ := rb.box.GetName()
		namint, _ := strconv.Atoi(name)
		fmt.Println(namint)

		if namint > mb.outPosition {
			//if e.IsDoubleClick(eb) {
			//
			//} else if mb.position != namint {
			mb.lastPosition = mb.position
			mb.position = namint
			mb.updateCursor()
			//}
		} else {
			fmt.Println("button pressed at header")
		}
		return false
	})

	rb.box.ShowAll()
	return rb
}
