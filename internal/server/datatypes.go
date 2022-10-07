package server

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type RouteRequest struct {
	Origin      Point `json:"origin"`
	Destination Point `json:"destination"`
}

type RouteResponse struct {
	Exists      bool    `json:"exists"`
	Path        Path    `json:"path,omitempty"`
	Time        int64   `json:"time"`
	SearchSpace []Point `json:"search_space,omitempty"`
}

type Path struct {
	Waypoints []Point `json:"waypoints"`
	Length    int     `json:"length"`
}
