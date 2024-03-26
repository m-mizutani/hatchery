package slack

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/hatchery/pkg/domain/config"
	"github.com/m-mizutani/hatchery/pkg/domain/model"
	"github.com/m-mizutani/hatchery/pkg/infra"
	"github.com/m-mizutani/hatchery/pkg/utils"
)

func Exec(ctx context.Context, clients *infra.Clients, req config.Slack) error {
	var nextCursor string
	now := utils.CtxNow(ctx)

	for seq := 0; req.GetMaxPages() == nil || seq < *req.GetMaxPages(); seq++ {
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

const (
	// Slack API endpoint for Business Plan
	baseURL = "https://api.slack.com/audit/v1/logs"
)

func crawl(ctx context.Context, clients *infra.Clients, req config.Slack, end time.Time, seq int, cursor string) (*string, error) {
	d := req.GetDuration().GoDuration()

	objName := model.CSObjectName(
		fmt.Sprintf("%s_%d_%08d.json.gz", end.Format("20060102T150304"), d/time.Second, seq),
	)
	objWriter := clients.CloudStorage().NewObjectWriter(ctx,
		model.CSBucket(req.GetBucket()),
		model.LogObjNamePrefix(req, end)+objName,
	)
	w := gzip.NewWriter(objWriter)

	startTime := end.Add(-d)
	qv := url.Values{}
	qv.Add("limit", fmt.Sprintf("%d", req.GetLimit()))
	qv.Add("oldest", fmt.Sprintf("%d", startTime.Unix()))

	if cursor != "" {
		qv.Add("cursor", cursor)
	}

	endpoint, err := url.Parse(baseURL)
	if err != nil {
		return nil, goerr.Wrap(err, "failed to parse URL").With("url", baseURL)
	}
	endpoint.RawQuery = qv.Encode()

	apiURL := endpoint.String()
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, goerr.Wrap(err, "failed to create HTTP request")
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+req.GetAccessToken())

	httpResp, err := clients.HTTPClient().Do(httpReq)
	if err != nil {
		return nil, goerr.Wrap(err, "failed to send HTTP request")
	}

	if httpResp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(httpResp.Body)
		return nil, goerr.New("unexpected status code").With("status", httpResp.Status).With("body", string(data))
	}

	body, err := io.ReadAll(httpResp.Body)
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

	if err := w.Close(); err != nil {
		return nil, goerr.Wrap(err, "failed to close gzip writer").With("object", objName)
	}
	if err := objWriter.Close(); err != nil {
		return nil, goerr.Wrap(err, "failed to close object writer").With("object", objName)
	}

	if resp.ResponseMetadata.NextCursor != "" {
		return &resp.ResponseMetadata.NextCursor, nil
	}

	return nil, nil
}

type apiResponse struct {
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}
