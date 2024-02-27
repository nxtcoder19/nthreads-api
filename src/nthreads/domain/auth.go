package domain

import (
	"context"
	"fmt"
	"github.com/nxtcoder19/nthreads-backend/package/functions"
	mongo_db "github.com/nxtcoder19/nthreads-backend/package/mongo-db"
	"github.com/nxtcoder19/nthreads-backend/src/nthreads/entities"
)

func (i *Impl) SignUp(ctx context.Context, firstName string, lastName string, email string, password string) (*entities.User, error) {
	id := i.db.NewId()
	verifiedPasswordString := "verify" + password
	encodedString := functions.GetMD5Hash(verifiedPasswordString)
	user := entities.User{
		Id:             id,
		FirstName:      firstName,
		LastName:       lastName,
		Email:          email,
		Password:       password,
		VerifyPassword: encodedString,
	}
	_, err := i.db.InsertRecord(ctx, AuthTable, user)
	//fmt.Println("user", nUser)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (i *Impl) Login(ctx context.Context, email string, password string) (res bool, err error) {
	verifyPassword := "verify" + password
	encodedString := functions.GetMD5Hash(verifyPassword)

	var user entities.User
	err = i.db.FindOne(ctx, AuthTable, &user, mongo_db.Filter{"email": email})
	if err != nil {
		return false, err
	}

	//fmt.Println("user is", user)
	if user.VerifyPassword != encodedString {
		return false, err
	}
	return true, nil
}

func (i *Impl) UpdateUser(ctx context.Context, email string, firstName string, lastName string) (*entities.User, error) {
	err := i.db.UpdateMany(
		ctx,
		AuthTable,
		mongo_db.Filter{"email": email},
		mongo_db.Filter{
			"firstname": firstName,
			"lastname":  lastName,
		},
	)

	var user entities.User
	err = i.db.FindOne(ctx, AuthTable, &user, mongo_db.Filter{"email": email})
	if err != nil {
		return nil, err
	}

	fmt.Println("user", user)
	return &user, nil
}

func (i *Impl) DeleteUser(ctx context.Context, email string) (string, error) {
	var user entities.User
	err := i.db.FindOne(ctx, AuthTable, &user, mongo_db.Filter{"email": email})
	if err != nil {
		return "", err
	}

	err = i.db.DeleteRecord(ctx, AuthTable, mongo_db.Filter{"email": email})
	if err != nil {
		return "", err
	}
	return "user deleted successfully", nil
}
