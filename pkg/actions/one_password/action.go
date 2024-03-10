package one_password

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/hatchery/pkg/domain/config"
	"github.com/m-mizutani/hatchery/pkg/domain/model"
	"github.com/m-mizutani/hatchery/pkg/infra"
	"github.com/m-mizutani/hatchery/pkg/utils"
)

const (
	// 1Password API endpoint for Business Plan
	// See https://developer.1password.com/docs/events-api/reference/
	APIEndpoint = "https://events.1password.com/api/v1/auditevents"

	// Time format for 1Password API
	// 2023-03-15T16:32:50-03:00
	timeFormat = "2006-01-02T15:04:05-07:00"
)

func Exec(ctx context.Context, clients *infra.Clients, req *config.OnePasswordImpl) error {
	var nextCursor string
	now := utils.CtxNow(ctx)

	for seq := 0; req.MaxPages == nil || seq < *req.MaxPages; seq++ {
		cursor, err := crawl(ctx, clients, req, now, seq, nextCursor)
		if err != nil {
			return goerr.Wrap(err, "failed to crawl 1Password logs").With("seq", seq).With("cursor", nextCursor).With("req", req)
		}
		if cursor == nil {
			break
		}
		nextCursor = *cursor
	}

	return nil
}

func crawl(ctx context.Context, clients *infra.Clients, req *config.OnePasswordImpl, end time.Time, seq int, cursor string) (*string, error) {
	d := req.GetDuration().GoDuration()

	objPrefix := config.ToObjNamePrefix(req, end)
	objName := model.CSObjectName(
		fmt.Sprintf("%s_%d_%d.json.gz", objPrefix, d/time.Second, seq),
	)
	objWriter := clients.CloudStorage().NewObjectWriter(ctx, model.CSBucket(req.GetBucket()), objName)
	w := gzip.NewWriter(objWriter)

	startTime := end.Add(-d)
	var body []byte
	if cursor != "" {
		raw, err := json.Marshal(apiResponseWithCursor{Cursor: cursor})
		if err != nil {
			return nil, goerr.Wrap(err, "failed to marshal API request")
		}
		body = raw
	} else {
		raw, err := json.Marshal(apiRequest{
			Limit:     req.GetLimit(),
			StartTime: startTime.Format(timeFormat),
			EndTime:   end.Format(timeFormat),
		})
		if err != nil {
			return nil, goerr.Wrap(err, "failed to marshal API request")
		}
		body = raw
	}
	reader := bytes.NewReader(body)

	httpReq, err := http.NewRequest(http.MethodPost, APIEndpoint, reader)
	if err != nil {
		return nil, goerr.Wrap(err, "failed to create HTTP request")
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+req.GetApiToken())

	httpResp, err := clients.HTTPClient().Do(httpReq)
	if err != nil {
		return nil, goerr.Wrap(err, "failed to send HTTP request")
	}

	if httpResp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(httpResp.Body)
		return nil, goerr.New("unexpected status code").With("status", httpResp.Status).With("body", string(data))
	}

	body, err = io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, goerr.Wrap(err, "failed to read response body")
	}

	var resp apiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, goerr.Wrap(err, "failed to unmarshal response body")
	}

	n, err := io.Copy(w, bytes.NewReader(body))
	if err != nil {
		return nil, goerr.Wrap(err, "failed to write response to object writer").With("bytes", n)
	}

	utils.CtxLogger(ctx).Info("harvested 1Password logs", "bytes", n, "object", objName, "cursor", resp.Cursor, "hasMore", resp.HasMore)

	if err := w.Close(); err != nil {
		return nil, goerr.Wrap(err, "failed to close gzip writer").With("object", objName)
	}
	if err := objWriter.Close(); err != nil {
		return nil, goerr.Wrap(err, "failed to close object writer").With("object", objName)
	}

	if resp.HasMore {
		return &resp.Cursor, nil
	}

	return nil, nil
}

type apiRequest struct {
	Limit     int    `json:"limit"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
}

type apiResponseWithCursor struct {
	Cursor string `json:"cursor"`
}

type apiResponse struct {
	Cursor  string `json:"cursor"`
	HasMore bool   `json:"has_more"`
}
