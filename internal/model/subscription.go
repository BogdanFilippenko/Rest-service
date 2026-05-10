package model

import (
	"github.com/google/uuid"
)

type Subscription struct{
	ID int 				`json:"id"`
	ServiceName string 	`json:"service_name"`
	Price int 				`json:"price"`
	UserId uuid.UUID			`json:"user_id"`
	StartDate string		`json:"start_date"`
	EndDate *string			`json:"end_date,omitempty"`


}