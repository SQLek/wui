// Package wui provides a Windows user interface library.
package wui

import (
	"syscall"
	"unsafe"

	"github.com/gonutz/w32/v2"
	w32v3 "github.com/gonutz/w32/v3" // indirect
)

// TreeViewItem represents a single item in the tree view control.
// Each item can have a text label, optional data, and child items.
type TreeViewItem struct {
	Text     string          // The text label displayed for this item
	Children []*TreeViewItem // Child items of this node
	Data     interface{}     // Optional user data associated with this item
}

// NewTreeView creates and returns a new TreeView control.
// The control is initialized with an empty list of items.
func NewTreeView() *TreeView {
	return &TreeView{
		items:    make([]*TreeViewItem, 0),
		selected: nil,
	}
}

// TreeView is a control that displays hierarchical data in a tree structure.
// It supports expanding/collapsing nodes and item selection.
type TreeView struct {
	textControl
	items    []*TreeViewItem
	selected *TreeViewItem
	onSelect func(item *TreeViewItem)
}

var _ Control = (*TreeView)(nil)

// canFocus returns true as TreeView can receive keyboard focus.
func (*TreeView) canFocus() bool {
	return true
}

// eatsTabs returns false as TreeView should not consume tab key events.
func (*TreeView) eatsTabs() bool {
	return false
}

// create initializes the TreeView control with the specified ID.
// It sets up the control's appearance and initializes its items.
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

// addItem adds a single item to the tree view.
// It returns the handle to the created tree item.
func (t *TreeView) addItem(parent w32v3.HTREEITEM, item *TreeViewItem) w32v3.HTREEITEM {
	text, _ := syscall.UTF16PtrFromString(item.Text)

	var insertStruct w32v3.TVINSERTSTRUCT
	insertStruct.Parent = parent
	insertStruct.InsertAfter = w32v3.TVI_LAST
	insertStruct.ItemEx.Mask = w32v3.TVIF_TEXT
	insertStruct.ItemEx.Text = text

	hItem := w32v3.HTREEITEM(w32.SendMessage(t.handle, w32v3.TVM_INSERTITEM, 0, uintptr(unsafe.Pointer(&insertStruct))))

	// Store the item data
	var tvItem w32v3.TVITEM
	tvItem.Mask = w32v3.TVIF_PARAM
	tvItem.Item = hItem
	tvItem.LParam = uintptr(unsafe.Pointer(item))
	w32.SendMessage(t.handle, w32v3.TVM_SETITEM, 0, uintptr(unsafe.Pointer(&tvItem)))

	// Add children recursively
	for _, child := range item.Children {
		t.addItem(hItem, child)
	}

	return hItem
}

// AddItem adds a new item to the root level of the tree view.
func (t *TreeView) AddItem(item *TreeViewItem) {
	t.items = append(t.items, item)
	if t.handle != 0 {
		t.addItem(w32v3.HTREEITEM(0), item)
	}
}

// Clear removes all items from the tree view.
func (t *TreeView) Clear() {
	t.items = nil
	t.selected = nil
	if t.handle != 0 {
		w32.SendMessage(t.handle, w32v3.TVM_DELETEITEM, 0, 0)
	}
}

// Items returns the current list of root-level items in the tree view.
func (t *TreeView) Items() []*TreeViewItem {
	return t.items
}

// SetItems replaces all items in the tree view with the provided items.
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

// SelectedItem returns the currently selected item in the tree view.
// Returns nil if no item is selected.
func (t *TreeView) SelectedItem() *TreeViewItem {
	if t.handle != 0 {
		hItem := w32v3.HTREEITEM(w32.SendMessage(t.handle, w32v3.TVM_GETNEXTITEM, w32v3.TVGN_CARET, 0))
		if hItem != 0 {
			var item w32v3.TVITEM
			item.Mask = w32v3.TVIF_PARAM
			item.Item = hItem
			w32.SendMessage(t.handle, w32v3.TVM_GETITEM, 0, uintptr(unsafe.Pointer(&item)))
			t.selected = (*TreeViewItem)(unsafe.Pointer(item.LParam))
		}
	}
	return t.selected
}

// SetSelectedItem sets the currently selected item in the tree view.
// Note: This is currently not implemented.
func (t *TreeView) SetSelectedItem(item *TreeViewItem) {
	t.selected = item
	if t.handle != 0 {
		// TODO: Implement finding and selecting the item
	}
}

// OnSelect returns the current selection change callback function.
func (t *TreeView) OnSelect() func(item *TreeViewItem) {
	return t.onSelect
}

// SetOnSelect sets the callback function that is called when the selection changes.
func (t *TreeView) SetOnSelect(f func(item *TreeViewItem)) {
	t.onSelect = f
}

// handleNotification processes control-specific notifications.
func (t *TreeView) handleNotification(cmd uintptr) {
	if cmd == w32v3.TVN_SELCHANGED && t.onSelect != nil {
		t.onSelect(t.SelectedItem())
	}
}
