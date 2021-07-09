package subscriptions

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type subscriptionResult struct {
	gophercloud.Result
}

func (r subscriptionResult) Extract() (*Subscription, error) {
	var s Subscription
	err := r.ExtractInto(&s)
	return &s, err
}

func (r subscriptionResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

func ExtractSubscriptionsInto(r pagination.Page, v interface{}) error {
	return r.(SubscriptionPage).Result.ExtractIntoSlicePtr(v, "subscriptions")
}

// Subscription represents a subscription in the OpenStack Bare Metal API.
type Subscription struct {
	// UUID for the resource.
	UUID string `json:"uuid"`

	// UUID of the Node this resource belongs to.
	NodeUUID string `json:"node_uuid"`

	Destination    string
	Context        string
	SubscriptionID string
}

// SubscriptionPage abstracts the raw results of making a List() request against
// the API.
type SubscriptionPage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if a page contains no Subscription results.
func (r SubscriptionPage) IsEmpty() (bool, error) {
	s, err := ExtractSubscriptions(r)
	return len(s) == 0, err
}

// NextPageURL uses the response's embedded link reference to navigate to the
// next page of results.
func (r SubscriptionPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"subscriptions_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// ExtractSubscriptions interprets the results of a single page from a List() call,
// producing a slice of Subscription entities.
func ExtractSubscriptions(r pagination.Page) ([]Subscription, error) {
	var s []Subscription
	err := ExtractSubscriptionsInto(r, &s)
	return s, err
}

// GetResult is the response from a Get operation. Call its Extract
// method to interpret it as a Subscription.
type GetResult struct {
	subscriptionResult
}

// CreateResult is the response from a Create operation.
type CreateResult struct {
	subscriptionResult
}

// UpdateResult is the response from an Update operation. Call its Extract
// method to interpret it as a Subscription.
type UpdateResult struct {
	subscriptionResult
}

// DeleteResult is the response from a Delete operation. Call its ExtractErr
// method to determine if the call succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}
