package search_index_group

import (
	"context"
	"google.golang.org/api/indexing/v3"
	"google.golang.org/api/option"
)

func NewSearchIndex() (*SearchIndex, error) {
	ctx := context.Background()
	service, err := indexing.NewService(ctx, option.WithCredentialsJSON([]byte("")))
	if err != nil {
		return nil, err
	}
	return &SearchIndex{
		googleSearchIndex: service,
	}, nil
}

type SearchIndex struct {
	googleSearchIndex *indexing.Service
}

func (s *SearchIndex) PublishUrl(url string) (*UrlNotification, error) {
	return nil, nil
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

	// ForceSendFields is a list of field names (e.g. "NotifyTime") to
	// unconditionally include in API requests. By default, fields with
	// empty or default values are omitted from API requests. However, any
	// non-pointer, non-interface field appearing in ForceSendFields will be
	// sent to the server regardless of whether the field is empty or not.
	// This may be used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "NotifyTime") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}
