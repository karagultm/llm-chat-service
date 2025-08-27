# Makefile

GO = go

.PHONY: all test mocks run

# Tüm testleri çalıştır
test:
	$(GO) test ./... -v

# Mockları yeniden oluştur
mocks: internal/chat/repository.go internal/chat/client.go
	@echo "Generating mocks..."
	mockgen -source=internal/chat/repository.go -destination=internal/chat/mock_repository.go -package=chat
	mockgen -source=internal/chat/client.go -destination=internal/chat/mock_client.go -package=chat

# Projeyi çalıştır
run:
	$(GO) run ./cmd/myapp/main.go

# Tüm adımları birleştir
all: mocks test run
