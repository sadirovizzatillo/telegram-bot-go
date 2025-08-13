package extractor

import (
	"bytes"
	"os/exec"
	"strings"
)

func GetVideoURL(link string) (string, error) {
	cmd := exec.Command("yt-dlp", "-g", link)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}
