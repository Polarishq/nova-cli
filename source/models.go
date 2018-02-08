package source

type StrMatrix [][]string

type MetricsLSResponse struct {
	Metrics []string `json:"metrics"`
}

type NovaSearchStatsQuery struct {
	Fields []string `json:"fields"`
	FieldsType string `json:"fields_type"`
	GroupBy []string `json:"group_by"`
	Stats []string `json:"stats"`
	Blocking bool `json:"blocking"`
	Span string `json:"span"`
}

type NovaSearchEventsQuery struct {
	Blocking bool `json:"blocking"`
	Mode string `json:"mode"`
	SearchTerms string `json:"search_terms"`
	Transforms []string `json:"transforms"`
	Reports []string `json:"reports"`
}

type NovaSearchResponse struct {
	Errors []string `json:"errors"`
	Items []interface{} `json:"items"`
	Job struct {
		SearchID string `json:"search_id"`
		SearchClass string `json:"search_class"`
		Status string `json:"status"`
	} `json:"job"`
	Metadata struct {
		Count int `json:"count"`
		Offset int `json:"offset"`
		TotalCount int `json:"total_count"`
	} `json:"metadata"`
}

type NovaOutgoingEventFormat struct {
	Source string            `json:"source"`
	Entity string            `json:"entity"`
	Event  map[string]string `json:"event"`
}

type NovaIncomingEventReporting []struct {
	Payload interface{} `json:"event"`
	Source string `json:"source"`
	Time string `json:"time"`
}

type NovaIncomingEventNonReporting []struct {
	Payload struct {
		Event struct {
			Raw string `json:"raw"`
		} `json:"event"`
	} `json:"event"`
	Source string `json:"source"`
	Time string `json:"time"`
}
