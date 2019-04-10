package app

import (
	"net/http"
)

func (a *App) cronWeatherHandler(w http.ResponseWriter, r *http.Request) *appError {
	message, err := a.weatherMessage()
	if err != nil {
		return &appError{err, "failed to fetch weather"}
	}
	message.Channel = a.weather["CHANNEL"]
	message.UserName = a.weather["USERNAME"]
	message.IconEmoji = a.weather["ICONEMOJI"]

	if err := a.sendMessage(message); err != nil {
		return &appError{err, "failed to send message"}
	}
	return nil
}

func (a *App) cronWishlistHandler(w http.ResponseWriter, r *http.Request) *appError {
	message, err := a.wishlistMessage()
	if err != nil {
		return &appError{err, "failed to fetch wishlist"}
	}
	message.Channel = a.wishlistChannel
	message.UserName = "wishlist"
	message.IconEmoji = ":memo:"

	if err := a.sendMessage(message); err != nil {
		return &appError{err, "failed to send message"}
	}
	return nil
}
