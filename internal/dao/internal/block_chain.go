// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"OnchainParser/internal/chain/ethEvm/consts"
	"OnchainParser/internal/model"
	"OnchainParser/internal/model/entity"
	"bytes"
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"strconv"
)

// BlockChainDao is the data access object for table block_chain.
type BlockChainDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns BlockChainColumns // columns contains all the column names of Table for convenient usage.
}

// BlockChainColumns defines and stores column names for table block_chain.
type BlockChainColumns struct {
	ChainName  string // 区块链名
	BaseCoin   string // e.g BNB
	ChainId    string // 链编号
	RpcUrls    string // RPC配置
	UpdateTime string // 修改时间(时间戳)
}

// blockChainColumns holds the columns for table block_chain.
var blockChainColumns = BlockChainColumns{
	ChainName:  "chain_name",
	BaseCoin:   "base_coin",
	ChainId:    "chain_id",
	RpcUrls:    "rpc_urls",
	UpdateTime: "update_time",
}

// NewBlockChainDao creates and returns a new DAO object for table data access.
func NewBlockChainDao() *BlockChainDao {
	return &BlockChainDao{
		group:   "default",
		table:   "block_chain",
		columns: blockChainColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *BlockChainDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *BlockChainDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *BlockChainDao) Columns() BlockChainColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *BlockChainDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *BlockChainDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *BlockChainDao) CtxWithDatabase(ctx context.Context, database string) *gdb.Model {
	return dao.DB().Model(database + "." + dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *BlockChainDao) Transaction(ctx context.Context, f func(ctx context.Context, tx *gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

func (dao *BlockChainDao) QueryOneByChainId(ctx context.Context, chainId int32, chainName string) (rcd entity.BlockChainEntity) {
	key := chainName + "_" + consts.BLOCK_CHAIN_DATA_KEY
	v, err := g.Redis().Do(ctx, "GET", key)
	if err == nil {
		serr := v.Struct(&rcd)
		if serr == nil && !g.IsEmpty(rcd.ChainId) {
			return rcd
		}
	}

	temp, err := dao.CtxWithDatabase(ctx, chainName).Where("chain_id =?", chainId).One()
	if temp == nil {
		g.Log().Debug(ctx, "error:", err)
		return rcd
	}
	json.Unmarshal([]byte(temp.Json()), &rcd)

	//设置缓存过期时间10天
	g.Redis().Do(ctx, "SETEX", key, 3600*24, rcd)
	return rcd
}

func (dao *BlockChainDao) QueryOneByChainName(ctx context.Context, chainName string) (rcd entity.BlockChainEntity) {
	key := chainName + "_" + consts.BLOCK_CHAIN_DATA_KEY
	v, err := g.Redis().Do(ctx, "GET", key)
	if err == nil {
		serr := v.Struct(&rcd)
		if serr == nil && !g.IsEmpty(rcd.ChainId) {
			return rcd
		}
	}

	temp, err := dao.CtxWithDatabase(ctx, chainName).Where("chain_name =?", chainName).One()
	if temp == nil {
		g.Log().Debug(ctx, "error:", err)
		return rcd
	}
	json.Unmarshal([]byte(temp.Json()), &rcd)

	//设置缓存过期时间10天
	g.Redis().Do(ctx, "SETEX", key, 3600*24, rcd)
	return rcd
}

func (dao *BlockChainDao) Update(ctx context.Context, chainName string, in model.BlockChainCreateInput) (err error) {
	dataModel := dao.CtxWithDatabase(ctx, chainName)
	if !g.IsEmpty(in.BaseCoin) && !g.IsEmpty(in.RpcUrls) {
		_, err = dataModel.Update("base_coin=?, rpc_urls=? ", "chain_name=? and chain_id=?", in.BaseCoin, in.RpcUrls, in.ChainName, in.ChainId)
	}
	if !g.IsEmpty(in.BaseCoin) && g.IsEmpty(in.RpcUrls) {
		_, err = dataModel.Update("base_coin=? ", "chain_name=? and chain_id=?", in.BaseCoin, in.ChainName, in.ChainId)
	}
	if g.IsEmpty(in.BaseCoin) && !g.IsEmpty(in.RpcUrls) {
		_, err = dataModel.Update(" rpc_urls=? ", "chain_name=? and chain_id=?", in.RpcUrls, in.ChainName, in.ChainId)
	}
	if err != nil {
		g.Log().Debug(ctx, "BlockChainDao Update error:", err, ",chain_name:", in.ChainName)
		return err
	}

	key := chainName + "_" + consts.BLOCK_CHAIN_DATA_KEY
	g.Redis().Do(ctx, "DEL", key)

	return nil
}

func (dao *BlockChainDao) ExecInsert(ctx context.Context, chainName string, in model.BlockChainCreateInput) error {
	var dimSqlDML bytes.Buffer
	dimSqlDML.WriteString("insert into ")
	dimSqlDML.WriteString(chainName + "." + dao.table)
	dimSqlDML.WriteString(" (chain_name, base_coin, chain_id, rpc_urls, update_time) values ( ")
	dimSqlDML.WriteString("'" + in.ChainName + "'")
	dimSqlDML.WriteString(",")
	dimSqlDML.WriteString("'" + in.BaseCoin + "'")
	dimSqlDML.WriteString(",")
	dimSqlDML.WriteString(strconv.FormatInt(int64(in.ChainId), 10))
	dimSqlDML.WriteString(",")
	dimSqlDML.WriteString("'" + in.RpcUrls + "'")
	dimSqlDML.WriteString(",")
	dimSqlDML.WriteString(strconv.FormatInt(int64(in.UpdateTime), 10))
	dimSqlDML.WriteString(")")

	_, err := dao.DB().Exec(ctx, dimSqlDML.String())
	if err != nil {
		g.Log().Debug(ctx, "insert error:", err)
		return err
	}

	key := chainName + "_" + consts.BLOCK_CHAIN_DATA_KEY + strconv.FormatInt(int64(in.ChainId), 10)
	_, err = g.Redis().Do(ctx, "DEL", key)
	if err != nil {
		g.Log().Debug(ctx, "Redis do error:", err)
		return err
	}
	return nil
}