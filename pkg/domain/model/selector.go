package model

import (
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/hatchery/pkg/domain/config"
	"github.com/m-mizutani/hatchery/pkg/domain/types"
)

type Selector struct {
	IDs  []string
	Tags []string
	All  bool
}

func (x *Selector) Validate() error {
	if x.All && len(x.IDs) > 0 {
		return goerr.Wrap(types.ErrInvalidOption, "cannot specify both all and ids")
	}
	if x.All && len(x.Tags) > 0 {
		return goerr.Wrap(types.ErrInvalidOption, "cannot specify both all and tags")
	}
	if len(x.IDs) == 0 && len(x.Tags) == 0 && !x.All {
		return goerr.Wrap(types.ErrInvalidOption, "must specify either ids, tags or all")
	}

	return nil
}

func (x *Selector) Contains(action config.Action) bool {
	if x.All {
		return true
	}

	if len(x.IDs) > 0 {
		for _, id := range x.IDs {
			if id == action.GetId() {
				return true
			}
		}
	}

	if len(x.Tags) > 0 {
		if action.GetTags() != nil {
			for _, tag := range *action.GetTags() {
				for _, xTag := range x.Tags {
					if tag == xTag {
						return true
					}
				}
			}
		}
	}

	return false
}
