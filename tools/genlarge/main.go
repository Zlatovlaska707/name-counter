package main

// Тут может быть немного говнокод, но это чисто инструментарная штука

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"name-counter/pkg/safeclose"
	"os"
	"sort"
	"time"

	"name-counter/pkg/errwriter"
)

type VerificationData struct {
	Seed          int64            `json:"seed"`
	TotalLines    int64            `json:"total_lines"`
	Frequencies   map[string]int64 `json:"frequencies"`
	WeightedNames []WeightedName   `json:"weighted_names,omitempty"`
}

type WeightedName struct {
	Name   string
	Weight int
}

func main() {
	mode := flag.String("mode", "generate", "Режим работы: generate (генерация) или verify (верификация)")

	outPath := flag.String("out", "../testdata/large_names.txt", "Путь к выходному файлу с именами")
	n := flag.Int64("n", 10_000_000, "Количество строк (имён) для генерации")
	flushEvery := flag.Int64("progress", 1_000_000, "Выводить прогресс каждые N строк (0 = отключено)")
	seed := flag.Int64("seed", 0, "Зерно генератора случайных чисел (0 = использовать текущее время)")

	// тут для режима верификации доп.
	verifyPath := flag.String("verify", "", "Путь к файлу с именами для проверки")
	statsPath := flag.String("stats", "../testdata/verification.json", "Путь к файлу с эталонными данными")

	flag.Parse()

	// основоной режим определяет --mode
	switch *mode {
	case "generate":
		if *n <= 0 {
			errwriter.PrintErrln("Ошибка: параметр -n должен быть больше 0")
			os.Exit(2)
		}

		rngSeed := *seed
		if rngSeed == 0 {
			rngSeed = time.Now().UnixNano()
		}

		// создаем пул имён с разными весами и получаем карту ожидаемых частот
		names, weights, freqMap := buildWeightedNamePool()

		verificationData := VerificationData{
			Seed:        rngSeed,
			TotalLines:  *n,
			Frequencies: freqMap,
		}

		if err := saveVerificationData(*statsPath, verificationData); err != nil {
			errwriter.PrintErrf("Предупреждение: не удалось сохранить эталонные данные: %v\n", err)
		} else {
			errwriter.PrintErrf("Эталонные данные сохранены в %s (seed: %d)\n", *statsPath, rngSeed)
		}

		if err := writeFileWeighted(*outPath, *statsPath, names, weights, *n, *flushEvery, rngSeed); err != nil {
			errwriter.PrintErrln(err)
			os.Exit(1)
		}

	case "verify":
		if *verifyPath == "" {
			errwriter.PrintErrln("Ошибка: в режиме verify необходимо указать -verify")
			os.Exit(2)
		}

		if err := verifyCounts(*verifyPath, *statsPath, *flushEvery); err != nil {
			errwriter.PrintErrln(err)
			os.Exit(1)
		}

	default:
		errwriter.PrintErrf("Неизвестный режим: %s. Используйте 'generate' или 'verify'\n", *mode)
		os.Exit(2)
	}
}

