package ext

import (
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TableDataProvider an interface for table data population.
type TableDataProvider interface {
	GetAll() []map[string]string
}

// Table an extensioun for tview.Table.
type Table struct {
	*tview.Table

	cols          []string
	orderCol      string
	reversedOrder bool
	dataProvider  TableDataProvider
}

// NewTable returns new extended Table.
func NewTable(cols []string, dataProvider TableDataProvider) *Table {
	t := &Table{
		Table:         tview.NewTable(),
		cols:          cols,
		orderCol:      "",
		reversedOrder: false,
		dataProvider:  dataProvider,
	}

	t.SetSelectable(true, false)
	t.Select(0, 0)
	t.SetFixed(1, 0)
	t.SetBorder(true)

	return t
}

// SetOrder sets the order of displayed table data.
func (t *Table) SetOrder(orderCol string, reversed bool) *Table {
	t.orderCol = orderCol
	t.reversedOrder = reversed
	return t
}

// GetSelectedRef returns a reference map[string]string of current selected row.
func (t *Table) GetSelectedRef() map[string]string {
	row, _ := t.GetSelection()
  if res, ok := t.GetCell(row, 0).GetReference().(map[string]string); ok {
    return res
  } else {
    return make(map[string]string)
  }
}

// Refresh refreshes the data and redraws the table.
func (t *Table) Refresh() *Table {
	t.Clear()

	for i, col := range t.cols {
		t.SetCell(0, i,
			tview.NewTableCell(col).
				SetSelectable(false).
				SetAttributes(tcell.AttrBold).
				SetExpansion(10))
	}

	rows := t.dataProvider.GetAll()
	t.sort(rows)

	for i, row := range rows {
		// setting reference for the row
		firtCol := row[t.cols[0]]
		t.SetCell(i+1, 0, tview.NewTableCell(firtCol).SetReference(row))

		// continue with other cols
		for j, col := range t.cols[1:] {
			t.SetCell(i+1, j+1, tview.NewTableCell(row[col]))
		}
	}

	return t
}

// sort internal function to perform sorting depending on state. State sets by SetOrder function.
func (t *Table) sort(s []map[string]string) {
	if t.orderCol != "" {
		sort.Slice(s, func(i, j int) bool {
			if t.reversedOrder {
				return s[i][t.orderCol] > s[j][t.orderCol]
			}
			return s[i][t.orderCol] < s[j][t.orderCol]
		})
	}
}
