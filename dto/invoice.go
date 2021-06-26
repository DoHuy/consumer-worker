package dto

import "time"

type Invoice struct {
	ID   		string      `json:"id"`
	ShipperID	string		`json:"shipperId"`
	PostID		int64 		`json:"postId"`
	State		int			`json:"state"`
	Time		time.Time	`json:"time"`
	Description	string		`json:"description"`
}