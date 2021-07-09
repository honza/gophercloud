package subscriptions

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToSubscriptionListQuery() (string, error)
	ToSubscriptionListDetailQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the node attributes you want to see returned. Marker and Limit are used
// for pagination.
type ListOpts struct {
	// Filter the list by the name or uuid of the Node
	Node string `q:"node"`

	// Filter the list by the Node uuid
	NodeUUID string `q:"node_uuid"`

	// One or more fields to be returned in the response.
	Fields []string `q:"fields"`

	// Requests a page size of items.
	Limit int `q:"limit"`

	// The ID of the last-seen item
	Marker string `q:"marker"`

	// Sorts the response by the requested sort direction.
	// Valid value is asc (ascending) or desc (descending). Default is asc.
	SortDir string `q:"sort_dir"`

	// Sorts the response by the this attribute value. Default is id.
	SortKey string `q:"sort_key"`
}

// ToSubscriptionListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToSubscriptionListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request against the API to list subscriptions accessible to you.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToSubscriptionListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return SubscriptionPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ToSubscriptionListDetailQuery formats a ListOpts into a query string for the list details API.
func (opts ListOpts) ToSubscriptionListDetailQuery() (string, error) {
	// Detail endpoint can't filter by Fields
	if len(opts.Fields) > 0 {
		return "", fmt.Errorf("fields is not a valid option when getting a detailed listing of subscriptions")
	}

	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListDetail - Return a list subscriptions with complete details.
// Some filtering is possible by passing in flags in "ListOpts",
// but you cannot limit by the fields returned.
func ListDetail(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listDetailURL(client)
	if opts != nil {
		query, err := opts.ToSubscriptionListDetailQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return SubscriptionPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get - requests the details off a subscription, by ID.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(getURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToSubscriptionCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies subscription creation parameters.
type CreateOpts struct {
	// UUID of the Node this resource belongs to.
	NodeUUID string `json:"node_uuid,omitempty"`
}

// ToSubscriptionCreateMap assembles a request body based on the contents of a CreateOpts.
func (opts CreateOpts) ToSubscriptionCreateMap() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Create - requests the creation of a subscription
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToSubscriptionCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(createURL(client), reqBody, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// TODO Update
type Patch interface {
	ToSubscriptionUpdateMap() map[string]interface{}
}

// UpdateOpts is a slice of Patches used to update a subscription
type UpdateOpts []Patch

type UpdateOp string

const (
	ReplaceOp UpdateOp = "replace"
	AddOp     UpdateOp = "add"
	RemoveOp  UpdateOp = "remove"
)

type UpdateOperation struct {
	Op    UpdateOp    `json:"op" required:"true"`
	Path  string      `json:"path" required:"true"`
	Value interface{} `json:"value,omitempty"`
}

func (opts UpdateOperation) ToSubscriptionUpdateMap() map[string]interface{} {
	return map[string]interface{}{
		"op":    opts.Op,
		"path":  opts.Path,
		"value": opts.Value,
	}
}

// Update - requests the update of a subscription
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOpts) (r UpdateResult) {
	body := make([]map[string]interface{}, len(opts))
	for i, patch := range opts {
		body[i] = patch.ToSubscriptionUpdateMap()
	}

	resp, err := client.Patch(updateURL(client, id), body, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete - requests the deletion of a subscription
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(deleteURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
