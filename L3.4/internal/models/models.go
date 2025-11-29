package models

type Task struct {
	TaskID      string     `json:"id"`
	ImageID     string     `json:"image_id"`
	Type        string     `json:"type"`
	Params      TaskParams `json:"params"`
	ResultPath  string     `json:"result_path"`
	Status      TaskStatus `json:"status"`
	Extension   string     `json:"extension"`
	ContentType string     `json:"content_type"`
}

type TaskStatus string

const (
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusDone       TaskStatus = "done"
	TaskStatusFailed     TaskStatus = "failed"
)

type TaskParams struct {
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Watermark string `json:"text,omitempty"`
}
