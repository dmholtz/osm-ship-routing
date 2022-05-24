// SPDX-License-Identifier: MIT

package openapi_server

// Point - Object representation of a point in the Geographic Coordinate System (GCS). 
type Point struct {

	// unit degree
	Lat float32 `json:"lat"`

	// unit degree
	Lon float32 `json:"lon"`
}

// AssertPointRequired checks if the required fields are not zero-ed
func AssertPointRequired(obj Point) error {
	elements := map[string]interface{}{
		"lat": obj.Lat,
		"lon": obj.Lon,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecursePointRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Point (e.g. [][]Point), otherwise ErrTypeAssertionError is thrown.
func AssertRecursePointRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aPoint, ok := obj.(Point)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertPointRequired(aPoint)
	})
}
