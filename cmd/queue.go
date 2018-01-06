package cmd

import "github.com/rebuy-de/aws-nuke/resources"

type ItemState int

const (
	ItemStateNew ItemState = iota
	ItemStatePending
	ItemStateWaiting
	ItemStateFailed
	ItemStateFiltered
	ItemStateFinished
)

type Item struct {
	Resource resources.Resource
	State    ItemState
	Region   Region
	Reason   string
	Type     string
}

func (i *Item) Print() {
	switch i.State {
	case ItemStateNew:
		Log(i.Region, i.Type, i.Resource, ReasonWaitPending, "would remove")
	case ItemStatePending:
		Log(i.Region, i.Type, i.Resource, ReasonWaitPending, "triggered remove")
	case ItemStateWaiting:
		Log(i.Region, i.Type, i.Resource, ReasonWaitPending, "waiting")
	case ItemStateFailed:
		Log(i.Region, i.Type, i.Resource, ReasonError, i.Reason)
	case ItemStateFiltered:
		Log(i.Region, i.Type, i.Resource, ReasonSkip, i.Reason)
	case ItemStateFinished:
		Log(i.Region, i.Type, i.Resource, ReasonSuccess, "removed")
	}
}

func (i *Item) List() ([]resources.Resource, error) {
	listers := resources.GetListers()
	return listers[i.Type](i.Region.Session)
}

type Queue []*Item

func (q Queue) CountTotal() int {
	return len(q)
}

func (q Queue) Count(states ...ItemState) int {
	count := 0
	for _, item := range q {
		for _, state := range states {
			if item.State == state {
				count = count + 1
				break
			}
		}
	}
	return count
}
