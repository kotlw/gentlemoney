package ext

import (
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DateField struct {
	*tview.Box

	location             *time.Location
	currentYear          int
	currentMonth         int
	currentDay           int
	yearDropDown         *tview.DropDown
	monthDropDown        *tview.DropDown
	datePalette          *tview.Table
	yearOpen             bool
	monthOpen            bool
	open                 bool
	label                string
	labelWidth           int
	labelColor           tcell.Color
	backgroundColor      tcell.Color
	fieldTextColor       tcell.Color
	fieldBackgroundColor tcell.Color
	finished             func(tcell.Key)
}

func NewDateField() *DateField {
	now := time.Now()
	y, m, d := now.Date()

	df := &DateField{
		Box: tview.NewBox(),

		location:             now.Location(),
		labelColor:           tview.Styles.SecondaryTextColor,
		fieldBackgroundColor: tview.Styles.ContrastBackgroundColor,
		fieldTextColor:       tview.Styles.PrimaryTextColor,
	}

	df.initYearDropDown(2017, y)
	df.initMonthDropDown(int(m))
	df.initDatePalette(d)

	return df
}

func (d *DateField) initYearDropDown(startYear, currentYear int) {
	d.currentYear = currentYear

	opts := rangeStrings(startYear, currentYear+1, "   ")
	d.yearDropDown = tview.NewDropDown().
		SetFieldWidth(10).
		SetOptions(opts, nil).
		SetCurrentOption(len(opts) - 1)

	d.yearDropDown.SetSelectedFunc(func(text string, index int) {
		d.currentYear = startYear + index
		d.refreshDatePalette()
	})
}

func (d *DateField) initMonthDropDown(currentMonth int) {
	d.currentMonth = currentMonth

	opts := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	d.monthDropDown = tview.NewDropDown().
		SetFieldWidth(10).
		SetOptions(opts, nil).
		SetCurrentOption(currentMonth - 1)

	d.monthDropDown.SetSelectedFunc(func(text string, index int) {
		d.currentMonth = index + 1
		d.refreshDatePalette()
	})
}

func (d *DateField) initDatePalette(currentDay int) {
	d.currentDay = currentDay

	d.datePalette = tview.NewTable().SetBorders(false).SetSelectable(true, true)
	d.datePalette.SetSelectionChangedFunc(func(row, column int) {
		day, ok := d.datePalette.GetCell(row, column).GetReference().(int)
		if ok {
			d.currentDay = day
		}
	})
	d.refreshDatePalette()
}

func (d *DateField) refreshDatePalette() {
	d.datePalette.Clear()

	daysInMonth := time.Date(d.currentYear, time.Month(d.currentMonth)+1, 0, 0, 0, 0, 0, d.location).Day()
	offset := (int(time.Date(d.currentYear, time.Month(d.currentMonth), 1, 0, 0, 0, 0, d.location).Weekday()) + 6) % 7

	// create content starting from weekdays titles then empty strings (if month starts not from monday) then days.
	content := make([]string, 7+offset+daysInMonth)
	copy(content, []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"})
	copy(content[7+offset:], rangeStrings(1, daysInMonth+1, ""))

	// set numbers as selectable cells with reference & make other cells not selectable
	for i, v := range content {
		r, c := i/7, i%7
		if i > 6+offset {
			d.datePalette.SetCell(r, c, tview.NewTableCell(v).SetAlign(tview.AlignRight).SetReference(i-offset-6))
		} else {
			d.datePalette.SetCell(r, c, tview.NewTableCell(v).SetSelectable(false))
		}
	}

	i := d.currentDay + offset + 6
	d.datePalette.Select(i/7, i%7)
}

// SetLabel sets the text to be displayed before the input area.
func (d *DateField) SetLabel(label string) *DateField {
	d.label = label
	return d
}

// GetLabel returns the text to be displayed before the input area.
func (d *DateField) GetLabel() string {
	return d.label
}

// SetFormAttributes sets attributes shared by all form items.
func (d *DateField) SetFormAttributes(labelWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) tview.FormItem {
	d.labelWidth = labelWidth
	d.labelColor = labelColor
	d.backgroundColor = bgColor
	d.fieldTextColor = fieldTextColor
	d.fieldBackgroundColor = fieldBgColor
	return d
}

// GetFieldWidth returns this primitive's field screen width.
func (d *DateField) GetFieldWidth() int {
	return 10
}

// GetFieldHeight returns this primitive's field height.
func (d *DateField) GetFieldHeight() int {
	return 1
}

// SetFinishedFunc sets a callback invoked when the user leaves this form item.
func (d *DateField) SetFinishedFunc(handler func(key tcell.Key)) tview.FormItem {
	d.finished = handler
	return d
}

// SetText sets the current text of the input field.
func (d *DateField) SetTextDate(text string) *DateField {
	dt, _ := time.Parse("2006-01-02", text)
	y, m, _d := dt.Date()

	d.currentYear = y
	d.currentMonth = int(m)
	d.currentDay = _d

	d.refreshDatePalette()

	return d
}

// GetText returns the current text of the date field.
func (d *DateField) GetTextDate() string {
	return time.Date(d.currentYear, time.Month(d.currentMonth), d.currentDay, 0, 0, 0, 0, d.location).Format("2006-01-02")
}

// Draw draws this primitive onto the screen.
func (d *DateField) Draw(screen tcell.Screen) {
	d.Box.DrawForSubclass(screen, d)

	// Prepare.
	x, y, width, height := d.GetInnerRect()
	rightLimit := x + width
	if height < 1 || rightLimit <= x {
		return
	}

	// Draw label.
	if d.labelWidth > 0 {
		labelWidth := d.labelWidth
		if labelWidth > rightLimit-x {
			labelWidth = rightLimit - x
		}
		tview.Print(screen, d.label, x, y, labelWidth, tview.AlignLeft, d.labelColor)
		x += labelWidth
	} else {
		_, drawnWidth := tview.Print(screen, d.label, x, y, rightLimit-x, tview.AlignLeft, d.labelColor)
		x += drawnWidth
	}

	// Draw selection area.
	fieldWidth := d.GetFieldWidth()
	if rightLimit-x < fieldWidth {
		fieldWidth = rightLimit - x
	}

	fieldStyle := tcell.StyleDefault.Background(d.fieldBackgroundColor)
	if d.HasFocus() && !d.open {
		fieldStyle = fieldStyle.Background(d.fieldTextColor)
	}
	for index := 0; index < fieldWidth; index++ {
		screen.SetContent(x+index, y, ' ', nil, fieldStyle)
	}

	// Draw selected date.
	color := d.fieldTextColor
	text := time.Date(d.currentYear, time.Month(d.currentMonth), d.currentDay, 0, 0, 0, 0, d.location).Format("2006-01-02")
	// Just show the current selection.
	if d.HasFocus() && !d.open {
		color = d.fieldBackgroundColor
	}
	tview.Print(screen, text, x, y, fieldWidth, tview.AlignLeft, color)

	// Draw calendar.
	if d.HasFocus() && d.open {
		lx := x
		ly := y + 1
		lwidth := 20
		lheight := 2 + d.datePalette.GetRowCount()
		swidth, sheight := screen.Size()
		// We prefer to align the left sides of the list and the main widget, but
		// if there is no space to the right, then shift the list to the left.
		if lx+lwidth >= swidth {
			lx = swidth - lwidth
			if lx < 0 {
				lx = 0
			}
		}
		// We prefer to drop down but if there is no space, maybe drop up?
		if ly+lheight >= sheight && ly-2 > lheight-ly {
			ly = y - lheight
			if ly < 0 {
				ly = 0
			}
		}
		if ly+lheight >= sheight {
			lheight = sheight - ly
		}

		d.yearDropDown.SetRect(lx, ly, 10, 1)
		d.monthDropDown.SetRect(lx+10, ly, 10, 1)
		d.datePalette.SetRect(lx, ly+1, lwidth, lheight-2)

		d.yearDropDown.SetFieldBackgroundColor(tview.Styles.MoreContrastBackgroundColor)
		d.monthDropDown.SetFieldBackgroundColor(tview.Styles.MoreContrastBackgroundColor)
		d.datePalette.SetBackgroundColor(tview.Styles.MoreContrastBackgroundColor)

		datePaletteStyle := tcell.StyleDefault.Background(tview.Styles.MoreContrastBackgroundColor)
		datePaletteStyle = datePaletteStyle.Foreground(tview.Styles.PrimaryTextColor)
		d.datePalette.SetSelectedStyle(datePaletteStyle)
		if d.datePalette.HasFocus() {
			d.datePalette.SetSelectedStyle(tcell.StyleDefault.Background(tview.Styles.PrimaryTextColor))
		}

		if d.yearOpen || d.monthOpen {
			d.yearDropDown.SetFieldBackgroundColor(d.fieldBackgroundColor)
			d.monthDropDown.SetFieldBackgroundColor(d.fieldBackgroundColor)
			d.datePalette.SetBackgroundColor(d.fieldBackgroundColor)
			d.datePalette.SetSelectedStyle(datePaletteStyle.Background(tview.Styles.ContrastBackgroundColor))
		}

		d.datePalette.Draw(screen)
		d.yearDropDown.Draw(screen)
		d.monthDropDown.Draw(screen)
	}
}

// InputHandler returns the handler for this primitive.
func (d *DateField) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {

		// handle opening and closing calendar
		if d.open {
			if event.Key() == tcell.KeyEsc && !d.yearOpen && !d.monthOpen {
				d.open = false
				setFocus(d)
				return
			}
			if event.Key() == tcell.KeyEnter && d.datePalette.HasFocus() {
				d.open = false
				setFocus(d)
				return
			}
		} else {
			if event.Key() == tcell.KeyEnter {
				d.open = true
				setFocus(d.datePalette)
				return
			}
			if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyTab || event.Key() == tcell.KeyBacktab {
				if d.finished != nil {
					d.finished(event.Key())
				}
			}
		}

		// handle year dropdown keys
		if d.yearDropDown.HasFocus() {
			// ignore runes to ignore prefix input
			if event.Key() == tcell.KeyRune {
				return
			}

			if !d.yearOpen {
				// handle keys if dropdown is closed
				if event.Key() == tcell.KeyEnter {
					d.yearOpen = true
				}
				if event.Key() == tcell.KeyTab {
					setFocus(d.monthDropDown)
					return
				}
				if event.Key() == tcell.KeyBacktab {
					setFocus(d.datePalette)
					return
				}
			} else {
				// handle keys if dropdown is opened
				if event.Key() == tcell.KeyEnter || event.Key() == tcell.KeyEsc {
					d.yearOpen = false
				}
			}

			if handler := d.yearDropDown.InputHandler(); handler != nil {
				handler(event, setFocus)
			}
			return
		}

		// handle month dropdown keys
		if d.monthDropDown.HasFocus() {
			// ignore runes to ignore prefix input
			if event.Key() == tcell.KeyRune {
				return
			}

			if !d.monthOpen {
				// handle keys if dropdown is closed
				if event.Key() == tcell.KeyEnter {
					d.monthOpen = true
				}
				if event.Key() == tcell.KeyTab {
					setFocus(d.datePalette)
				}
				if event.Key() == tcell.KeyBacktab {
					setFocus(d.yearDropDown)
				}

			} else {
				// handle keys if dropdown is opened
				if event.Key() == tcell.KeyEnter || event.Key() == tcell.KeyEsc {
					d.monthOpen = false
				}
			}

			if handler := d.monthDropDown.InputHandler(); handler != nil {
				handler(event, setFocus)
			}
			return
		}

		// handle date palette keys
		if d.datePalette.HasFocus() {

			if event.Key() == tcell.KeyTab {
				setFocus(d.yearDropDown)
			}
			if event.Key() == tcell.KeyBacktab {
				setFocus(d.monthDropDown)
			}

			if handler := d.datePalette.InputHandler(); handler != nil {
				handler(event, setFocus)
			}
			return
		}

	})
}

// HasFocus returns whether or not this primitive has focus.
func (d *DateField) HasFocus() bool {
	if d.open {
		return d.yearDropDown.HasFocus() || d.monthDropDown.HasFocus() || d.datePalette.HasFocus()
	}
	return d.Box.HasFocus()
}

func rangeStrings(start, end int, prefix string) []string {
	res := make([]string, 0, end-start)
	for i := start; i < end; i++ {
		res = append(res, prefix+strconv.Itoa(i))
	}
	return res
}
