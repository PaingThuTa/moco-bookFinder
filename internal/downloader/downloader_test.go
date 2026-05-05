package downloader

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

type mockDownloader struct {
	name string
	file *DownloadedFile
	err  error
}

func (m *mockDownloader) SourceName() string {
	return m.name
}

func (m *mockDownloader) DownloadFile(ctx context.Context, url string) (*DownloadedFile, error) {
	return m.file, m.err
}

func TestSourceManager_DownloadFile_Found(t *testing.T) {
	expected := &DownloadedFile{
		Data:     []byte("fake-epub-content"),
		Filename: "test.epub",
		FileType: "epub",
	}
	dls := []FileDownloader{
		&mockDownloader{name: "TestSource", file: expected},
	}
	mgr := NewSourceManager(dls, &http.Client{})

	result, err := mgr.DownloadFile(context.Background(), "TestSource", "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Filename != "test.epub" {
		t.Errorf("expected filename 'test.epub', got %q", result.Filename)
	}
	if result.FileType != "epub" {
		t.Errorf("expected type 'epub', got %q", result.FileType)
	}
}

func TestSourceManager_DownloadFile_UnknownSource(t *testing.T) {
	mgr := NewSourceManager(nil, &http.Client{})

	_, err := mgr.DownloadFile(context.Background(), "UnknownSource", "https://example.com")
	if err == nil {
		t.Fatal("expected error for unknown source")
	}
}

func TestSourceManager_DownloadFile_Error(t *testing.T) {
	dls := []FileDownloader{
		&mockDownloader{name: "FailSource", err: ErrNoFileFound},
	}
	mgr := NewSourceManager(dls, &http.Client{})

	_, err := mgr.DownloadFile(context.Background(), "FailSource", "https://example.com")
	if err == nil {
		t.Fatal("expected error from failing downloader")
	}
	if !errors.Is(err, ErrNoFileFound) {
		t.Errorf("expected ErrNoFileFound, got %v", err)
	}
}
