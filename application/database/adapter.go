package database

import (
	"github.com/cocktailcloud/acloud-alarm-collector/application/model"
	db "github.com/cocktailcloud/acloud-monitoring-common/v2/database"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/database/impl"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/logger"
)

type DBAdapter struct {
	ready bool
	DB    db.DB
}

func NewAdapter(conf *impl.Config) (DBAdapter, error) {
	adapter := DBAdapter{
		ready: false,
	}

	db, err := impl.NewConnection(conf)
	if err != nil {
		logger.WithError(err).Fatal("Could not open database connection")
		return adapter, err
	}
	adapter.DB = db
	adapter.ready = true
	return adapter, nil
}

func (d *DBAdapter) OrmMapping() {
	d.DB.AddDbMap(model.AlarmHistory{}, "alarm_history")
	d.DB.AddDbMap(model.AlarmCurrent{}, "alarm_current")
	d.DB.AddDbMap(model.KubernetesEvent{}, "kubernetes_event")
}

func (d *DBAdapter) GetConn() db.DB {
	return d.DB
}

func (d *DBAdapter) IsNotReady() bool {
	return !d.ready
}

func (d *DBAdapter) Shutdown() {
	if d.DB != nil {
		d.DB.CloseConnection()
	}
	logger.Infof("Shutdown DB Adapter...")
}
