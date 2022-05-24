// SPDX-License-Identifier: MIT

package openapi_server

// Path - A path is described by sequence of points as well as its total length.
type Path struct {

	// A path is an ordered list of points.
	Waypoints []Point `json:"waypoints"`

	// unit meters
	Length int32 `json:"length"`
}

// AssertPathRequired checks if the required fields are not zero-ed
func AssertPathRequired(obj Path) error {
	elements := map[string]interface{}{
		"waypoints": obj.Waypoints,
		"length": obj.Length,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.Waypoints {
		if err := AssertPointRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertRecursePathRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Path (e.g. [][]Path), otherwise ErrTypeAssertionError is thrown.
func AssertRecursePathRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aPath, ok := obj.(Path)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertPathRequired(aPath)
	})
}
