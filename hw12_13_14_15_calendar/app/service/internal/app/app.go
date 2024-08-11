package app

import (
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/lib/storage"
)

type (
	Storage = storage.Storage
)

type App struct {
	storage Storage
}

func New(storage Storage) *App {
	return &App{storage: storage}
}

func (a *App) Storage() storage.Storage {
	return a.storage
}
