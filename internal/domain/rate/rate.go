package rate

type Rate struct {
	Date   string  `json:"date"`
	Base   string  `json:"base"`
	Target string  `json:"target"`
	Value  float64 `json:"value"`
}
