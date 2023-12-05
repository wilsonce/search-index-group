package search_index_group

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"google.golang.org/api/indexing/v3"
	"google.golang.org/api/option"
	"net/url"
	"os"
	"time"
)

func NewSearchIndex() (*SearchIndex, error) {
	ctx := context.Background()

	gTokenFile := viper.GetString("token.google.searchFile")
	gToken, err := os.ReadFile(gTokenFile)
	if err != nil {
		return nil, err
	}
	service, err := indexing.NewService(ctx, option.WithCredentialsJSON(gToken))
	if err != nil {
		return nil, err
	}
	return &SearchIndex{
		googleSearchIndex: service,
		restyClient:       resty.New(),
		bingApiHost:       "https://ssl.bing.com",
		bingApiKey:        viper.GetString("token.bing.apiKey"),
	}, nil
}

type SearchIndex struct {
	googleSearchIndex *indexing.Service

	restyClient *resty.Client

	bingApiHost string
	bingApiKey  string
}

type Platform int

const (
	Google Platform = iota
	Bing
)

type UpdateIndexResult struct {
	Notification *UrlNotification
	Error        error
}

func (s *SearchIndex) PublishUrl(uri string) []*UpdateIndexResult {
	var result []*UpdateIndexResult
	if s.googleSearchIndex != nil {
		notification, err := s.googleUpdateIndex(uri)
		result = append(result, &UpdateIndexResult{
			Notification: notification,
			Error:        err,
		})
	}
	if len(s.bingApiKey) > 0 {
		notification, err := s.bingUpdateIndex(uri)
		result = append(result, &UpdateIndexResult{
			Notification: notification,
			Error:        err,
		})
	}
	return result
}

func (s *SearchIndex) googleUpdateIndex(uri string) (*UrlNotification, error) {
	updateTime := time.Now().Format(time.RFC3339)
	reqNotification := &indexing.UrlNotification{
		Url:        uri,
		NotifyTime: updateTime,
		Type:       "URL_UPDATED",
	}
	publish := s.googleSearchIndex.UrlNotifications.Publish(reqNotification)
	rsp, err := publish.Do()
	if err != nil {
		return nil, err
	}
	rspNotification := &UrlNotification{
		Url:            uri,
		NotifyTime:     updateTime,
		Type:           "URL_UPDATED",
		HTTPStatusCode: rsp.HTTPStatusCode,
		Platform:       Google,
	}

	return rspNotification, nil
}

type UrlNotification struct {
	// NotifyTime: Creation timestamp for this notification. Users should
	// _not_ specify it, the field is ignored at the request time.
	NotifyTime string `json:"notifyTime,omitempty"`

	// Type: The URL life cycle event that Google is being notified about.
	//
	// Possible values:
	//   "URL_NOTIFICATION_TYPE_UNSPECIFIED" - Unspecified.
	//   "URL_UPDATED" - The given URL (Web document) has been updated.
	//   "URL_DELETED" - The given URL (Web document) has been deleted.
	Type string `json:"type,omitempty"`

	// Url: The object of this notification. The URL must be owned by the
	// publisher of this notification and, in case of `URL_UPDATED`
	// notifications, it _must_ be crawlable by Google.
	Url string `json:"url,omitempty"`

	HTTPStatusCode int `json:"http_status_code"`

	Platform Platform `json:"platform"`
}

func (s *SearchIndex) bingUpdateIndex(uri string) (*UrlNotification, error) {
	tmpUrl, _ := url.Parse(uri)
	body := map[string]string{
		"siteUrl": fmt.Sprintf("%s://%s", tmpUrl.Scheme, tmpUrl.Host),
		"url":     uri,
	}
	rsp, err := s.restyClient.
		SetBaseURL(s.bingApiHost).
		R().
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{"apikey": s.bingApiKey}).
		SetBody(body).
		Post("/webmaster/api.svc/json/SubmitUrl")
	if err != nil {
		return nil, err
	}
	tmp := &UrlNotification{
		HTTPStatusCode: rsp.StatusCode(),
		Url:            uri,
		Type:           "URL_UPDATED",
		NotifyTime:     time.Now().String(),
		Platform:       Bing,
	}
	return tmp, nil
}
