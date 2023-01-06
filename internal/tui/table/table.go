package table

import (
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DataProvider interface {
	GetAll() []map[string]string
}

type Table struct {
	*tview.Table

	cols          []string
	dataProvider  DataProvider
	orderCol      string
	reversedOrder bool
}

func New(cols []string, dataProvider DataProvider) *Table {
	t := &Table{
		Table:         tview.NewTable(),
		cols:          cols,
		dataProvider:  dataProvider,
		orderCol:      "",
		reversedOrder: false,
	}

	t.SetSelectable(true, false)
	t.Select(0, 0)
	t.SetFixed(1, 0)
	t.SetBorder(true)

	return t
}

func (t *Table) SetOrder(orderCol string, reversed bool) *Table {
	t.orderCol = orderCol
	t.reversedOrder = reversed
	return t
}

func (t *Table) GetSelectedRef() map[string]string {
	return t.GetCell(t.GetSelection()).GetReference().(map[string]string)
}

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