func buildWeightedNamePool() ([]string, []int, map[string]int64) {

	weightedNames := []WeightedName{
		// Очень популярные имена (вес 100)
		{"Александр", 100}, {"Алексей", 100}, {"Андрей", 100}, {"Дмитрий", 100},
		{"Елена", 100}, {"Ирина", 100}, {"Мария", 100}, {"Михаил", 100},
		{"Наталья", 100}, {"Ольга", 100}, {"Сергей", 100}, {"Татьяна", 100},
		{"Владимир", 100}, {"Екатерина", 100}, {"Николай", 100}, {"Анна", 100},

		// Популярные имена (вес 50)
		{"Артём", 50}, {"Виктор", 50}, {"Виталий", 50}, {"Галина", 50},
		{"Григорий", 50}, {"Дарья", 50}, {"Евгений", 50}, {"Евгения", 50},
		{"Иван", 50}, {"Игорь", 50}, {"Ксения", 50}, {"Людмила", 50},
		{"Максим", 50}, {"Надежда", 50}, {"Олег", 50}, {"Павел", 50},
		{"Пётр", 50}, {"Роман", 50}, {"Светлана", 50}, {"Юлия", 50},

		// Среднепопулярные имена (вес 20)
		{"Алла", 20}, {"Анжелика", 20}, {"Борис", 20}, {"Валентин", 20},
		{"Валерий", 20}, {"Вера", 20}, {"Владислав", 20}, {"Даниил", 20},
		{"Денис", 20}, {"Диана", 20}, {"Жанна", 20}, {"Захар", 20},
		{"Зоя", 20}, {"Инна", 20}, {"Кирилл", 20}, {"Кристина", 20},
		{"Лариса", 20}, {"Леонид", 20}, {"Лидия", 20}, {"Лилия", 20},

		// Редкие имена (вес 5)
		{"Авдей", 5}, {"Алевтина", 5}, {"Аполлон", 5}, {"Василиса", 5},
		{"Георгий", 5}, {"Евдоким", 5}, {"Ерофей", 5}, {"Ипполит", 5},
		{"Кузьма", 5}, {"Лукьян", 5}, {"Мефодий", 5}, {"Никодим", 5},
		{"Оксана", 5}, {"Пантелей", 5}, {"Раиса", 5}, {"Спиридон", 5},
		{"Трофим", 5}, {"Ульяна", 5}, {"Феоктист", 5}, {"Эльвира", 5},

		// Очень редкие имена (вес 1)
		{"Аристарх", 1}, {"Вавила", 1}, {"Герасим", 1}, {"Досифей", 1},
		{"Евстигней", 1}, {"Исидор", 1}, {"Капитон", 1}, {"Лаврентий", 1},
		{"Михей", 1}, {"Никанор", 1}, {"Онисим", 1}, {"Пафнутий", 1},
		{"Родион", 1}, {"Серафим", 1}, {"Терентий", 1}, {"Феодосий", 1},
	}

	names := make([]string, 0, len(weightedNames))
	weights := make([]int, 0, len(weightedNames))

	freqMap := make(map[string]int64)

	for _, wn := range weightedNames {
		names = append(names, wn.Name)
		weights = append(weights, wn.Weight)
		freqMap[wn.Name] = 0 // пока на стадии инициализации ставим ноль, потому что хз какая частотность
	}

	totalWeight := 0
	for _, w := range weights {
		totalWeight += w
	}

	errwriter.PrintErrf("Создан пул из %d уникальных имён с разными весами\n", len(names))
	errwriter.PrintErrf("Общий вес: %d\n", totalWeight)
	errwriter.PrintErrf("Ожидаемое распределение: очень частые (~%.1f%%), частые (~%.1f%%), средние (~%.1f%%), редкие (~%.1f%%), очень редкие (~%.1f%%)\n",
		100.0*100/float64(totalWeight),
		100.0*50/float64(totalWeight),
		100.0*20/float64(totalWeight),
		100.0*5/float64(totalWeight),
		100.0*1/float64(totalWeight))

	return names, weights, freqMap
}

func saveVerificationData(path string, data VerificationData) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer safeclose.CloseFile(f, path)

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")

	return encoder.Encode(data)
}

func loadVerificationData(path string) (*VerificationData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer safeclose.CloseFile(f, path)

	var data VerificationData

	decoder := json.NewDecoder(f)

	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func writeFileWeighted(path, statsPath string, names []string, weights []int, total, progressEvery, seed int64) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer safeclose.CloseFile(f, path)

	w := bufio.NewWriterSize(f, 4<<20) // 4 МБ (4 * 2^20)

	start := time.Now()

	totalWeight := 0
	for _, w := range weights {
		totalWeight += w
	}

	rng := rand.New(rand.NewSource(seed))

	// карта для подсчёта реальных частот, нужно для статистики
	actualCounts := make(map[string]int64)

	for i := int64(0); i < total; i++ {

		idx := weightedRandom(rng, weights, totalWeight)
		name := names[idx]

		if _, err := w.WriteString(name); err != nil {
			return err
		}

		if err := w.WriteByte('\n'); err != nil {
			return err
		}

		actualCounts[name]++

		if progressEvery > 0 && (i+1)%progressEvery == 0 {
			elapsed := time.Since(start).Round(time.Millisecond)
			errwriter.PrintErrf("Записано %d / %d строк (прошло %s)\n", i+1, total, elapsed)
		}
	}

	if err := w.Flush(); err != nil {
		return err
	}

	if err := f.Sync(); err != nil {
		return err
	}

	updateVerificationDataWithActualCounts(statsPath, actualCounts, total, seed)

	printGenerationStats(actualCounts, total, path, start)

	return nil
}

func updateVerificationDataWithActualCounts(statsPath string, actualCounts map[string]int64, totalLines, seed int64) {
	var data VerificationData

	if file, err := os.Open(statsPath); err == nil {
		defer safeclose.CloseFile(file, statsPath)
		decoder := json.NewDecoder(file)
		_ = decoder.Decode(&data)
	}

	data.Seed = seed
	data.TotalLines = totalLines
	data.Frequencies = actualCounts

	file, err := os.Create(statsPath)
	if err != nil {
		errwriter.PrintErrf("Предупреждение: не удалось обновить эталонные данные: %v\n", err)
		return
	}
	defer safeclose.CloseFile(file, statsPath)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		errwriter.PrintErrf("Предупреждение: не удалось сохранить эталонные данные: %v\n", err)
	}
}

