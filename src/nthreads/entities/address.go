package entities

import mongodb "github.com/nxtcoder19/nthreads-backend/package/mongo-db"

type Address struct {
	Id               mongodb.ID  `json:"id"`
	Email            string      `json:"email"`
	Name             string      `json:"name"`
	MobileNo         int         `json:"mobileNo" bson:"mobileNo"`
	PinCode          int         `json:"pinCode" bson:"pinCode"`
	Locality         string      `json:"locality"`
	Area             string      `json:"area"`
	City             string      `json:"city"`
	State            string      `json:"state"`
	Landmark         string      `json:"landmark,omitempty"`
	AlternatePhoneNo int         `json:"alternatePhoneNo,omitempty" bson:"alternatePhoneNo"`
	AddressType      AddressType `json:"addressType" bson:"addressType"`
}

type AddressType string

const (
	HomeAddressType AddressType = "Home"
	WorkAddressType AddressType = "Work"
)
