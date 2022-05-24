// SPDX-License-Identifier: MIT

package openapi_server

// RouteRequest - Request a route from a origin to a destination. 
type RouteRequest struct {

	Origin Point `json:"origin"`

	Destination Point `json:"destination"`
}

// AssertRouteRequestRequired checks if the required fields are not zero-ed
func AssertRouteRequestRequired(obj RouteRequest) error {
	elements := map[string]interface{}{
		"origin": obj.Origin,
		"destination": obj.Destination,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	if err := AssertPointRequired(obj.Origin); err != nil {
		return err
	}
	if err := AssertPointRequired(obj.Destination); err != nil {
		return err
	}
	return nil
}

// AssertRecurseRouteRequestRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of RouteRequest (e.g. [][]RouteRequest), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseRouteRequestRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aRouteRequest, ok := obj.(RouteRequest)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertRouteRequestRequired(aRouteRequest)
	})
}
