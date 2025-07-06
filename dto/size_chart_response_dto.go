package dto

type Component struct {
	Table [][]string `json:"table"`
}

type Section struct {
	Name       string     `json:"name"`
	Components *Component `json:"component"`
}

type SizeChartDetailsResponse struct {
	Data []*Section `json:"sections"`
}
