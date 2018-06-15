package cmd

import (
	"github.com/rebuy-de/aws-nuke/resources"
	"github.com/sirupsen/logrus"
)

type ItemState int

// States of Items based on the latest request to AWS.
const (
	ItemStateNew ItemState = iota
	ItemStatePending
	ItemStateWaiting
	ItemStateFailed
	ItemStateFiltered
	ItemStateFinished
)

// An Item describes an actual AWS resource entity with the current state and
// some metadata.
type Item struct {
	Resource resources.Resource

	State  ItemState
	Reason string

	Region Region
	Type   string
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

// List gets all resource items of the same resource type like the Item.
func (i *Item) List() ([]resources.Resource, error) {
	listers := resources.GetListers()
	return listers[i.Type](i.Region.Session)
}

func (i *Item) GetProperty(key string) string {
	if key == "" {
		return i.Resource.String()
	}

	getter, ok := i.Resource.(resources.ResourcePropertyGetter)
	if !ok {
		logrus.Warnf("%T does not support custom properties", i.Resource)
		return ""
	}

	return getter.Properties().Get(key)
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
