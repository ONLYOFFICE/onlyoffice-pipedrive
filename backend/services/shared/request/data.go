package request

type DataRequest struct {
	UserID int    `json:"user_id"`
	FileID string `json:"file_id"`
	DealID string `json:"deal_id"`
}
