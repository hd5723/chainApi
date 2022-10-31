package contract

import (
	"OnchainParser/internal/dao"
	"OnchainParser/internal/model"
	"OnchainParser/internal/model/entity"
	"OnchainParser/internal/service"
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

type (
	sProtocol struct{}
)

func init() {
	service.RegisterProtocol(New())
}

func New() *sProtocol {
	return &sProtocol{}
}

func (s *sProtocol) QueryOneByProtocolCode(ctx context.Context, protocolCode string, chainName string) (rcd entity.ProtocolEntity, err error) {
	return dao.Protocol.QueryOneByProtocolCode(ctx, protocolCode, chainName)
}

func (s *sProtocol) ToAuditList(ctx context.Context, chainName string) (result []entity.ProtocolEntity, err error) {
	return dao.Protocol.ToAuditList(ctx, chainName)
}

func (s *sProtocol) DeleteByProtocolCode(ctx context.Context, protocolCode string, chainName string) (err error) {
	sql := " drop database IF EXISTS " + protocolCode
	if _, err = g.DB().Exec(ctx, sql); err != nil {
		return err
	}
	return dao.Protocol.DeleteByProtocolCode(ctx, protocolCode, chainName)
}

func (s *sProtocol) QueryEnableList(ctx context.Context, chainName string) (result []entity.ProtocolEntity, err error) {
	return dao.Protocol.QueryEnableList(ctx, chainName)
}

func (s *sProtocol) QueryList(ctx context.Context, protocol_code string, is_valid string, chainName string) (result []entity.ProtocolEntity, err error) {
	return dao.Protocol.QueryList(ctx, protocol_code, is_valid, chainName)
}

func (s *sProtocol) ExecInsert(ctx context.Context, chainName string, in model.ProtocolCreateInput) (err error) {
	return dao.Protocol.ExecInsert(ctx, chainName, in)
}

func (s *sProtocol) UpdateValid(ctx context.Context, protocolCode string, isVaild bool, chainName string) (err error) {
	return dao.Protocol.UpdateValid(ctx, protocolCode, isVaild, chainName)
}
