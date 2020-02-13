package mbox

import (
	"errors"
	//"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const hma string = "<span><tt><b>"
const hmb string = "</b></tt></span>"
const dmae string = "<span background=\"white\"><tt>"
const dmbe string = "</tt></span>"
const dmao string = "<span background=\"#e8e8e8\"><tt>"
const dmbo string = "</tt></span>"
const cma string = "<span foreground=\"white\" background=\"#6666dd\"><tt>"
const cmb string = "</tt></span>"
const OUTPOS int = -1

// Global vars
var headMarkup = [2]string{hma, hmb}
var dataMarkupE = [2]string{dmae, dmbe}
var dataMarkupO = [2]string{dmao, dmbo}
var cursorMarkup = [2]string{cma, cmb}

var RowSep = 2
var ColumnSep = 2
var LeftRightMargin = 1

var cursorPosition = OUTPOS
var lastPosition = OUTPOS

// Functions
func generateX(elem, align string, grow int) string {
	if grow < 0 {
		grow = 0
	}

	sep := strings.Repeat(" ", LeftRightMargin)

	if align == "left" {
		return sep + elem + strings.Repeat(" ", grow) + sep

	} else if align == "rigth" {
		return sep + strings.Repeat(" ", grow) + elem + sep

	} else if align == "center" {
		a := grow / 2
		b := grow / 2
		if grow%2 != 0 {
			b++
		}
		return sep + strings.Repeat(" ", a) + elem + strings.Repeat(" ", b) + sep
	} else {
		return sep + elem + strings.Repeat(" ", grow) + sep
	}
}

func updateColsWidth(elems *[]string, colsWidth *[]int) {
	for i := 0; i < len(*elems); i++ {
		if len(*colsWidth) < len(*elems) {
			*colsWidth = append(*colsWidth, 0)
		}
		if (*colsWidth)[i] < utf8.RuneCountInString((*elems)[i]) {
			(*colsWidth)[i] = utf8.RuneCountInString((*elems)[i])
		}
	}
}

func defaultAligns(n int) []string {
	var aligns []string

	for i := 0; i < n; i++ {
		aligns = append(aligns, "left")
	}

	return aligns
}

func setDataMarkup(label *gtk.Label, namex string, n int) {
	if n % 2 == 0 {
		label.SetMarkup(dataMarkupE[0] + namex + dataMarkupE[1])
	} else {
		label.SetMarkup(dataMarkupO[0] + namex + dataMarkupO[1])
	}
}

func setDataMarkups(n int, row *RowBox2) {
	if n % 2 == 0 {
		for i := 0; i < len(row.Items); i++ {
			row.Items[i].Label.SetMarkup(dataMarkupE[0] + row.Items[i].Namex + dataMarkupE[1])
		}

	} else {
		for i := 0; i < len(row.Items); i++ {
			row.Items[i].Label.SetMarkup(dataMarkupO[0] + row.Items[i].Namex + dataMarkupO[1])
		}
	}
}

// Head
type HeadItem struct {
	Name     string
	Namex    string
	Align    string
	EventBox *gtk.EventBox
	Label    *gtk.Label
}

func NewHeadItem(name string, width int, align string) HeadItem {
	grow := width - utf8.RuneCountInString(name)
	namex := generateX(name, align, grow)

	ebox, _ := gtk.EventBoxNew()
	label, _ := gtk.LabelNew(name)
	label.SetMarkup(headMarkup[0] + namex + headMarkup[1])
	ebox.Add(label)

	return HeadItem{name, namex, align, ebox, label}
}

func (hi *HeadItem) SetAlign(align string) {
	hi.Align = align
	grow := utf8.RuneCountInString(hi.Namex) -
		utf8.RuneCountInString(hi.Name) - (LeftRightMargin * 2)

	hi.Namex = generateX(hi.Name, hi.Align, grow)
	hi.Label.SetMarkup(headMarkup[0] + hi.Namex + headMarkup[1])
}

func (hi *HeadItem) refreshWidth(width int) {
	grow := width - utf8.RuneCountInString(hi.Name)
	hi.Namex = generateX(hi.Name, hi.Align, grow)
	hi.Label.SetMarkup(headMarkup[0] + hi.Namex + headMarkup[1])
}

type Header struct {
	Items []HeadItem
	Box   *gtk.Box
}

func NewHeader(names []string, widths []int, aligns []string) Header {
	var items []HeadItem
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, ColumnSep)
	box.SetMarginBottom(RowSep)

	for i, name := range names {
		items = append(items, NewHeadItem(name, widths[i], aligns[i]))
	}

	for _, item := range items {
		box.Add(item.EventBox)
	}

	return Header{items, box}
}

func (h *Header) reset() {
	for i := 0; i < len(h.Items); i++ {
		h.Items[i].refreshWidth(0)
	}
}

