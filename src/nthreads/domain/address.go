package domain

import (
	"context"
	"fmt"
	mongo_db "github.com/nxtcoder19/nthreads-backend/package/mongo-db"
	"github.com/nxtcoder19/nthreads-backend/src/nthreads/entities"
)

func (i *Impl) CreateAddress(ctx context.Context, addressIn *entities.Address) (*entities.Address, error) {
	id := i.db.NewId()
	session := getSessionData(ctx)
	nAddress := &entities.Address{
		Id:               id,
		Email:            session.UserEmail,
		Name:             addressIn.Name,
		MobileNo:         addressIn.MobileNo,
		PinCode:          addressIn.PinCode,
		Locality:         addressIn.Locality,
		Area:             addressIn.Area,
		City:             addressIn.City,
		State:            addressIn.State,
		Landmark:         addressIn.Landmark,
		AlternatePhoneNo: addressIn.AlternatePhoneNo,
		AddressType:      addressIn.AddressType,
	}

	if addressIn.AddressType == "Home" {
		nAddress.AddressType = entities.HomeAddressType
	} else {
		nAddress.AddressType = entities.WorkAddressType
	}

	_, err := i.db.InsertRecord(ctx, AddressTable, nAddress)
	if err != nil {
		return nil, err
	}
	return nAddress, nil
}

func (i *Impl) UpdateAddress(ctx context.Context, id string, addressIn *entities.Address) (*entities.Address, error) {
	err := i.db.UpdateMany(
		ctx,
		AddressTable,
		mongo_db.Filter{"id": id},
		mongo_db.Filter{
			"name":     addressIn.Name,
			"mobileNo": addressIn.MobileNo,
			"pinCode":  addressIn.PinCode,
			"locality": addressIn.Locality,
			"area":     addressIn.Area,
		},
	)

	var address entities.Address
	err = i.db.FindOne(ctx, AddressTable, &address, mongo_db.Filter{"id": id})
	if err != nil {
		return nil, err
	}

	return &address, nil
}

func (i *Impl) GetAddress(ctx context.Context, id string) (*entities.Address, error) {
	var address entities.Address
	err := i.db.FindOne(ctx, AddressTable, &address, mongo_db.Filter{"id": id})
	if err != nil {
		return nil, err
	}

	return &address, nil
}

func (i *Impl) GetAddresses(ctx context.Context) ([]*entities.Address, error) {
	addresses := make([]*entities.Address, 0)
	session := getSessionData(ctx)
	cursor, err := i.db.Find(ctx, AddressTable, mongo_db.Filter{"email": session.UserEmail})
	if err != nil {
		return nil, err
	}
	defer func() {
		if cer := cursor.Close(ctx); cer != nil {
			fmt.Println(cer)
		}
	}()
	for cursor.Next(ctx) {
		var address entities.Address
		if err := cursor.Decode(&address); err != nil {
			return nil, err
		}
		addresses = append(addresses, &address)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return addresses, nil
}

func (i *Impl) DeleteAddress(ctx context.Context, id string) (bool, error) {
	var address entities.Address
	err := i.db.FindOne(ctx, AddressTable, &address, mongo_db.Filter{"id": id})
	if err != nil {
		return false, err
	}

	err = i.db.DeleteRecord(ctx, AddressTable, mongo_db.Filter{"id": id})
	if err != nil {
		return false, err
	}
	return true, nil
}
