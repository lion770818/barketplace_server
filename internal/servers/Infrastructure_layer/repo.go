package Infrastructure_layer

import (
	"fmt"
	"marketplace_server/config"
	"time"

	//"marketplace_server/internal/backpack"
	Infrastructure_backpack "marketplace_server/internal/backpack/Infrastructure_layer"

	model_backpack "marketplace_server/internal/backpack/model"

	model_transaction "marketplace_server/internal/bill/model"
	model_user "marketplace_server/internal/user/model"

	Infrastructure_bill "marketplace_server/internal/bill/Infrastructure_layer"
	Infrastructure_product "marketplace_server/internal/product/Infrastructure_layer"
	model_product "marketplace_server/internal/product/model"
	Infrastructure_user "marketplace_server/internal/user/Infrastructure_layer"

	"marketplace_server/internal/common/logs"
	"marketplace_server/internal/common/mysql"
	"marketplace_server/internal/common/rabbitmqx"

	"marketplace_server/internal/common/redis"

	"github.com/jinzhu/gorm"

	//  mysql driver
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// 持久化管理物件
type RepositoriesManager struct {
	AuthRepo        Infrastructure_user.AuthInterface    // 驗證
	UserRepo        Infrastructure_user.UserRepo         // 用戶
	TransactionRepo Infrastructure_bill.TransactionRepo  // 交易
	ProductRepo     Infrastructure_product.ProductRepo   // 產品持久層
	BackpackRepo    Infrastructure_backpack.BackpackRepo // 背包持久層
	db              *gorm.DB
}

// 建立持久化管理物件
func NewRepositories(cfg *config.Config) *RepositoriesManager {

	// 持久化类型的 repo
	mysqlCfg := mysql.Config{
		LogMode:  cfg.Mysql.LogMode,
		Driver:   cfg.Mysql.Driver,
		Host:     cfg.Mysql.Host,
		Port:     cfg.Mysql.Port,
		Database: cfg.Mysql.Database,
		User:     cfg.Mysql.User,
		Password: cfg.Mysql.Password,
	}
	logs.Debugf("mysqlCfg=%+v", mysqlCfg)
	db := mysql.NewDB(mysqlCfg)

	// 建立redis連線
	redisCfg := &redis.RedisParameter{
		Network:      "tcp",
		Address:      fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           0,
		DialTimeout:  time.Second * time.Duration(10),
		ReadTimeout:  time.Second * time.Duration(10),
		WriteTimeout: time.Second * time.Duration(10),
		PoolSize:     10,
	}
	logs.Debugf("redisCfg=%+v", redisCfg)
	redisClient, err := redis.NewRedis(redisCfg)
	if err != nil {
		logs.Errorf("newRedis error=%v", err)
		return nil
	}
	// 初始化 rabbit mq
	err = rabbitmqx.Init(
		cfg.RabbitMq.Host,
		cfg.RabbitMq.Port,
		cfg.RabbitMq.User,
		cfg.RabbitMq.Password,
		cfg.RabbitMq.ConnectNum,
		cfg.RabbitMq.ChannelNum)
	if err != nil {
		logs.Errorf("rabbitmqx Init err:%v", err)
		return nil
	}

	transactionRepo := Infrastructure_bill.NewMysqlTransactionRepo(db)
	protuctRepo := Infrastructure_product.NewProductRepoManager(db, redisClient.GetClient())
	// user 和 產品
	userRepo := Infrastructure_user.NewMysqlUserRepo(db, redisClient.GetClient())
	backpackRepo := Infrastructure_backpack.NewMysqlBackpackRepo(db)

	// auth 策略
	var authRepo Infrastructure_user.AuthInterface
	if cfg.Auth.Active == "redis" {
		logs.Debugf("使用redis當驗證緩存")
		authRepo = Infrastructure_user.NewRedisAuthRepo(redisClient.GetClient(), cfg.AuthExpireTime)
	} else {
		logs.Debugf("使用jwt當驗證緩存")
		authRepo = Infrastructure_user.NewJwtAuth(cfg.Auth.PrivateKey, cfg.AuthExpireTime)
	}

	return &RepositoriesManager{
		AuthRepo:        authRepo,
		UserRepo:        userRepo,
		TransactionRepo: transactionRepo,
		ProductRepo:     protuctRepo,
		BackpackRepo:    backpackRepo,
		db:              db,
	}
}

// closes the  database connection
func (s *RepositoriesManager) Close() error {
	return s.db.Close()
}

func (s *RepositoriesManager) GetDB() *gorm.DB {
	return s.db
}

// This migrate all tables
func (s *RepositoriesManager) Automigrate() error {
	return s.db.AutoMigrate(&model_user.UserPO{},
		&model_transaction.Transaction_PO{},
		&model_product.Product_PO{},
		&model_backpack.Backpack_PO{}).Error
}
