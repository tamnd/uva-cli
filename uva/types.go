package uva

// Problem is a UVa Online Judge problem.
type Problem struct {
	Rank        int    `json:"rank"`
	PID         int    `json:"pid"`
	Num         int    `json:"num"`
	Title       string `json:"title"`
	AC          int    `json:"ac"`
	WA          int    `json:"wa"`
	TLE         int    `json:"tle"`
	DACU        int    `json:"dacu"`
	TimeLimitMS int    `json:"time_limit_ms"`
	URL         string `json:"url"`
}
