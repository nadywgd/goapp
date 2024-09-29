package watcher

type Counter struct {
	Iteration int `json:"iteration"`
}

type CounterReset struct {
	Action string `json:"action"`
	Value  int    `json:"value"`
}
