package bootstrap

import (
	"context"
	"time"
)

func Shutdown(app *App) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	if app.DB != nil {
		app.DB.Close()
	}

	<-ctx.Done()
}
