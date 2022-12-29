package persistence

import (
	"signupin-api/internal/app/api/dto"
	"signupin-api/internal/pkg/user"

	"github.com/kamva/mgm/v3"

	"github.com/kkodecaffeine/go-common/core/database/mongo/errortype"
	"github.com/kkodecaffeine/go-common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepo struct {
	client *mongo.Client
	mapper entityMapper
}

var _ user.Repository = &userRepo{}

func (r *userRepo) SaveOne(model *user.User) (string, error) {
	coll := mgm.Coll(model)
	err := coll.Create(model)
	if err != nil {
		return "", errortype.ParseAndReturnDBError(err, coll.Name(), nil, nil, nil)
	}

	insertedID := utils.MapToStringID(model.ID)
	return insertedID, nil
}

func (r *userRepo) GetAuthNumber() (string, error) {
	found := &user.AuthNumber{}
	filter := bson.D{}

	err := mgm.Coll(found).FindOne(mgm.Ctx(), filter).Decode(found)
	if err != nil {
		return "", errortype.ParseAndReturnDBError(err, mgm.CollName(found), filter, nil, nil)
	}

	return found.AuthNumber, nil
}

func (r *userRepo) GetOne(identifier string, password ...string) (*dto.GetUserWithTokenResponse, error) {
	found := &user.User{}
	var filter primitive.M

	if len(password) == 0 {
		filter = bson.M{"email": identifier}
	} else {
		filter = bson.M{
			"$and": []bson.M{
				{"password": password[0]},
				{
					"$or": []bson.M{
						{"email": identifier},
						{"phone": identifier},
					}},
			},
		}
	}

	err := mgm.Coll(found).FindOne(mgm.Ctx(), filter).Decode(found)
	if err != nil {
		return nil, errortype.ParseAndReturnDBError(err, mgm.CollName(found), filter, nil, nil)
	}

	result := r.mapper.toDomainProps2(found.ID, found)

	return result, nil
}

func (r *userRepo) GetOneByID(ID string) (*dto.GetUserResponse, error) {
	objectID, err := utils.MapToObjectID(ID)
	if err != nil {
		return nil, err
	}

	found := &user.User{}
	filter := bson.M{"_id": objectID}

	coll := mgm.Coll(found)
	err = coll.FindOne(mgm.Ctx(), filter).Decode(found)
	if err != nil {
		return nil, errortype.ParseAndReturnDBError(err, coll.Name(), filter, nil, nil)
	}

	result := r.mapper.toDomainProps(found.ID, found)

	return result, nil
}

func (r *userRepo) UpdatePassword(ID primitive.ObjectID, newpassword string) (*dto.GetUserResponse, error) {
	found := &user.User{}
	filter := bson.D{{Key: "_id", Value: ID}}

	coll := mgm.Coll(found)
	coll.FindOne(mgm.Ctx(), filter).Decode(&found)

	found.Password = newpassword

	err := coll.Update(found)
	if err != nil {
		return nil, errortype.ParseAndReturnDBError(err, coll.Name(), nil, nil, nil)
	}

	result := r.mapper.toDomainProps(found.ID, found)

	return result, nil
}

func (r *userRepo) UpsertAuthNumber(model *user.AuthNumber) (string, error) {
	found := &user.AuthNumber{}
	filter := bson.D{}

	coll := mgm.Coll(found)
	coll.FindOne(mgm.Ctx(), filter).Decode(&found)

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	err := coll.Update(model, &opt)
	if err != nil {
		return "", errortype.ParseAndReturnDBError(err, coll.Name(), nil, nil, nil)
	}

	return model.AuthNumber, nil
}

func New(client *mongo.Client) user.Repository {
	return &userRepo{client, entityMapper{}}
}
