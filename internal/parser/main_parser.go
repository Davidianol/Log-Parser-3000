package parser

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log_parser3000/internal/domain"
	"log_parser3000/internal/parser/raw"
	"log_parser3000/internal/service"
	"os"
	"path/filepath"
	"strings"
)

type Repository interface {
	SaveParsedLog(parsed *domain.ParsedLog) (int, error)
	SaveLogError(filename string, errMsg string) (int, error)
}

type MainParser struct {
	repo    Repository
	dataDir string
}

func NewMainParser(repo Repository, dataDir string) *MainParser {
	return &MainParser{repo: repo, dataDir: dataDir}
}

type inputFile struct {
	Name    string
	Content []byte
}

func (p *MainParser) ParseFromDataPath(relPath string) (int, error) {
	cleanPath := filepath.Clean(relPath)
	if strings.Contains(cleanPath, "..") {
		return 0, fmt.Errorf("invalid path")
	}

	fullPath := filepath.Join(p.dataDir, cleanPath)
	info, err := os.Stat(fullPath)
	if err != nil {
		return 0, fmt.Errorf("error with file %s: %w", relPath, service.ErrNotFound)
	}

	var files []inputFile
	var filename string

	switch {
	case info.IsDir():
		filename = cleanPath
		files, err = collectFromDir(fullPath)
	case isZip(fullPath):
		filename = cleanPath
		files, err = collectFromZip(fullPath)
	case isTarGz(fullPath):
		filename = cleanPath
		files, err = collectFromTarGz(fullPath)
	default:
		return 0, fmt.Errorf("unsupported path: must be a directory, .zip or .tar.gz")
	}

	if err != nil {
		logID, _ := p.repo.SaveLogError(cleanPath, err.Error())
		return logID, err
	}

	return p.parseFiles(files, filename)
}

func (p *MainParser) parseFiles(files []inputFile, filename string) (int, error) {
	var dbContent []byte
	var sharpContent []byte

	for _, f := range files {
		switch {
		case strings.EqualFold(f.Name, "ibdiagnet2.db_csv"):
			dbContent = f.Content
		case looksLikeSharpInfoFile(f.Name):
			sharpContent = f.Content
		}
	}

	if dbContent == nil {
		msg := "ibdiagnet2.db_csv not found"
		logID, _ := p.repo.SaveLogError(filename, msg)
		return logID, fmt.Errorf(msg)
	}

	dbData, err := ParseDBCSV(bytes.NewReader(dbContent))
	if err != nil {
		logID, _ := p.repo.SaveLogError(filename, err.Error())
		return logID, err
	}

	var sharpData []raw.SharpInfo
	if sharpContent != nil {
		sharpData, err = ParseSharpInfo(bytes.NewReader(sharpContent))
		if err != nil {
			logID, _ := p.repo.SaveLogError(filename, err.Error())
			return logID, err
		}
	}

	parsed, err := MapParsedDataToDomain(filename, dbData, sharpData)
	if err != nil {
		logID, _ := p.repo.SaveLogError(filename, err.Error())
		return logID, err
	}

	logID, err := p.repo.SaveParsedLog(parsed)
	if err != nil {
		return 0, err
	}

	return logID, nil
}

func collectFromDir(dirPath string) ([]inputFile, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}

	var files []inputFile
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		full := filepath.Join(dirPath, entry.Name())
		content, err := os.ReadFile(full)
		if err != nil {
			return nil, fmt.Errorf("read file %s: %w", entry.Name(), err)
		}
		files = append(files, inputFile{Name: entry.Name(), Content: content})
	}
	return files, nil
}

func collectFromZip(path string) ([]inputFile, error) {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}
	defer zr.Close()

	var files []inputFile
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("open zip entry %s: %w", f.Name, err)
		}

		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("read zip entry %s: %w", f.Name, err)
		}

		files = append(files, inputFile{
			Name:    filepath.Base(f.Name),
			Content: content,
		})
	}
	return files, nil
}

func collectFromTarGz(path string) ([]inputFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open tar.gz: %w", err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("init gzip reader: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)

	var files []inputFile
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read tar entry: %w", err)
		}
		if hdr.Typeflag == tar.TypeDir {
			continue
		}

		content, err := io.ReadAll(tr)
		if err != nil {
			return nil, fmt.Errorf("read tar file %s: %w", hdr.Name, err)
		}

		files = append(files, inputFile{
			Name:    filepath.Base(hdr.Name),
			Content: content,
		})
	}
	return files, nil
}

func isZip(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".zip")
}

func isTarGz(path string) bool {
	lower := strings.ToLower(path)
	return strings.HasSuffix(lower, ".tar.gz") || strings.HasSuffix(lower, ".tgz")
}

func looksLikeSharpInfoFile(name string) bool {
	lower := strings.ToLower(name)
	return strings.Contains(lower, "sharp_an_info")
}
