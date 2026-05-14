package cli

import (
	"errors"
	"flag"
	"fmt"
	"name-counter/pkg/errwriter"
	"os"

	"name-counter/internal/domain"
)

type Config struct {
	FilePath  string
	SortMode  domain.SortMode
	PprofAddr string
	ShowHelp  bool
}

type Parser struct {
	// Флаги командной строки
	filePath  *string
	sortMode  *string
	pprofAddr *string
	help      *bool
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse() (*Config, error) {
	// Инициализируем флаги
	p.filePath = flag.String("file", "", "путь к входному файлу (одно имя на строку, UTF-8)")
	p.sortMode = flag.String("sort", "freq", "порядок сортировки: freq (по частоте убыв.) или name (алфавитный)")
	p.pprofAddr = flag.String("pprof", "", "адрес для pprof (например :6060)")
	p.help = flag.Bool("help", false, "показать справку")

	flag.Parse()

	// Проверяем запрос на справку
	if *p.help {
		return &Config{ShowHelp: true}, nil
	}

	// Определяем путь к файлу
	filePath := p.determineFilePath()
	if filePath == "" {
		return nil, errors.New("не указан путь к файлу")
	}

	// Парсим режим сортировки
	sortMode, err := p.parseSortMode(*p.sortMode)
	if err != nil {
		return nil, err
	}

	return &Config{
		FilePath:  filePath,
		SortMode:  sortMode,
		PprofAddr: *p.pprofAddr,
		ShowHelp:  false,
	}, nil
}

func (p *Parser) determineFilePath() string {
	if *p.filePath != "" {
		return *p.filePath
	}
	if flag.NArg() == 1 {
		return flag.Arg(0)
	}
	return ""
}

func (p *Parser) parseSortMode(s string) (domain.SortMode, error) {
	switch s {
	case "freq", "frequency":
		return domain.SortByFrequencyDesc, nil
	case "name":
		return domain.SortByNameAsc, nil
	default:
		return 0, fmt.Errorf("неверный -sort %q: используйте freq или name", s)
	}
}

func (p *Parser) PrintUsage() {
	errwriter.PrintErrf("Использование: %s [-флаг] <путь_к_файлу>\n", os.Args[0])
	errwriter.PrintErrln("\nфлажки:")
	flag.PrintDefaults()
	errwriter.PrintErrln("\nПримеры:")
	errwriter.PrintErrln("  run names.txt")
	errwriter.PrintErrln("  run -file names.txt -sort name")
	errwriter.PrintErrln("  run -file names.txt -pprof :6060")
}
