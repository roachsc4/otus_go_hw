package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int
	Front() *listItem
	Back() *listItem
	PushFront(v interface{}) *listItem
	PushBack(v interface{}) *listItem
	Remove(i *listItem)
	MoveToFront(i *listItem)
}

type listItem struct {
	value interface{}
	Next  *listItem
	Prev  *listItem
	list  *list
}

type list struct {
	head   *listItem
	tail   *listItem
	length int
}

func (l list) Len() int {
	return l.length
}

func (l *list) Front() *listItem {
	return l.head
}

func (l *list) Back() *listItem {
	return l.tail
}

// PushFront - moves item to the front of list.
func (l *list) PushFront(v interface{}) *listItem {
	// In the end of operation we need to increase list length
	defer func() { l.length++ }()

	newItem := &listItem{value: v, list: l}

	// Consider the case when head of list is nil
	if l.head == nil {
		// If head is nil, new item must be set as head and tail
		l.head = newItem
		l.tail = newItem // add
		return newItem
	}
	// Change current head item with target one
	currentFront := l.head
	newItem.Next = currentFront
	if currentFront != nil {
		currentFront.Prev = newItem
	}

	l.head = newItem
	return newItem
}

// PushBack - moves item to the back of list.
func (l *list) PushBack(v interface{}) *listItem {
	defer func() { l.length++ }()

	newItem := &listItem{value: v, list: l}
	// Consider the case when tail of list is nil
	if l.tail == nil {
		// If tail is nil, new item must be set as head and tail
		l.tail = newItem
		l.head = newItem
		return newItem
	}
	// Change current tail item with target one
	currentBack := l.tail
	newItem.Prev = currentBack
	if currentBack != nil {
		currentBack.Next = newItem
	}
	l.tail = newItem
	return newItem
}

// Remove - removes item from list.
// Removal is achieved by redirecting links of neighbour items and clearing links of target item.
func (l *list) Remove(i *listItem) {
	// List can remove only it's own items
	if i.list != l {
		return
	}
	// If target item is not head, than its "prev" item should be linked to target item's "next" item.
	// Otherwise head should be set to target item's next item
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.head = i.Next
	}

	// If target item is not tail, than its "next" item should be back-linked to target item's "prev" item.
	// Otherwise tail should be set to target item's "prev" item
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.tail = i.Prev
	}

	// Clear references
	i.Next = nil
	i.Prev = nil

	l.length--
}

func (l *list) MoveToFront(i *listItem) {
	// List can move only it's own items
	if i.list != l {
		return
	}
	// If target item is head item - no need to do anything
	if i == l.head {
		return
	}

	// If previous item is not nil, than it should be redirected to target item's next item
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	// If next item is not nil, than it's "prev" should be linked with target item's "prev"
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	// If target item is tail and it is not only element in list - new tail should be set
	if i == l.tail && l.tail.Prev != nil {
		l.tail = l.tail.Prev
	}

	l.head.Prev = i
	i.Next = l.head
	l.head = i
}

func NewList() List {
	return &list{}
}
