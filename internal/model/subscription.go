package model

import (
	"github.com/google/uuid"
)

type Subscription struct{
	ServiceName string 	`json:"service_name"`
	Price int 				`json:"price"`
	Id uuid.UUID			`json:"user_id"`
	DataStart string		`json:"start_date"`
	DataEnd *string			`json:"end_date,omitempty"`


}