// Tablam
type DataItem struct {
	Name     string
	Namex    string
	Align    string
	EventBox *gtk.EventBox
	Label    *gtk.Label
}

func NewDataItem(name, align string, n, width int) DataItem {
	grow := width - utf8.RuneCountInString(name)
	namex := generateX(name, align, grow)

	ebox, _ := gtk.EventBoxNew()
	label, _ := gtk.LabelNew(namex)
	setDataMarkup(label, namex, n)
	ebox.Add(label)

	return DataItem{name, namex, align, ebox, label}
}

func (di *DataItem) refreshWidth(n, width int) {
	grow := width - utf8.RuneCountInString(di.Name)
	di.Namex = generateX(di.Name, di.Align, grow)
	setDataMarkup(di.Label, di.Namex, n)
}

func (di *DataItem) edit(name string, width int) (bool, int) {
	var changed bool
	di.Name = name

	nwidth := utf8.RuneCountInString(di.Name)
	if nwidth > width {
		width = nwidth
		changed = true
	}

	grow := width - nwidth
	di.Namex = generateX(di.Name, di.Align, grow)
	di.Label.SetMarkup(cursorMarkup[0] + di.Namex + cursorMarkup[1])

	return changed, nwidth
}

type RowBox2 struct {
	Items []DataItem
	Box   *gtk.Box
}

func NewRowBox2(id int, items []DataItem, t *Tablam) RowBox2 {
	var rb = RowBox2{items, nil}

	rb.Box, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, ColumnSep)
	rb.Box.SetName(strconv.Itoa(id))

	for _, item := range rb.Items {
		rb.Box.Add(item.EventBox)
	}

	rb.Box.Connect("button-press-event", func(box *gtk.Box, e *gdk.Event) bool {
		name, _ := box.GetName()
		namint, _ := strconv.Atoi(name)
		lastPosition = cursorPosition
		cursorPosition = namint

		t.updateCursor()
		return false
	})

	rb.Box.ShowAll()
	return rb
}

type Tablam struct {
	head        Header
	rows        []RowBox2
	colsWidth   []int
	colsChanged []int
	aligns      []string
	Grid        *gtk.Grid
	Box         *gtk.Box
}

func NewTablam(titles, aligns []string) Tablam {
	var head Header

	var t = Tablam{
		head,
		nil,
		nil,
		nil,
		aligns,
		nil,
		nil,
	}

	updateColsWidth(&titles, &t.colsWidth)

	if t.aligns == nil && t.colsWidth != nil {
		t.aligns = defaultAligns(len(t.colsWidth))
	}

	t.Box, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	t.Box.SetHAlign(gtk.ALIGN_CENTER)

	t.Grid, _ = gtk.GridNew()
	t.Grid.SetHAlign(gtk.ALIGN_CENTER)
	t.Grid.SetRowSpacing(uint(RowSep))

	if titles != nil {
		t.head = NewHeader(titles, t.colsWidth, t.aligns)
		t.Box.Add(t.head.Box)
	}
	t.Box.Add(t.Grid)

	return t
}

func (t *Tablam) CursorDown() int {
	if len(t.rows) > 0 {
		lastPosition = cursorPosition

		if cursorPosition < len(t.rows)-1 {
			cursorPosition++
		} else {
			cursorPosition = 0
		}

		t.updateCursor()
	}

	return cursorPosition
}

func (t *Tablam) CursorUp() int {
	lastPosition = cursorPosition
	cursorPosition--

	if cursorPosition < 0 {
		cursorPosition = len(t.rows) - 1
	}
	if cursorPosition >= 0 {
		t.updateCursor()
	}

	return cursorPosition
}

func (t *Tablam) updateCursor() {
	if lastPosition > OUTPOS {
		setDataMarkups(lastPosition, &t.rows[lastPosition])
	}

	if cursorPosition >= OUTPOS {
		for i := 0; i < len(t.rows[0].Items); i++ {
			t.rows[cursorPosition].Items[i].Label.SetMarkup(
				cursorMarkup[0] + t.rows[cursorPosition].Items[i].Namex + cursorMarkup[1])
		}
	}
}

func (t *Tablam) markupActiveRow() {
	for i := 0; i < len(t.rows[0].Items); i++ {
		t.rows[cursorPosition].Items[i].Label.SetMarkup(
			cursorMarkup[0] + t.rows[cursorPosition].Items[i].Namex + cursorMarkup[1])
	}
}

func (t *Tablam) SetHeadAligns(aligns []string) error {
	if t.head.Items != nil {
		for i, align := range aligns {
			t.head.Items[i].SetAlign(align)
		}
		return nil

	} else {
		return errors.New("Header not defined")
	}
}

