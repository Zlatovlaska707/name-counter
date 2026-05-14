package main

import (
	"log"
	"name-counter/pkg/errwriter"
	"os"
	"os/signal"
	"syscall"

	"name-counter/internal/application"
	"name-counter/internal/cli"
	"name-counter/internal/domain"
	"name-counter/internal/pprof"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Парсим аргументы командной строки
	parser := cli.NewParser()
	cfg, err := parser.Parse()
	if err != nil {
		errwriter.PrintErrln("Ошибка:", err)
		parser.PrintUsage()
		os.Exit(2)
	}

	// Показываем справку если нужно
	if cfg.ShowHelp {
		parser.PrintUsage()
		return
	}

	// Запускаем профилировщик
	var pprofServer *pprof.Server
	if cfg.PprofAddr != "" {
		pprofServer = pprof.NewServer(cfg.PprofAddr)
		if err := pprofServer.Start(); err != nil {
			log.Printf("Предупреждение: не удалось запустить pprof сервер: %v", err)
		}
	}

	app := application.New(&domain.Config{
		FilePath: cfg.FilePath,
		SortMode: cfg.SortMode,
	})

	if err := app.Run(os.Stdout, os.Stderr); err != nil {
		errwriter.PrintErrln("Ошибка выполнения:", err)
		if pprofServer != nil {
			if stopErr := pprofServer.Stop(); stopErr != nil {
				log.Printf("Предупреждение: ошибка при остановке pprof сервера: %v", stopErr)
			}
		}
		os.Exit(1)
	}

	waitForPprofShutdown(pprofServer)
}

func waitForPprofShutdown(pprofServer *pprof.Server) {
	if pprofServer == nil {
		return
	}

	log.Println("Основная работа завершена. pprof сервер продолжает работу.")
	log.Println("Для остановки нажмите Ctrl+C или отправьте SIGTERM")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Получен сигнал завершения, пока-пока сервер...")

	if err := pprofServer.Stop(); err != nil {
		log.Printf("Ошибка при остановке pprof сервера: %v", err)
	}
	log.Println("Программа завершена")

}
