package data

import (
	"cache/config"
	"cache/internal/interface/getter"
	"cache/internal/model/information"
	clog "cache/pkg/log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

var CDB *gorm.DB

func Sqlinit() {
	path := config.ConfigPath()
	viper.SetConfigName("app")
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		clog.Panic("[mysql ]", err)
	}
	// var err error
	CDB, err = gorm.Open("mysql", viper.GetString("mysql.dns"))
	if err != nil {
		clog.Panic("[mysql ]", err)
	}
	// CDB.LogMode(true)
	clog.Info("[mysql] 数据库连接成功")
}

func Get() getter.GetterFunc {
	return getter.GetterFunc(func(key string) ([]byte, error) {
		Info := &information.Information{}
		// log.Println("[slowdb]====",)
		result := CDB.Where("k = ?", key).First(Info)

		// 判断查询结果
		if result.Error == gorm.ErrRecordNotFound {
			// 没有找到数据
			clog.Info("[mysql] 查询 |key:" + "数据库中相关无数据|")
			return []byte("数据库中相关无数据"), nil
		} else if result.Error != nil {
			// 查询过程中发生错误
			return nil, result.Error
		} else {
			// 找到数据，可以在user中访问
			clog.Info("[mysql] 查询 |key: ", key, " Value: ", Info.V, "|")
			return []byte(Info.V), nil
		}
	})
}