func printGenerationStats(actualCounts map[string]int64, total int64, path string, start time.Time) {
	errwriter.PrintErrln("\n--- Статистика генерации ---")
	errwriter.PrintErrf("Всего строк: %d\n", total)
	errwriter.PrintErrf("Уникальных имён: %d\n", len(actualCounts))
	errwriter.PrintErrln("Топ 10 самых частых имён:")

	// Сортируем и выводим топ
	type countPair struct {
		name  string
		count int64
	}

	pairs := make([]countPair, 0, len(actualCounts))
	for name, count := range actualCounts {
		pairs = append(pairs, countPair{name, count})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	// Выводим топ 10
	topN := 10
	for i := 0; i < topN && i < len(pairs); i++ {
		percentage := float64(pairs[i].count) * 100 / float64(total)
		errwriter.PrintErrf("  %d. %s: %d (%.2f%%)\n", i+1, pairs[i].name, pairs[i].count, percentage)
	}

	errwriter.PrintErrf("Готово: %d строк -> %s за %v\n", total, path, time.Since(start).Round(time.Millisecond))
}

func weightedRandom(rng *rand.Rand, weights []int, totalWeight int) int {
	r := rng.Intn(totalWeight)

	for i, w := range weights {
		if r < w {
			return i
		}
		r -= w
	}

	return 0
}

func verifyCounts(dataFilePath, verificationPath string, progressEvery int64) error {
	expected, err := loadVerificationData(verificationPath)
	if err != nil {
		return fmt.Errorf("не удалось загрузить эталонные данные: %w", err)
	}

	errwriter.PrintErrln("Начинаем верификацию подсчёта имён...")
	errwriter.PrintErrf("Ожидаемое общее количество строк: %d\n", expected.TotalLines)

	f, err := os.Open(dataFilePath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл с данными: %w", err)
	}
	defer safeclose.CloseFile(f, dataFilePath)

	scanner := bufio.NewScanner(f)

	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	// карта для реальных частот
	actualCounts := make(map[string]int64)
	totalLines := int64(0)

	for scanner.Scan() {
		name := scanner.Text()
		actualCounts[name]++
		totalLines++

		// Выводим прогресс каждые n строк по флагу
		if progressEvery > 0 && totalLines%progressEvery == 0 {
			errwriter.PrintErrf("Обработано %d строк...\n", totalLines)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ошибка при чтении файла: %w", err)
	}

	errwriter.PrintErrln("\n=== Результаты верификации ===")
	errwriter.PrintErrf("Всего строк в файле: %d\n", totalLines)
	errwriter.PrintErrf("Ожидалось: %d\n", expected.TotalLines)

	if totalLines != expected.TotalLines {
		errwriter.PrintErrln("ВНИМАНИЕ: Несоответствие количества строк!")
		errwriter.PrintErrf("   Разница: %d строк\n", totalLines-expected.TotalLines)
	} else {
		errwriter.PrintErrln("Количество строк совпадает")
	}

	errwriter.PrintErrln("\n--- Выборочная проверка частот имён ---")

	// примеры имён из разных весовых категорий для проверки, вот так по простому xd
	samples := []string{
		"Александр",
		"Ольга",
		"Артём",
		"Ксения",
		"Борис",
		"Василиса",
		"Трофим",
		"Аристарх",
	}

	// счётчик ошибок
	mismatchCount := 0
	totalChecked := 0

	for _, sample := range samples {
		actual := actualCounts[sample]
		expectedCount, exists := expected.Frequencies[sample]

		totalChecked++

		if !exists {
			errwriter.PrintErrf("%s: имя не найдено в эталонных данных\n", sample)
			mismatchCount++
			continue
		}

		// Сравниваем ожидаемое и реальное значения
		if actual == expectedCount {
			errwriter.PrintErrf("%s: ожидалось %d, реально %d ✓\n", sample, expectedCount, actual)
		} else {
			errwriter.PrintErrf("%s: ожидалось %d, реально %d ✗ (разница: %d)\n",
				sample, expectedCount, actual, actual-expectedCount)
			mismatchCount++
		}
	}

	errwriter.PrintErrln("\n--- Итог выборочной проверки ---")
	if mismatchCount == 0 {
		errwriter.PrintErrln("+ Все проверенные имена совпадают с эталоном")
	} else {
		errwriter.PrintErrf("- Обнаружено расхождений: %d из %d\n", mismatchCount, totalChecked)
	}

	errwriter.PrintErrln("\n--- Топ 10 самых частых имён (реальные данные) ---")

	type countPair struct {
		name  string
		count int64
	}

	// преобразуем карту в слайс для сортировки
	pairs := make([]countPair, 0, len(actualCounts))
	for name, count := range actualCounts {
		pairs = append(pairs, countPair{name, count})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	// Выводим топ 10
	topN := 10
	for i := 0; i < topN && i < len(pairs); i++ {
		percentage := float64(pairs[i].count) * 100 / float64(totalLines)
		errwriter.PrintErrf("%d. %s: %d (%.2f%%)\n", i+1, pairs[i].name, pairs[i].count, percentage)
	}

	errwriter.PrintErrln("\nВерификация завершена")
	return nil
}
