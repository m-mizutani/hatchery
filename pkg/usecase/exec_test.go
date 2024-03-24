package usecase

import (
	"context"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/hatchery/pkg/domain/config"
	"github.com/m-mizutani/hatchery/pkg/domain/model"
	"github.com/m-mizutani/hatchery/pkg/infra"
)

func tags(tagSet ...string) *[]string {
	return &tagSet
}

func TestExecute(t *testing.T) {
	actions := []config.Action{
		&config.FalconDataReplicatorImpl{
			Id:   "falcon1",
			Tags: tags("tag1", "falcon"),
		},
		&config.OnePasswordImpl{
			Id:   "onepass1",
			Tags: tags("tag1"),
		},
		&config.FalconDataReplicatorImpl{
			Id:   "falcon2",
			Tags: tags("tag2", "falcon"),
		},
		&config.FalconDataReplicatorImpl{
			Id: "falcon3",
		},
	}

	testCases := map[string]struct {
		selector *model.Selector
		expected []string
	}{
		"Select by ID": {
			selector: &model.Selector{
				IDs: []string{"falcon2", "onepass1"},
			},
			expected: []string{"falcon2", "onepass1"},
		},
		"Select by tag": {
			selector: &model.Selector{
				Tags: []string{"tag1"},
			},
			expected: []string{"falcon1", "onepass1"},
		},
		"Select by all": {
			selector: &model.Selector{
				All: true,
			},
			expected: []string{"falcon1", "onepass1", "falcon2", "falcon3"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var executed []string
			execFn := func(ctx context.Context, clients *infra.Clients, action config.Action) error {
				executed = append(executed, action.GetId())
				return nil
			}

			if err := Execute(context.Background(), &infra.Clients{}, actions, tc.selector, WithExecFn(execFn)); err != nil {
				t.Fatal(err)
			}

			at := gt.A(t, executed).Length(len(tc.expected))
			for _, id := range tc.expected {
				at.Have(id)
			}
		})
	}
}
