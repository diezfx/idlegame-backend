package item

type ItemContainer struct {
	items map[string]ItemDefinition
}
type ItemDefinition struct {
	ID   string
	Tags []string
}

func NewContainer() *ItemContainer {
	return &ItemContainer{
		items: map[string]ItemDefinition{},
	}
}

func (ic *ItemContainer) AddItemDefinition(item ItemDefinition) {
	ic.items[item.ID] = item
}
