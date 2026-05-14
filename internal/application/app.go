package application

import (
	"fmt"
	"io"
	"os"

	"name-counter/internal/domain"
)

type App struct {
	config  *domain.Config
	counter *domain.Counter
}

func New(config *domain.Config) *App {
	return &App{
		config:  config,
		counter: domain.NewCounter(),
	}
}

func (app *App) Run(out io.Writer, errOut io.Writer) error {
	file, err := os.Open(app.config.FilePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл %s: %w", app.config.FilePath, err)
	}
	defer file.Close()

	results, err := app.counter.Count(file, app.config.SortMode)
	if err != nil {
		return fmt.Errorf("ошибка при подсчете: %w", err)
	}

	for _, result := range results {
		_, err := fmt.Fprintf(out, "%s: %d\n", result.Name, result.Count)
		if err != nil {
			return fmt.Errorf("ошибка вывода: %w", err)
		}
	}

	return nil
}

func (app *App) RunWithDefaultOutput() error {
	return app.Run(os.Stdout, os.Stderr)
}
