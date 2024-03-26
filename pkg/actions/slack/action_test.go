package slack_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/apple/pkl-go/pkl"
	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/hatchery/pkg/actions/slack"
	"github.com/m-mizutani/hatchery/pkg/domain/config"
	"github.com/m-mizutani/hatchery/pkg/infra"
	"github.com/m-mizutani/hatchery/pkg/infra/cs"
	"github.com/m-mizutani/hatchery/pkg/utils"
)

func TestAction(t *testing.T) {
	mock := cs.NewMock()

	clients := infra.New(infra.WithCloudStorage(mock))

	ctx := context.Background()
	now := time.Now().Add(-time.Hour * 24)
	ctx = utils.CtxWithNow(ctx, func() time.Time { return now })

	maxPages := 2
	req := &config.SlackImpl{
		AccessToken: utils.LoadEnv(t, "TEST_SLACK_ACCESS_TOKEN"),
		Bucket:      "test-bucket",
		Duration: &pkl.Duration{
			Value: 24,
			Unit:  pkl.Hour,
		},
		Limit:    10,
		MaxPages: &maxPages,
	}

	gt.NoError(t, slack.Exec(ctx, clients, req)).Must()
	expectedPath := fmt.Sprintf(
		"logs/%04d/%02d/%02d/%02d/%s_86400_00000000.json.gz",
		now.Year(), now.Month(), now.Day(), now.Hour(),
		now.Format("20060102T150405"),
	)
	r0 := mock.Results[0]
	gt.Equal(t, r0.Bucket, "test-bucket")
	gt.Equal(t, string(r0.Object), expectedPath)
	gt.Equal(t, r0.Body.Closed, true)

	var resp apiResponse
	r := gt.R1(gzip.NewReader(bytes.NewReader(r0.Body.Bytes()))).NoError(t)
	gt.NoError(t, json.NewDecoder(r).Decode(&resp))

	gt.A(t, resp.Entries).Longer(0)
	gt.V(t, resp.Entries[0].Action).NotEqual("")
}

type apiResponse struct {
	Entries []struct {
		Action string `json:"action"`
	}
}

func TestIntegration(t *testing.T) {
	prefix := "slack-2/"
	req := &config.SlackImpl{
		AccessToken: utils.LoadEnv(t, "TEST_SLACK_ACCESS_TOKEN"),
		Bucket:      utils.LoadEnv(t, "TEST_BUCKET"),
		Prefix:      &prefix,
		Duration: &pkl.Duration{
			Value: 1,
			Unit:  pkl.Hour,
		},
		Limit: 1000,
	}

	ctx := context.Background()
	csClient := gt.R1(cs.New(ctx)).NoError(t)
	clients := infra.New(infra.WithCloudStorage(csClient))

	gt.NoError(t, slack.Exec(ctx, clients, req)).Must()
}
