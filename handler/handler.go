package handler

import (
	"bytebox/configparser"
	"bytebox/database"
	"bytebox/logger"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gorm.io/gorm"
)

type handlerConfig struct {
	workDirName        string
	storageDirName     string
	maxHeaderMegabytes int64
}

var (
	handlerConfigInstanceOnce sync.Once
	handlerConfigInstance     handlerConfig
)

func LoadDatabaseConfig() {
	handlerConfigInstanceOnce.Do(func() {
		configInstance := configparser.GetConfigInstance()
		handlerConfigInstance.workDirName = configInstance["apihandler"].(map[string]interface{})["workDirName"].(string)
		handlerConfigInstance.storageDirName = configInstance["apihandler"].(map[string]interface{})["storageDirName"].(string)
		handlerConfigInstance.maxHeaderMegabytes = int64(configparser.GetConfigInstance()["server"].(map[string]interface{})["maxHeaderMegabytes"].(float64))
		storageDir := filepath.Join(handlerConfigInstance.workDirName, handlerConfigInstance.storageDirName)
		if err := os.MkdirAll(storageDir, 0755); err != nil {
			logger.GetLoggerInstance().Fatalf("create storage dir %s failed", storageDir)
		}
	})
}

func GetHandlerConfigInstance() handlerConfig {
	handlerConfigInstanceOnce.Do(func() {
		handlerConfigInstance = handlerConfig{
			workDirName:        ".temp",
			storageDirName:     "storage",
			maxHeaderMegabytes: 10,
		}
		storageDir := filepath.Join(handlerConfigInstance.workDirName, handlerConfigInstance.storageDirName)
		if err := os.MkdirAll(storageDir, 0755); err != nil {
			logger.GetLoggerInstance().Fatalf("create storage dir %s failed", storageDir)
		}
	})
	return handlerConfigInstance
}

func UploadHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	handlerConfigInstance := GetHandlerConfigInstance()
	maxHeaderMegabytes := handlerConfigInstance.maxHeaderMegabytes
	storageDir := filepath.Join(handlerConfigInstance.workDirName, handlerConfigInstance.storageDirName)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// parse the multipart form data
	err := r.ParseMultipartForm(maxHeaderMegabytes << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	fileNameOriginal := handler.Filename

	// generate a SHA-256 hash for the file
	fileHash := sha256.New()
	salt := fmt.Sprintf("%x", time.Now().UnixNano())
	fileHash.Write([]byte(salt))
	_, err = io.Copy(fileHash, file)
	if err != nil {
		http.Error(w, "Unable to generate hash for the file", http.StatusInternalServerError)
		return
	}

	// reset the file reader to the beginning
	file.Seek(0, 0)

	// create a new file with the hash value and original file extension
	hashString := hex.EncodeToString(fileHash.Sum(nil))
	fileExtension := filepath.Ext(handler.Filename)
	fileName := hashString + fileExtension

	tempFile, err := os.Create(filepath.Join(storageDir, fileName))
	if err != nil {
		http.Error(w, "Unable to create file on server", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	// copy the uploaded file content to the newly created file
	if _, err := io.Copy(tempFile, file); err != nil {
		http.Error(w, "Unable to save the file", http.StatusInternalServerError)
		return
	}

	// insert data into database
	fileModel := database.File{Hash: hashString, FileName: fileNameOriginal}
	result := db.Create(&fileModel)
	if result.Error != nil {
		http.Error(w, "Database insert failed", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"message": "File uploaded successfully, hash prefix: " + hashString[:8],
		"hash":    hashString,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hashPrefix := r.PostFormValue("hashcode")

	var files []database.File
	result := db.Where("hash LIKE ?", hashPrefix+"%").Find(&files)
	if result.Error != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}

	switch len(files) {
	case 0:
		w.Header().Set("Content-Type", "application/json")
		data := map[string]interface{}{
			"message": "No file found with the given hash prefix.",
		}
		json.NewEncoder(w).Encode(data)
	case 1:
		storageDir := filepath.Join(GetHandlerConfigInstance().workDirName, GetHandlerConfigInstance().storageDirName)
		fileHash := files[0].Hash
		fileNameOriginal := files[0].FileName
		fileExt := filepath.Ext(fileNameOriginal)
		filePath := filepath.Join(storageDir, fileHash+fileExt)
		file, err := os.Open(filePath)

		if err != nil {
			http.Error(w, "Failed to open file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			http.Error(w, "Failed to get file info", http.StatusInternalServerError)
			return
		}

		encodedFileName := url.QueryEscape(fileNameOriginal)
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, encodedFileName, encodedFileName))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", (fileInfo.Size())))
		io.Copy(w, file)

		os.Remove(filePath)
		result := db.Delete(&files[0])
		if result.Error != nil {
			http.Error(w, "Failed to delete file from database", http.StatusInternalServerError)
			return
		}
	default:
		w.Header().Set("Content-Type", "application/json")
		data := map[string]interface{}{
			"message": "Multiple files found with the given hash prefix. Please provide a more specific hash.",
		}
		json.NewEncoder(w).Encode(data)
	}
}
