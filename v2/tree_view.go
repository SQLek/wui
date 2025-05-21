package wui

import (
	"syscall"
	"unsafe"

	"github.com/gonutz/w32/v2"
	w32v3 "github.com/gonutz/w32/v3" // indirect
)

var treeViewItems = make(map[uintptr]*TreeViewItem)

// TreeViewItem represents a single item in the tree view
type TreeViewItem struct {
	Text     string
	Children []*TreeViewItem
	Data     interface{}
}

func NewTreeView() *TreeView {
	return &TreeView{
		items:    make([]*TreeViewItem, 0),
		selected: nil,
	}
}

type TreeView struct {
	textControl
	items    []*TreeViewItem
	selected *TreeViewItem
	onSelect func(item *TreeViewItem)
}

var _ Control = (*TreeView)(nil)

func (*TreeView) canFocus() bool {
	return true
}

func (*TreeView) eatsTabs() bool {
	return false
}

func (t *TreeView) create(id int) {
	t.textControl.create(
		id,
		w32.WS_EX_CLIENTEDGE,
		"SysTreeView32",
		w32.WS_TABSTOP|w32v3.TVS_HASLINES|w32v3.TVS_LINESATROOT|w32v3.TVS_HASBUTTONS,
	)
	// Set extended styles
	w32.SendMessage(t.handle, w32v3.TVM_SETEXTENDEDSTYLE, 0,
		w32v3.TVS_EX_DOUBLEBUFFER|w32v3.TVS_EX_FADEINOUTEXPANDOS)

	// Add all items
	for _, item := range t.items {
		t.addItem(w32v3.HTREEITEM(0), item)
	}
}

func (t *TreeView) addItem(parent w32v3.HTREEITEM, item *TreeViewItem) w32v3.HTREEITEM {
	text, _ := syscall.UTF16PtrFromString(item.Text)

	var insertStruct w32v3.TVINSERTSTRUCT
	insertStruct.Parent = parent
	insertStruct.InsertAfter = w32v3.TVI_LAST
	insertStruct.ItemEx.Mask = w32v3.TVIF_TEXT
	insertStruct.ItemEx.Text = text

	hItem := w32v3.HTREEITEM(w32.SendMessage(t.handle, w32v3.TVM_INSERTITEM, 0, uintptr(unsafe.Pointer(&insertStruct))))
	treeViewItems[uintptr(hItem)] = item

	// Add children recursively
	for _, child := range item.Children {
		t.addItem(hItem, child)
	}

	return hItem
}

func (t *TreeView) AddItem(item *TreeViewItem) {
	t.items = append(t.items, item)
	if t.handle != 0 {
		t.addItem(w32v3.HTREEITEM(0), item)
	}
}

func (t *TreeView) Clear() {
	t.items = nil
	t.selected = nil
	if t.handle != 0 {
		w32.SendMessage(t.handle, w32v3.TVM_DELETEITEM, 0, 0)
		treeViewItems = make(map[uintptr]*TreeViewItem)
	}
}

func (t *TreeView) Items() []*TreeViewItem {
	return t.items
}

func (t *TreeView) SetItems(items []*TreeViewItem) {
	t.items = items
	t.selected = nil
	if t.handle != 0 {
		w32.SendMessage(t.handle, w32v3.TVM_DELETEITEM, 0, 0)
		for _, item := range t.items {
			t.addItem(w32v3.HTREEITEM(0), item)
		}
	}
}

func (t *TreeView) SelectedItem() *TreeViewItem {
	if t.handle != 0 {
		hItem := w32v3.HTREEITEM(w32.SendMessage(t.handle, w32v3.TVM_GETNEXTITEM, w32v3.TVGN_CARET, 0))
		if hItem != 0 {
			t.selected = treeViewItems[uintptr(hItem)]
		}
	}
	return t.selected
}

func (t *TreeView) SetSelectedItem(item *TreeViewItem) {
	t.selected = item
	if t.handle != 0 {
		// TODO: Implement finding and selecting the item
	}
}

func (t *TreeView) OnSelect() func(item *TreeViewItem) {
	return t.onSelect
}

func (t *TreeView) SetOnSelect(f func(item *TreeViewItem)) {
	t.onSelect = f
}

func (t *TreeView) handleNotification(cmd uintptr) {
	if cmd == w32v3.TVN_SELCHANGED && t.onSelect != nil {
		t.onSelect(t.SelectedItem())
	}
}
