package service

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

const (
	outdir = "./public/pdf"
)

// SrvConvertFileService describes the service.
type SrvConvertFileService interface {
	// Add your methods here
	WordToPDF(ctx context.Context, file io.Reader, name, mimeType string) (rs string, err error)
}

type basicSrvConvertFileService struct{}

func (b *basicSrvConvertFileService) WordToPDF(ctx context.Context, file io.Reader, name, mimeType string) (rs string, err error) {
	fileName := uuid.New().String()
	fileDir := "./public/upload/" + fileName
	if err != nil {
		fmt.Println(err)
	}
	// defer tempFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	err = ioutil.WriteFile(fileDir, fileBytes, 0644)
	if err != nil {
		fmt.Println(err)
	}
	rs, err = converter(fileDir, fileName)
	return rs, err
}

func converter(path string, fileName string) (rs string, err error) {
	cmd := exec.Command("soffice", "--headless", "--convert-to", "pdf", "./"+path, "--outdir", outdir)
	stderr, err := cmd.StderrPipe()
	log.SetOutput(os.Stderr)

	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	slurp, _ := ioutil.ReadAll(stderr)
	fmt.Printf("%s\n", slurp)

	if err := cmd.Wait(); err != nil {
		return "", err
	}
	return "/static/" + fileName + ".pdf", nil
}

// NewBasicSrvConvertFileService returns a naive, stateless implementation of SrvConvertFileService.
func NewBasicSrvConvertFileService() SrvConvertFileService {
	return &basicSrvConvertFileService{}
}

// New returns a SrvConvertFileService with all of the expected middleware wired in.
func New(middleware []Middleware) SrvConvertFileService {
	var svc SrvConvertFileService = NewBasicSrvConvertFileService()
	for _, m := range middleware {
		svc = m(svc)
	}
	return svc
}
