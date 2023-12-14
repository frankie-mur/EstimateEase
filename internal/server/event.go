package server

// Represent an HTMX websocket event
type Event struct {
	Header  HTMXHeader `json:"HEADERS"`
	Payload string     `json:"payload"`
}

type HTMXHeader struct {
	HXRequest     string `json:"HX-Request"`
	HXTrigger     string `json:"HX-Trigger"`
	HXTriggerName string `json:"HX-Trigger-Name"`
	HXTarget      string `json:"HX-Target"`
	HXCurrentURL  string `json:"HX-Current-URL"`
}
