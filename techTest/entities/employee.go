package entities

import (
	"errors"
	"strings"

	"github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Employee struct {
	ID          primitive.ObjectID `bson:"_id, omitempty"`
	Name        string             `bson:"name" valid:"required, length(2|50)"`
	Email       string             `bson:"email" valid:"required, email, length(5|255)"`
	PhoneNumber interface{}        `bson:"phone_number, omitempty" valid:"length(8|16)"`
	JobTitle    string             `bson:"jobtitle" valid:"required"`
	CompanyID   primitive.ObjectID `bson:"company_id" valid:"required"`
}

var JobTitleEnum = []string{"manager", "director", "staff"}

func (e *Employee) BeforeCreate() (err error) {

	e.JobTitle = strings.ToLower(e.JobTitle)
	err = errors.New("job title invalid")
	for _, value := range JobTitleEnum {
		if e.JobTitle == value {
			err = nil
		}
	}
	if err != nil {
		return
	}

	if e.PhoneNumber == "" {
		e.PhoneNumber = nil
	}

	_, err = govalidator.ValidateStruct(e)
	if err != nil {
		return
	}

	return nil
}
