.PHONY: all build run test test-race bench vet tidy clean docker-build docker-run docker-pprof \
	testdata-large clean-testdata-large help

BINARY := namecounter
CMD := ./cmd/run
SAMPLE := testdata/sample.txt
LARGE := testdata/large_names.txt
# Override: LINES=20000000 make testdata-large
LINES ?= 10000000

all: vet test build

build:
	go build -o $(BINARY) $(CMD)

run: build
	./$(BINARY) $(SAMPLE)

test:
	go test -count=1 ./...

test-race:
	go test -race -count=1 ./...

bench:
	go test -bench=. -benchmem ./internal/domain/...

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -f $(BINARY) $(BINARY).exe

docker-build:
	docker compose build

docker-run:
	docker compose run --rm namecounter

docker-pprof:
	docker compose --profile pprof up namecounter-pprof

# по дефолту 10M строк (+-100+ Мб).
testdata-large:
	go run ./tools/genlarge -out $(LARGE) -n $(LINES)

clean-testdata-large:
	rm -f $(LARGE)

help:
	@echo "Targets:"
	@echo "  make build        - go build -> ./$(BINARY)"
	@echo "  make run          - build and count $(SAMPLE)"
	@echo "  make test         - go test ./..."
	@echo "  make test-race    - go test -race ./..."
	@echo "  make bench        - benchmarks (internal/domain)"
	@echo "  make vet          - go vet ./..."
	@echo "  make tidy         - go mod tidy"
	@echo "  make clean        - remove built binary"
	@echo "  make docker-build      - docker compose build"
	@echo "  make docker-run        - docker compose run --rm namecounter"
	@echo "  make docker-pprof      - compose up namecounter-pprof (profile pprof)"
	@echo "  make testdata-large    - generate $(LARGE) ($(LINES) lines; set LINES=...)"
	@echo "  make clean-testdata-large - remove $(LARGE)"
