package entities

import (
	"github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Company struct {
	ID              primitive.ObjectID `bson:"_id, omitempty"`
	CompanyName     string             `bson:"company_name" valid:"required, length(3|50)"`
	TelephoneNumber interface{}        `bson:"telephone_number, omitempty" valid:"length(8|16)"`
	IsActive        bool               `bson:"is_active"`
	Address         string             `bson:"address" valid:"length(10|50)"`
}

func (c *Company) BeforeCreate() (err error) {
	if c.TelephoneNumber == "" {
		c.TelephoneNumber = nil
	}

	_, err = govalidator.ValidateStruct(c)
	if err != nil {
		return
	}

	return nil
}
