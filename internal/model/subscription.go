package model

import (
	"github.com/google/uuid"
)

type Subscription struct{
	ID int 				`json:"id"`
	ServiceName string 	`json:"service_name"`
	Price int 				`json:"price"`
	UserId uuid.UUID			`json:"user_id"`
	StartData string		`json:"start_data"`
	EndData *string			`json:"end_data,omitempty"`


}