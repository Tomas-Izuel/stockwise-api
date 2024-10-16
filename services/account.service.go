package services

import (
	"context"
	"invest/errors"
	"invest/lib"
	"invest/models"
	"invest/models/dto"
	"invest/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllAcountsByUserID(ctx context.Context, userID string) ([]models.Account, error) {
	accounts, err := repository.GetAllAccountsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func CreateAccount(ctx context.Context, id string, accountDTO *dto.CreateAccountDTO) (*mongo.InsertOneResult, error) {
	if err := validate.Struct(accountDTO); err != nil {
		return nil, errors.NewValidationError(lib.MapValidationErrors(err))
	}

	userExist, err := repository.FindUserByID(ctx, id)
	if err != nil || userExist == nil {
		return nil, errors.Wrap(404, "user not found", err)
	}

	typeExists, err := repository.FindAccountTypeByID(ctx, accountDTO.TypeID)
	if err != nil || typeExists == nil {
		return nil, errors.Wrap(404, "account type not found", err)
	}

	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Wrap(400, "invalid user ID format", err)
	}

	typeID, err := primitive.ObjectIDFromHex(accountDTO.TypeID)
	if err != nil {
		return nil, errors.Wrap(400, "invalid account type ID format", err)
	}

	account := &models.Account{
		UserID: userID,
		Period: accountDTO.Period,
		Type:   typeID,
	}

	user, err := repository.InsertAccount(ctx, account)
	if err != nil {
		return nil, errors.Wrap(500, "failed to create account", err)
	}

	return user, nil
}

func UpdateAccount(ctx context.Context, id string, accountDTO *dto.UpdateAccountDTO) (*mongo.UpdateResult, error) {
	if err := validate.Struct(accountDTO); err != nil {
		return nil, errors.NewValidationError(lib.MapValidationErrors(err))
	}

	accountExist, err := repository.FindAccountByID(ctx, id)
	if err != nil || accountExist == nil {
		return nil, errors.Wrap(404, "account not found", err)
	}

	updateData := map[string]interface{}{
		"period": accountDTO.Period,
	}

	userUpdated, err := repository.UpdateAccount(ctx, id, updateData)
	if err != nil {
		return nil, errors.Wrap(500, "failed to update account", err)
	}

	return userUpdated, nil
}

func DeleteAccount(ctx context.Context, id string) error {
	accountExist, err := repository.FindAccountByID(ctx, id)
	if err != nil || accountExist == nil {
		return errors.Wrap(404, "account not found", err)
	}

	if err := repository.DeleteAccount(ctx, id); err != nil {
		return errors.Wrap(500, "failed to delete account", err)
	}

	return nil
}

func GetAccountByID(ctx context.Context, id string) (*models.Account, error) {
	account, err := repository.FindAccountByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return account, nil
}
