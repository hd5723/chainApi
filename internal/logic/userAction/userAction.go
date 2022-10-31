package contract

import (
	"OnchainParser/internal/service"
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type (
	sUserAction struct{}
)

func init() {
	service.RegisterUserAction(New())
}

func New() *sUserAction {
	return &sUserAction{}
}

func (s *sUserAction) QueryTablesByProtocolCodeAndContractCode(ctx context.Context, protocolCode string, contractCode string) (result []gdb.Value, err error) {
	table := "information_schema.tables"
	queryTableName := protocolCode + "_" + contractCode + "%"
	result, err = g.DB().Model(table).Safe().Ctx(ctx).Where("table_schema=?", protocolCode).WhereLike("table_name", queryTableName).Limit(0, 1000).Array("CONCAT( 'drop table IF EXISTS ',  CONCAT(table_schema,'.',table_name) , ';' ) as tablename")
	if err != nil {
		g.Log().Debug(ctx, "QueryTablesByProtocolCodeAndContractCode error:", err)
		return nil, err
	}
	return result, nil
}

func (s *sUserAction) Excute(ctx context.Context, sqls []gdb.Value) (err error) {
	for i := 0; i < len(sqls); i++ {
		sql := sqls[i].String()
		if _, err = g.DB().Exec(ctx, sql); err != nil {
			return err
		}
	}
	return nil
}
