package service

import (
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git/internal/config"
	"github.com/Baraulia/goLab/IndTask.git/internal/myErrors"
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
	"io/ioutil"
)

type FileService struct {
	cfg    *config.Config
	logger logging.Logger
}

func NewFileService(cfg *config.Config, logger logging.Logger) *FileService {
	return &FileService{cfg: cfg, logger: logger}
}

func (f *FileService) GetFile(path string) ([]byte, error) {
	filePath := fmt.Sprintf("%s/%s", f.cfg.FilePath, path)
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Errorf("Can not read file with path: %s: %s", filePath, err)
		return nil, &myErrors.MyError{Err: fmt.Errorf("can not read file with path: %s: %w", filePath, err), Code: 500}
	}
	return file, nil
}
