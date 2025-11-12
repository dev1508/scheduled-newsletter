package response

// DeliveryStats represents delivery statistics for content
type DeliveryStats struct {
	ContentID    string `json:"content_id"`
	TotalSent    int    `json:"total_sent"`
	Pending      int    `json:"pending"`
	Delivered    int    `json:"delivered"`
	Failed       int    `json:"failed"`
	Bounced      int    `json:"bounced"`
	DeliveryRate float64 `json:"delivery_rate"`
}
