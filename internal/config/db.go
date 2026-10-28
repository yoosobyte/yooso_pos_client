package config

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"yoosoPosV3/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// GetAppDataDir 获取应用数据目录，Windows 优先 D 盘。
// 依次尝试各候选目录，返回第一个真正可写的；全部不可写时回退到系统临时目录，
// 避免目录创建失败后 SQLite 只报一句 "unable to open database file" 就退出。
func GetAppDataDir() string {
	for _, dir := range appDataDirs() {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("[数据目录] 不可用 %s: %v", dir, err)
			continue
		}
		if err := checkWritable(dir); err != nil {
			log.Printf("[数据目录] 不可写 %s: %v", dir, err)
			continue
		}
		return dir
	}

	dir := filepath.Join(os.TempDir(), "yoosoPosDb")
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("[数据目录] 临时目录也不可用 %s: %v", dir, err)
	}
	return dir
}

// appDataDirs 按优先级返回候选数据目录
func appDataDirs() []string {
	dirs := make([]string, 0, 2)
	if runtime.GOOS == "darwin" {
		if home, err := os.UserHomeDir(); err == nil {
			dirs = append(dirs, filepath.Join(home, "yoosoPosDb"))
		}
	} else if _, err := os.Stat("D:\\"); err == nil {
		dirs = append(dirs, "D:\\yoosoPosDb")
	}
	if cfg, err := os.UserConfigDir(); err == nil {
		dirs = append(dirs, filepath.Join(cfg, "yoosoPosDb"))
	}
	return dirs
}

// checkWritable 通过创建并删除临时文件确认目录可写
func checkWritable(dir string) error {
	f, err := os.CreateTemp(dir, ".write_test_*")
	if err != nil {
		return err
	}
	name := f.Name()
	f.Close()
	return os.Remove(name)
}

func getDBPath() string {
	return filepath.Join(GetAppDataDir(), "yooso-pos-free.db")
}

func InitDB() error {
	dbPath := getDBPath()
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		}),
	})
	if err != nil {
		log.Fatalf("数据库连接失败 (%s): %v", dbPath, err)
		return err
	}

	err = DB.AutoMigrate(
		&model.Setting{},
		&model.Category{},
		&model.Unit{},
		&model.Product{},
		&model.Customer{},
		&model.Supplier{},
		&model.PurchaseOrder{},
		&model.PurchaseOrderItem{},
		&model.SalesOrder{},
		&model.SalesOrderItem{},
		&model.InventoryLog{},
	)
	if err != nil {
		log.Fatalf("数据库迁移失败 (%s): %v", dbPath, err)
		return err
	}

	// 首次运行时执行默认数据初始化
	var count int64
	DB.Model(&model.Setting{}).Where("key = ?", "init_done").Count(&count)
	if count == 0 {
		initDefaultData()
	}
	// 确保必要设置项存在（兼容旧数据库升级）
	ensureSetting("init_done", "1")
	ensureSetting("store_name", "悦速用户")
	ensureSetting("avatar", "")
	ensureSetting("secondary_password", "123456")
	ensureSetting("client_id", generateUUID())

	log.Println("Database connected and migrated successfully!")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

// ensureSetting 确保设置项存在，若不存在则用默认值创建（兼容旧数据库升级）
func ensureSetting(key, value string) {
	var count int64
	DB.Model(&model.Setting{}).Where("key = ?", key).Count(&count)
	if count == 0 {
		DB.Create(&model.Setting{Key: key, Value: value})
	}
}

// generateUUID 生成一个简单的 UUID v4 格式字符串
func generateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// 极端情况回退：用时间戳
		return fmt.Sprintf("%x%x", time.Now().UnixNano(), time.Now().UnixMilli())
	}
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func initDefaultData() {
	defaultUnits := []model.Unit{
		{Name: "个"},
		{Name: "箱"},
		{Name: "件"},
		{Name: "斤"},
		{Name: "kg"},
		{Name: "袋"},
		{Name: "包"},
		{Name: "瓶"},
		{Name: "台"},
		{Name: "套"},
	}
	for _, u := range defaultUnits {
		var existing model.Unit
		err := DB.Unscoped().Where("name = ?", u.Name).First(&existing).Error
		if err == nil {
			// 存在（包括已逻辑删除的），恢复已删除的记录
			if existing.DeletedAt.Valid {
				DB.Unscoped().Model(&existing).Update("deleted_at", nil)
			}
		} else {
			DB.Create(&u)
		}
	}

	defaultCategories := []model.Category{
		{Name: "通用", Sort: 1},
	}
	for _, c := range defaultCategories {
		var existing model.Category
		err := DB.Unscoped().Where("name = ?", c.Name).First(&existing).Error
		if err == nil {
			if existing.DeletedAt.Valid {
				DB.Unscoped().Model(&existing).Update("deleted_at", nil)
			}
		} else {
			DB.Create(&c)
		}
	}
}
