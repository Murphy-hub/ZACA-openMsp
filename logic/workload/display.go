// Package workload 展示层
package workload

import (
	"math/big"
	"strings"

	"github.com/araddon/dateparse"

	"github.com/pkg/errors"
	"gitlab.oneitfarm.com/bifrost/capitalizone/database/mysql/cfssl-model/dao"
	"gitlab.oneitfarm.com/bifrost/capitalizone/database/mysql/cfssl-model/model"
	"gitlab.oneitfarm.com/bifrost/capitalizone/logic/schema"
	"gorm.io/gorm"
)

type CertListParams struct {
	// 查询条件
	CertSN         string
	Role, UniqueID string
	// 分页条件
	Page, PageSize                 int
	Status                         string
	Order                          string
	ExpiryStartTime, ExpiryEndTime string
}

type CertListResult struct {
	CertList []*schema.FullCert
	Total    int64
}

// CertList 获取证书列表
func (l *Logic) CertList(params *CertListParams) (*CertListResult, error) {
	query := l.db.Session(&gorm.Session{})
	if params.CertSN != "" {
		sn := new(big.Int)
		i, ok := sn.SetString(params.CertSN, 10)
		if !ok {
			// try hex
			i, ok = sn.SetString(params.CertSN, 16)
			if !ok {
				return nil, errors.New("sn invalid")
			}
		}
		query = query.Where("serial_number = ?", i.String())
	}
	if params.Role != "" {
		query = query.Where("ca_label = ?", strings.ToLower(params.Role))
	}
	if params.UniqueID != "" {
		query = query.Where("common_name = ?", params.UniqueID)
	}
	if params.Order == "" {
		params.Order = "issued_at desc"
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.ExpiryStartTime != "" {
		date, err := dateparse.ParseAny(params.ExpiryStartTime)
		if err != nil {
			return nil, errors.Wrap(err, "过期时间错误")
		}
		query = query.Where("expiry > ?", date)
	}
	if params.ExpiryEndTime != "" {
		date, err := dateparse.ParseAny(params.ExpiryEndTime)
		if err != nil {
			return nil, errors.Wrap(err, "过期时间错误")
		}
		query = query.Where("expiry < ?", date)
	}

	query = query.Select(
		"ca_label",
		"common_name",
		"issued_at",
		"serial_number",
		"authority_key_identifier",
		"status",
		"not_before",
		"expiry",
		"revoked_at",
		"pem",
	)

	list, total, err := dao.GetAllCertificates(query, params.Page, params.PageSize, params.Order)
	if err != nil {
		return nil, errors.Wrap(err, "数据库查询错误")
	}
	var result CertListResult
	result.CertList = make([]*schema.FullCert, 0, len(list))
	for _, row := range list {
		cert, err := schema.GetFullCertByModelCert(row)
		if err != nil {
			continue
		}
		result.CertList = append(result.CertList, cert)
	}
	result.Total = total
	return &result, nil
}

type CertDetailParams struct {
	SN  string
	AKI string
}

func (l *Logic) CertDetail(params *CertDetailParams) (*schema.FullCert, error) {
	db := l.db.Session(&gorm.Session{})
	row := &model.Certificates{}
	if err := db.Where(&model.Certificates{
		SerialNumber:           params.SN,
		AuthorityKeyIdentifier: params.AKI,
	}).First(&row).Error; err != nil {
		return nil, errors.Wrap(err, "数据库查询错误")
	}
	cert, err := schema.GetFullCertByModelCert(row)
	if err != nil {
		return nil, err
	}
	return cert, nil
}
