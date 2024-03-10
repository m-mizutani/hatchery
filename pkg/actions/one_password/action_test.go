package one_password_test

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
	"github.com/m-mizutani/hatchery/pkg/actions/one_password"
	"github.com/m-mizutani/hatchery/pkg/domain/config"
	"github.com/m-mizutani/hatchery/pkg/infra"
	"github.com/m-mizutani/hatchery/pkg/infra/cs"
	"github.com/m-mizutani/hatchery/pkg/utils"
)

func TestHarvester(t *testing.T) {
	mock := cs.NewMock()

	clients := infra.New(infra.WithCloudStorage(mock))

	ctx := context.Background()
	now := time.Now().Add(-time.Hour * 24)
	ctx = utils.CtxWithNow(ctx, func() time.Time { return now })

	maxPages := 2
	req := &config.OnePasswordImpl{
		ApiToken: utils.LoadEnv(t, "TEST_1PASSWORD_API_TOKEN"),
		Bucket:   "test-bucket",
		Duration: &pkl.Duration{
			Value: 24,
			Unit:  pkl.Hour,
		},
		Limit:    10,
		MaxPages: &maxPages,
	}

	gt.NoError(t, one_password.Exec(ctx, clients, req)).Must()

	expectedPath := fmt.Sprintf(
		"logs/%04d/%02d/%02d/%02d/%s_86400_0.json.gz",
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

	gt.NotEqual(t, resp.Cursor, "")
	gt.A(t, resp.Items).Longer(0).At(0, func(t testing.TB, v Item) {
		gt.NotEqual(t, v.UUID, "")
	})
}

type apiResponse struct {
	Cursor  string `json:"cursor"`
	HasMore bool   `json:"has_more"`
	Items   []Item `json:"items"`
}

type Item struct {
	Action        string   `json:"action"`
	ActorDetails  Details  `json:"actor_details"`
	ActorUUID     string   `json:"actor_uuid"`
	AuxID         int64    `json:"aux_id"`
	AuxInfo       string   `json:"aux_info"`
	AuxUUID       string   `json:"aux_uuid"`
	Location      Location `json:"location"`
	ObjectDetails Details  `json:"object_details"`
	ObjectType    string   `json:"object_type"`
	ObjectUUID    string   `json:"object_uuid"`
	Session       Session  `json:"session"`
	Timestamp     string   `json:"timestamp"`
	UUID          string   `json:"uuid"`
}

type Location struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Region    string  `json:"region"`
}

type Session struct {
	DeviceUUID string `json:"device_uuid"`
	IP         string `json:"ip"`
	LoginTime  string `json:"login_time"`
	UUID       string `json:"uuid"`
}

type Details struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	UUID  string `json:"uuid"`
}
