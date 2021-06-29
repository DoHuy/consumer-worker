package dto

import "time"

type OrderEvent struct {
	ID           string    `json:"id,ignoreempty"`
	PartnerCode  string    `json:"partnerCode,ignoreempty"`
	ShipperID    string    `json:"shipperId,ignoreempty"`
	ShipperName  string    `json:"shipperName,ignoreempty"`
	ShipperPhone string    `json:"shipperPhone,ignoreempty"`
	PostID       string    `json:"postId,ignoreempty"`
	State        int       `json:"state,ignoreempty"`
	Time         time.Time `json:"time,ignoreempty"`
	Description  string    `json:"description,ignoreempty"`
}
