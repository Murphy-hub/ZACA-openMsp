package dao

import (
	"context"
	"time"

	"gitlab.oneitfarm.com/bifrost/capitalizone/examples/cfssl-model/model"

	"github.com/guregu/null"
	"github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = null.Bool{}
	_ = uuid.UUID{}
)

// GetAllOcspResponses is a function to get a slice of record(s) from ocsp_responses table in the cap database
// params - page     - page requested (defaults to 0)
// params - pagesize - number of records in a page  (defaults to 20)
// params - order    - db sort order column
// error - ErrNotFound, db Find error
func GetAllOcspResponses(ctx context.Context, page, pagesize int64, order string) (results []*model.OcspResponses, totalRows int, err error) {

	resultOrm := DB.Model(&model.OcspResponses{})
	resultOrm.Count(&totalRows)

	if page > 0 {
		offset := (page - 1) * pagesize
		resultOrm = resultOrm.Offset(offset).Limit(pagesize)
	} else {
		resultOrm = resultOrm.Limit(pagesize)
	}

	if order != "" {
		resultOrm = resultOrm.Order(order)
	}

	if err = resultOrm.Find(&results).Error; err != nil {
		err = ErrNotFound
		return nil, -1, err
	}

	return results, totalRows, nil
}

// GetOcspResponses is a function to get a single record from the ocsp_responses table in the cap database
// error - ErrNotFound, db Find error
func GetOcspResponses(ctx context.Context, argSerialNumber string, argAuthorityKeyIdentifier string) (record *model.OcspResponses, err error) {
	record = &model.OcspResponses{}
	if err = DB.First(record, argSerialNumber, argAuthorityKeyIdentifier).Error; err != nil {
		err = ErrNotFound
		return record, err
	}

	return record, nil
}

// AddOcspResponses is a function to add a single record to ocsp_responses table in the cap database
// error - ErrInsertFailed, db save call failed
func AddOcspResponses(ctx context.Context, record *model.OcspResponses) (result *model.OcspResponses, RowsAffected int64, err error) {
	db := DB.Save(record)
	if err = db.Error; err != nil {
		return nil, -1, ErrInsertFailed
	}

	return record, db.RowsAffected, nil
}

// UpdateOcspResponses is a function to update a single record from ocsp_responses table in the cap database
// error - ErrNotFound, db record for id not found
// error - ErrUpdateFailed, db meta data copy failed or db.Save call failed
func UpdateOcspResponses(ctx context.Context, argSerialNumber string, argAuthorityKeyIdentifier string, updated *model.OcspResponses) (result *model.OcspResponses, RowsAffected int64, err error) {

	result = &model.OcspResponses{}
	db := DB.First(result, argSerialNumber, argAuthorityKeyIdentifier)
	if err = db.Error; err != nil {
		return nil, -1, ErrNotFound
	}

	if err = Copy(result, updated); err != nil {
		return nil, -1, ErrUpdateFailed
	}

	db = db.Save(result)
	if err = db.Error; err != nil {
		return nil, -1, ErrUpdateFailed
	}

	return result, db.RowsAffected, nil
}

// DeleteOcspResponses is a function to delete a single record from ocsp_responses table in the cap database
// error - ErrNotFound, db Find error
// error - ErrDeleteFailed, db Delete failed error
func DeleteOcspResponses(ctx context.Context, argSerialNumber string, argAuthorityKeyIdentifier string) (rowsAffected int64, err error) {

	record := &model.OcspResponses{}
	db := DB.First(record, argSerialNumber, argAuthorityKeyIdentifier)
	if db.Error != nil {
		return -1, ErrNotFound
	}

	db = db.Delete(record)
	if err = db.Error; err != nil {
		return -1, ErrDeleteFailed
	}

	return db.RowsAffected, nil
}