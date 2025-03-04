package handlers

import "gopkg.in/telebot.v3"

func GetFileURL(bot *telebot.Bot, fileID string) (string, error) {
	file, err := bot.FileByID(fileID)
	if err != nil {
		return "", err
	}

	return "https://api.telegram.org/file/bot" + bot.Token + "/" + file.FilePath, nil
}