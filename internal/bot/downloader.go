package bot

import (
	"fmt"
	"io"
	"net/http"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type FileDownloader struct {
	bot *gotgbot.Bot
	hc  *http.Client
}

func NewFileDownloader(bot *gotgbot.Bot) *FileDownloader {
	hc := &http.Client{}

	return &FileDownloader{
		bot: bot,
		hc:  hc,
	}
}

func (fd *FileDownloader) GetFileURL(fileID string) (string, error) {
	file, err := fd.bot.GetFile(fileID, nil)
	if err != nil {
		return "", err
	}

	return file.URL(fd.bot, nil), nil
}

func (fd *FileDownloader) DownloadFile(fileURL string) ([]byte, error) {
	resp, err := fd.hc.Get(fileURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}