func (t *Tablam) AddRow(rdata []string) {
	if t.head.Items == nil && t.rows == nil {
		updateColsWidth(&rdata, &t.colsWidth)

		if t.aligns == nil {
			t.aligns = defaultAligns(len(t.colsWidth))
		}
	}

	var rowItems []DataItem

	for i, elem := range rdata {
		rowItems = append(rowItems, NewDataItem(elem, t.aligns[i], len(t.rows), t.colsWidth[i]))
	}

	row := NewRowBox2(len(t.rows), rowItems, t)
	t.Grid.Attach(row.Box, 0, len(t.rows), 1, 1)
	t.rows = append(t.rows, row)

	t.colsChanged = []int{}

	for i, rd := range rdata {
		drunes := utf8.RuneCountInString(rd)

		if drunes > t.colsWidth[i] {
			t.colsWidth[i] = drunes
			t.colsChanged = append(t.colsChanged, i)
		}
	}

	t.refreshLabels()

	if cursorPosition > OUTPOS && cursorPosition < len(t.rows) {
		t.updateCursor()
	}
}

func (t *Tablam) EditActiveRow(edata []string) {
	t.colsChanged = nil

	for i := 0; i < len(t.rows[cursorPosition].Items); i++ {
		changed, nwidth := t.rows[cursorPosition].Items[i].edit(edata[i], t.colsWidth[i])

		if changed {
			t.colsWidth[i] = nwidth
			t.colsChanged = append(t.colsChanged, i)
		}
	}

	t.refreshLabels()

	if cursorPosition > OUTPOS {
		t.updateCursor()
	}
}

func (t *Tablam) DeleteActiveRow() {
	if cursorPosition > OUTPOS {

		t.rows = append(t.rows[:cursorPosition], t.rows[cursorPosition+1:]...)
		t.Grid.RemoveRow(cursorPosition)

		if len(t.rows) == OUTPOS+1 {
			cursorPosition = OUTPOS
		} else if cursorPosition == len(t.rows) {
			cursorPosition--
		}

		t.refreshDataMarkup()

		if cursorPosition > OUTPOS {
			t.markupActiveRow()
			t.UpdateBoxNames()
		}
	}
}

func (t *Tablam) DeleteAll() {
	if len(t.rows) > 0 {

		for i := 0; i < len(t.rows); i++ {
			t.Grid.RemoveRow(0)
		}

		t.rows = []RowBox2{}
		cursorPosition = OUTPOS
		lastPosition = OUTPOS

		for i := 0; i < len(t.colsWidth); i++ {
			t.colsWidth[i] = 0
		}

		t.colsChanged = []int{}
		t.head.reset()
	}
}

func (t *Tablam) refreshLabels() {
	for _, n := range t.colsChanged {
		if t.head.Items != nil {
			t.head.Items[n].refreshWidth(t.colsWidth[n])
		}

		for i := 0; i < len(t.rows); i++ {
			t.rows[i].Items[n].refreshWidth(i, t.colsWidth[n])
		}
	}
}

func (t *Tablam) CursorIsActive() bool {
	if cursorPosition > OUTPOS {
		return true
	} else {
		return false
	}
}

func (t *Tablam) ClearCursor() int {
	if cursorPosition > OUTPOS {
		setDataMarkups(cursorPosition, &t.rows[cursorPosition])
		cursorPosition = OUTPOS
	}

	return cursorPosition
}

func (t *Tablam) ActiveData() []string {
	if cursorPosition > OUTPOS {
		var rowData []string

		for _, item := range t.rows[cursorPosition].Items {
			rowData = append(rowData, item.Name)
		}
		return rowData

	} else {
		return nil
	}
}

func (t *Tablam) UpdateBoxNames() {
	for i := 0; i < len(t.rows); i++ {
		t.rows[i].Box.SetName(strconv.Itoa(i))
	}
}

func (t Tablam) SetHeadMarkup(a, b string) {
	headMarkup = [2]string{a, b}
}

func (t Tablam) SetDataMarkupEven(a, b string) {
	dataMarkupE = [2]string{a, b}
}

func (t Tablam) SetDataMarkupOdd(a, b string) {
	dataMarkupO = [2]string{a, b}
}

func (t Tablam) SetCursorMarkup(a, b string) {
	cursorMarkup = [2]string{a, b}
}

func (t Tablam) SetRowSeparation(sep int) {
	RowSep = sep
}

func (t Tablam) SetColumnSeparation(sep int) {
	ColumnSep = sep
}

func (t Tablam) SetLeftAndRightMargin(margin int) {
	LeftRightMargin = margin
}

func (t Tablam) GetCursorPosition() int {
	return cursorPosition
}

func (t *Tablam) refreshDataMarkup() {
	for i := 0; i < len(t.rows); i++ {
		setDataMarkups(i, &t.rows[i])
	}
}
