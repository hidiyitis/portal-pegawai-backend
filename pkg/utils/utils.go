package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func SaveFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	var baseDir, _ = os.Getwd()
	var folder = "document"
	fileExt := filepath.Ext(header.Filename)
	if fileExt == ".jpg" || fileExt == ".png" || fileExt == ".jpeg" {
		folder = "images"
	}
	var publicDir = filepath.Join(baseDir, "public", folder)
	var fileName = fmt.Sprintf("%s%s", time.Now().Format("20060102150405"), filepath.Ext(header.Filename))
	filePath := filepath.Join(publicDir, fileName)

	if err := os.MkdirAll(publicDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create public directory: %v", err)
	}

	fileTemp, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	fmt.Println(fileTemp.Name())
	defer func(fileTemp *os.File) {
		err := fileTemp.Close()
		if err != nil {
			return
		}
	}(fileTemp)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		errStr := fmt.Sprintf("Error in reading the file buffer %s\n", err)
		return "", errors.New(errStr)
	}

	_, err = fileTemp.Write(fileBytes)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

func ReadFile(filename string) ([]byte, error) {
	baseDir, _ := os.Getwd()
	var publicDir = filepath.Join(baseDir, "public")
	filePath := filepath.Join(publicDir, filename)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}
