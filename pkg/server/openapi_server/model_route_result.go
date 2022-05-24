// SPDX-License-Identifier: MIT

package openapi_server

type RouteResult struct {

	Origin Point `json:"origin"`

	Destination Point `json:"destination"`

	// States whether a route from origin to destination exists
	Reachable bool `json:"reachable"`

	Path Path `json:"path,omitempty"`
}

// AssertRouteResultRequired checks if the required fields are not zero-ed
func AssertRouteResultRequired(obj RouteResult) error {
	elements := map[string]interface{}{
		"origin": obj.Origin,
		"destination": obj.Destination,
		"reachable": obj.Reachable,
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
	if err := AssertPathRequired(obj.Path); err != nil {
		return err
	}
	return nil
}

// AssertRecurseRouteResultRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of RouteResult (e.g. [][]RouteResult), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseRouteResultRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aRouteResult, ok := obj.(RouteResult)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertRouteResultRequired(aRouteResult)
	})
}
