package bot

import (
	"fmt"
	"io"
	"net/http"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type fileDownloader struct {
	bot *gotgbot.Bot
}

func NewFileDownloader(bot *gotgbot.Bot) *fileDownloader {
	return &fileDownloader{bot: bot}
}

func (fd *fileDownloader) GetFileURL(fileID string) (string, error) {
	file, err := fd.bot.GetFile(fileID, nil)
	if err != nil {
		return "", err
	}
	return file.URL(fd.bot, nil), nil
}

func (fd *fileDownloader) DownloadFile(fileURL string) ([]byte, error) {
	resp, err := http.Get(fileURL)
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
