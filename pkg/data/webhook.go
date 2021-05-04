package data

type Webhook struct {
	Event       string      `json:"event"`
	Timestamp   string      `json:"timestamp"`
	Model       string      `json:"model"`
	Username    string      `json:"username"`
	RequestId   string      `json:"request_id"`
	Data        WebhookData `json:"data"`
	Created     string      `json:"created"`
	LastUpdated string      `json:"last_updated"`
}

type WebhookData struct {
	ID             string             `json:"id"`
	Url            string             `json:"url"`
	Address        string             `json:"address"`
	AssignedObject AssignedObjectType `json:"assigned_object"`
}

type AssignedObjectType struct {
	ID               string               `json:"id"`
	AssignedObjectId string               `json:"assigned_object_id"`
	AssignedObject   AssignedObjectDevice `json:"assigned_object"`
	Name             string               `json:"name"`
}

type AssignedObjectDevice struct {
	ID          string `json:"id"`
	Url         string `json:"url"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}
