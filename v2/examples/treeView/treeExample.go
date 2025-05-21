package main

import (
	"github.com/gonutz/wui/v2"
)

func main() {
	window := wui.NewWindow()
	window.SetTitle("TreeView Example")
	window.SetInnerSize(600, 600)

	// Create TreeView
	tree := wui.NewTreeView()
	tree.SetBounds(10, 10, 580, 580)

	// Add some sample items
	root1 := &wui.TreeViewItem{Text: "Root 1"}
	root1.Children = []*wui.TreeViewItem{
		{Text: "Child 1.1"},
		{Text: "Child 1.2", Children: []*wui.TreeViewItem{
			{Text: "Grandchild 1.2.1"},
			{Text: "Grandchild 1.2.2"},
		}},
	}

	root2 := &wui.TreeViewItem{Text: "Root 2"}
	root2.Children = []*wui.TreeViewItem{
		{Text: "Child 2.1"},
		{Text: "Child 2.2"},
	}

	tree.AddItem(root1)
	tree.AddItem(root2)

	// Handle selection
	tree.SetOnSelect(func(item *wui.TreeViewItem) {
		if item != nil {
			window.SetTitle("Selected: " + item.Text)
		}
	})

	window.Add(tree)
	window.Show()
}
