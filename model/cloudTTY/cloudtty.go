package cloudTTY

type CloudTTY struct{}

type PodMessage struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Container string `json:"container"`
}
