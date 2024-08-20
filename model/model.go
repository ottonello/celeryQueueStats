package model

// Define the structs to match the JSON structure
type DeliveryInfo struct {
	Exchange   string `json:"exchange"`
	RoutingKey string `json:"routing_key"`
}

type Properties struct {
	CorrelationID string       `json:"correlation_id"`
	ReplyTo       string       `json:"reply_to"`
	DeliveryMode  int          `json:"delivery_mode"`
	DeliveryInfo  DeliveryInfo `json:"delivery_info"`
	Priority      int          `json:"priority"`
	BodyEncoding  string       `json:"body_encoding"`
	DeliveryTag   string       `json:"delivery_tag"`
}

type Headers struct {
	Lang            string         `json:"lang"`
	Task            string         `json:"task"`
	ID              string         `json:"id"`
	Shadow          interface{}    `json:"shadow"`
	Eta             string         `json:"eta"`
	Expires         interface{}    `json:"expires"`
	Group           interface{}    `json:"group"`
	GroupIndex      interface{}    `json:"group_index"`
	Retries         int            `json:"retries"`
	TimeLimit       [2]interface{} `json:"timelimit"`
	RootID          string         `json:"root_id"`
	ParentID        interface{}    `json:"parent_id"`
	ArgsRepr        string         `json:"argsrepr"`
	KwargsRepr      string         `json:"kwargsrepr"`
	Origin          string         `json:"origin"`
	IgnoreResult    bool           `json:"ignore_result"`
	CeleryAppliedAt string         `json:"celery_applied_at"`
}

type Message struct {
	Body            string     `json:"body"`
	ContentEncoding string     `json:"content-encoding"`
	ContentType     string     `json:"content-type"`
	Headers         Headers    `json:"headers"`
	Properties      Properties `json:"properties"`
}
