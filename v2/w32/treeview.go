package w32

// TreeView styles
const (
	TVS_HASLINES             = 0x0002
	TVS_LINESATROOT          = 0x0004
	TVS_HASBUTTONS           = 0x0001
	TVS_EX_DOUBLEBUFFER      = 0x0004
	TVS_EX_FADEINOUTEXPANDOS = 0x0040
)

// TreeView messages
const (
	TVM_FIRST            = 0x1100
	TVM_INSERTITEM       = TVM_FIRST + 50
	TVM_DELETEITEM       = TVM_FIRST + 1
	TVM_GETNEXTITEM      = TVM_FIRST + 10
	TVM_GETITEM          = TVM_FIRST + 62
	TVM_SETITEM          = TVM_FIRST + 13
	TVM_SETEXTENDEDSTYLE = TVM_FIRST + 44
)

// TreeView item flags
const (
	TVIF_TEXT  = 0x0001
	TVIF_PARAM = 0x0004
)

// TreeView navigation flags
const (
	TVGN_ROOT            = 0x0000
	TVGN_NEXT            = 0x0001
	TVGN_PREVIOUS        = 0x0002
	TVGN_PARENT          = 0x0003
	TVGN_CHILD           = 0x0004
	TVGN_FIRSTVISIBLE    = 0x0005
	TVGN_NEXTVISIBLE     = 0x0006
	TVGN_PREVIOUSVISIBLE = 0x0007
	TVGN_DROPHILITE      = 0x0008
	TVGN_CARET           = 0x0009
	TVGN_LASTVISIBLE     = 0x000A
)

// TreeView notifications
const (
	TVN_FIRST      = -400
	TVN_SELCHANGED = TVN_FIRST - 1
)

// TreeView item insertion constants
const (
	TVI_ROOT  = 0xFFFF0000
	TVI_FIRST = 0xFFFF0001
	TVI_LAST  = 0xFFFF0002
	TVI_SORT  = 0xFFFF0003
)

// HTREEITEM represents a handle to a tree view item
type HTREEITEM uintptr

// TVINSERTSTRUCT contains information used to add a new item to a tree-view control
type TVINSERTSTRUCT struct {
	HParent      HTREEITEM
	HInsertAfter HTREEITEM
	Item         TVITEM
}

// TVITEM specifies or receives attributes of a tree-view item
type TVITEM struct {
	Mask           uint32
	HItem          HTREEITEM
	State          uint32
	StateMask      uint32
	PszText        *uint16
	CchTextMax     int32
	IImage         int32
	ISelectedImage int32
	CChildren      int32
	LParam         uintptr
}
