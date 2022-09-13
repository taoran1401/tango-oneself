package repo

import (
	"github.com/jinzhu/gorm"
	"taogin/config/global"
)

var (
	r *Repository
)

type Repository struct {
	db *gorm.DB
}

func GetRepository() *Repository {
	if r == nil {
		return &Repository{
			db: global.DB["colorful"],
		}
	}
	return r
}

type RepositoryInterface interface {
	// 根据 id 查询
	GetById(id uint64) (interface{}, error)

	//  根据 ids 查询列表
	GetListByIds(ids []uint64) (res interface{}, err error)

	// 创建
	Create(data *interface{}) (uint64, error)

	// 编辑
	Update(data map[string]interface{}) error

	// barch编辑
	UpdateBatch(datas []map[string]interface{}) error

	// 删除
	Delete(ids []uint64, updatedBy uint64) error

	// 简单条件查询
	QueryByMaps(data map[string]interface{}, pageSize, current int, order string) (res []interface{}, err error)

	// 简单条件查询数量
	QueryTotalByMaps(data map[string]interface{}) (res int64, err error)

	// 简单条件查询
	QueryByWhere(pageSize, current int, order, where string, param ...interface{}) (res []interface{}, err error)

	// 简单条件查询数量
	QueryTotalByWhere(where string, param ...interface{}) (res int64, err error)
}
