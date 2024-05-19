package raft

import "sync"

type ApplicationInterface interface {
	Ping() string
	Get(key any) string
	Set(key any, val string) string
	Strln(key any) int
	Del(key any) string
	Append(key any, val string) string
}

type Application struct {
	// https://pkg.go.dev/sync#Map
	// sync.Map usage anticipates multiple
	// goroutines read, write, and overwrite
	// entries for disjoint sets of keys
	data sync.Map
}

func (app *Application) Ping() string {
	return "PONG"
}

func (app *Application) Get(key any) string {
	if val, ok := app.data.Load(key); ok {
		return val.(string)
	}

	return ""
}

func (app *Application) Set(key any, val string) string {
	app.data.Store(key, val)
	return "OK"
}

func (app *Application) Strln(key any) int {
	if val, ok := app.data.Load(key); ok {
		return len(val.(string))
	}

	return 0
}

func (app *Application) Del(key any) string {
	if val, ok := app.data.Load(key); ok {
		app.data.Delete(key)
		return val.(string)
	}

	return ""
}

func (app *Application) Append(key any, val string) string {
	if v, ok := app.data.Load(key); ok {
		newVal := v.(string) + val // Concatenate the strings
		app.data.Store(key, newVal)
		return "OK"
	}

	// Set the value to empty string first before appending
	app.data.Store(key, "")
	return app.Append(key, val)
}
