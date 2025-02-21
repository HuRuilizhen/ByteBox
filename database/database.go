package database

import (
	"bytebox/logger"
	"fmt"
	"os"
	"path"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

/* database table struct */

type File struct {
	Hash     string `gorm:"primaryKey"`
	FileName string
}

/* path parser function */

func getUserWorkDir(workDirName string) (string, error) {
	loggerInstance := logger.GetLoggerInstance()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		loggerInstance.Error(fmt.Sprintf("get user home dir %s failed", homeDir))
		return "", err
	}

	workDir := path.Join(homeDir, workDirName)
	if err := os.MkdirAll(workDir, 0755); err != nil {
		loggerInstance.Error(fmt.Sprintf("get user work dir %s failed", workDir))
		return "", err
	}

	return workDir, nil
}

func getDatabaseFile(workDirName string, databaseDirName string, databaseName string) (string, error) {
	loggerInstance := logger.GetLoggerInstance()

	workDir, err := getUserWorkDir(workDirName)
	if err != nil {
		loggerInstance.Error(fmt.Sprintf("get user work dir %s failed", workDir))
		return "", err
	}

	dbDir := path.Join(workDir, databaseDirName)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		loggerInstance.Error(fmt.Sprintf("get database dir %s failed", dbDir))
		return "", err
	}

	return path.Join(dbDir, databaseName), nil
}

func existDatabaseFile(workDirName string, databaseDirName string, databaseName string) (bool, error) {
	loggerInstance := logger.GetLoggerInstance()

	databaseFile, err := getDatabaseFile(workDirName, databaseDirName, databaseName)
	if err != nil {
		loggerInstance.Error(fmt.Sprintf("get database file %s failed", databaseFile))
		return false, err
	}

	_, err = os.Stat(databaseFile)
	if os.IsNotExist(err) {
		loggerInstance.Info("database file not found, try to initialize...")
		return false, nil
	}

	return true, nil
}

/* init database func */

func initDatabaseFile(workDirName string, databaseDirName string, databaseName string) *gorm.DB {
	loggerInstance := logger.GetLoggerInstance()
	databaseFile, err := getDatabaseFile(workDirName, databaseDirName, databaseName)

	if err != nil {
		loggerInstance.Fatal(fmt.Sprintf("get database file %s failed", databaseFile))
	}

	if exist, err := existDatabaseFile(workDirName, databaseDirName, databaseName); err != nil {
		loggerInstance.Fatal("check database file failed")
	} else if !exist {
		db, err := gorm.Open(sqlite.Open(databaseFile), &gorm.Config{})
		if err != nil {
			loggerInstance.Fatal(fmt.Sprintf("open database file %s failed", databaseFile))
		}
		if err := db.AutoMigrate(&File{}); err != nil {
			loggerInstance.Fatal("migrate database failed")
		}
		loggerInstance.Info("initialize database file successful")
		return db
	}

	db, err := gorm.Open(sqlite.Open(databaseFile), &gorm.Config{})
	if err != nil {
		loggerInstance.Fatal("open database failed")
	}

	loggerInstance.Info("open database successful")
	return db
}

/* database config setting */

type DatabaseConfig struct {
	workDirName     string
	databaseDirName string
	databaseName    string
}

var (
	databaseConfigInstance     DatabaseConfig
	databaseConfigInstanceOnce sync.Once
)

func GetDatabaseConfigInstance() DatabaseConfig {
	databaseConfigInstanceOnce.Do(func() {
		databaseConfigInstance.workDirName = "bytebox"
		databaseConfigInstance.databaseDirName = "db"
		databaseConfigInstance.databaseName = "bytebox.db"
	})
	return databaseConfigInstance
}

func (databaseConfig *DatabaseConfig) SetWorkDirName(workDirName string) {
	databaseConfig.workDirName = workDirName
}

func (databaseConfig *DatabaseConfig) SetDatabaseDirName(databaseDirName string) {
	databaseConfig.databaseDirName = databaseDirName
}

func (databaseConfig *DatabaseConfig) SetDatabaseName(databaseName string) {
	databaseConfig.databaseName = databaseName
}

/* init database by database config */
func initDatabase() *gorm.DB {
	databaseConfigInstance := GetDatabaseConfigInstance()
	workDirName := databaseConfigInstance.workDirName
	databaseDirName := databaseConfigInstance.databaseDirName
	databaseName := databaseConfigInstance.databaseName
	return initDatabaseFile(workDirName, databaseDirName, databaseName)
}

/* database instance */
var (
	databaseInstance     *gorm.DB
	databaseInstanceOnce sync.Once
)

func GetDatabaseInstance() *gorm.DB {
	databaseInstanceOnce.Do(func() {
		databaseInstance = initDatabase()
	})
	return databaseInstance
}
