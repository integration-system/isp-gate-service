package model

import (
	"github.com/go-pg/pg"
	"github.com/integration-system/isp-lib/database"
	"isp-gate-service/entity"

	"github.com/go-pg/pg/orm"
)

type (
	SnapshotRepository interface {
		GetByApplication(int32) (*entity.Snapshot, error)
		Update([]entity.Snapshot) error
	}

	snapshotRepository struct {
		DB       orm.DB
		rxClient *database.RxDbClient
	}
)

func (r snapshotRepository) Update(list []entity.Snapshot) error {
	_, err := r.getDb().Model(&list).OnConflict("(app_id) DO UPDATE WHERE EXCLUDED.version > version").Insert()
	return err
}

func (r snapshotRepository) GetByApplication(appId int32) (*entity.Snapshot, error) {
	model := new(entity.Snapshot)
	if err := r.getDb().Model(model).Where("app_id = ?", appId).Select(); err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	} else {
		return model, nil
	}
}

func (r snapshotRepository) getDb() orm.DB {
	if r.DB != nil {
		return r.DB
	}
	return r.rxClient.Unsafe()
}
