package task

type Message struct {
	TaskID  string         `json:"task_id"`
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}
