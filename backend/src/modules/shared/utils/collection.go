package shared_utils

type EntityCollection[Item Entity] struct {
	items        map[string]Item
	newItems     []Item
	dirtyItems   []Item
	removedItems []Item
}

type Entity interface {
	GetId() string
}

func (collection *EntityCollection[Item]) Add(item Item) {
	if _, ok := collection.items[item.GetId()]; ok {
		collection.items[item.GetId()] = item
		collection.dirtyItems = append(collection.dirtyItems, item)
		return
	}

	collection.items[item.GetId()] = item
	collection.newItems = append(collection.newItems, item)
}

func (collection *EntityCollection[Item]) Remove(item Item) {
	if _, ok := collection.items[item.GetId()]; !ok {
		return
	}

	delete(collection.items, item.GetId())
	collection.removedItems = append(collection.removedItems, item)
}

func (collection *EntityCollection[Item]) GetItems() []Item {
	items := make([]Item, len(collection.items))

	for _, item := range collection.items {
		items = append(items, item)
	}

	return items
}

func (collection EntityCollection[Item]) GetDirtyItems() []Item {
	return collection.dirtyItems
}

func (collection *EntityCollection[Item]) GetRemovedItems() []Item {
	return collection.removedItems
}

func (collection *EntityCollection[Item]) GetNewItems() []Item {
	return collection.newItems
}

func (collection EntityCollection[Item]) Length() int {
	return len(collection.items)
}

func (collection EntityCollection[Item]) Get(id string) *Item {
	item, ok := collection.items[id]

	if !ok {
		return nil
	}

	return &item
}

func NewCollection[Item Entity](items []Item) EntityCollection[Item] {
	itemsMap := make(map[string]Item, len(items))

	for _, item := range items {
		itemsMap[item.GetId()] = item
	}

	return EntityCollection[Item]{items: itemsMap}
}
