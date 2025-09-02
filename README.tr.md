## Proje Açıklaması
Bu projede, OpenAI tabanlı bir LLM (Large Language Model) backend servisi geliştirildi. Servis iki ana işlev sunar:
- LLM'e mesaj gönderme (chat)
- LLM ile olan konuşma geçmişini session bazlı görüntüleme

Kullanıcılar, bir sessionId ile veya yeni bir session başlatarak LLM'e mesaj gönderebilir. Tüm konuşma geçmişi MySQL veritabanında saklanır ve istenildiğinde görüntülenebilir.


## 📂 Proje Yapısı

Proje **cmd / pkg / internal** standardına göre organize edilmiştir:

```
.
├── cmd/
│   └── myapp/             # main.go (entrypoint)
├── internal/chat/         # domain katmanı
│   ├── handler.go         # HTTP handler'lar
│   ├── handler_test.go    # handler unit testleri
│   ├── service.go         # business logic
│   ├── service_test.go    # service unit testleri
│   ├── repository.go      # MySQL repository (GORM)
│   ├── model.go           # veri modelleri
│   ├── client.go          # OpenAI API client
│   └── mock_*             # gomock ile üretilen mock'lar
├── pkg/
│   ├── config/            # env & config (dotenv ile)
│   ├── database/          # GORM ile MySQL bağlantısı
│   └── logger/            # zap logging
├── .env.example           # örnek environment değişkenleri
├── Makefile               # build & test & run komutları
└── go.mod / go.sum
```

## ⚙️ Kullanılan Teknolojiler ve Kütüphaneler
- **Go** (>=1.22)
- **Echo**: REST API için
- **GORM**: MySQL veritabanı işlemleri için
- **MySQL** (>=8): Konuşma geçmişi ve session yönetimi için
- **OpenAI Go SDK**: LLM API entegrasyonu için
- **Zap**: Loglama için
- **Gomock (uber-go/mock)**: Unit testlerde dışa bağımlılıkları mocklamak için
- **dotenv**: .env dosyasından yapılandırma okumak için

## ⚡ Kurulum ve Çalıştırma

1) Repoyu klonlayın:
    ```bash
    git clone https://github.com/karagultm/llm-chat-service
    cd llm-chat-service
    ```
2) Bağımlılıkları indirin:
    ```bash
    go mod tidy
    ```
3) Gerekli ortam değişkenlerini `.env` dosyasına ekleyin. Örnek için `.env.example` dosyasını kullanabilirsiniz:
	```env
	APP_ENV=dev
	APP_PORT=3000
	DATABASE_URL=your-database-url
	OPENAI_API_KEY=your-api-key-here
	```
4) MySQL veritabanınızı oluşturun ve bağlantı bilgilerini `.env` dosyasına girin.
5) Bağımlılıkları yükleyin:
	```sh
	make mocks
	```
6) Testleri çalıştırın:
	```sh
	make test
	```
7) Uygulamayı başlatın:
	```sh
	make run
	```
> Uygulama başlarken GORM ile MySQL’e bağlanır ve `AutoMigrate` ile tabloyu (yoksa) oluşturmaya çalışır. Bunun çalışabilmesi için `.env` içindeki **DATABASE_URL**’ın doğru olması ve hedef veritabanının (ör. `chatdb`) hazır bulunması gerekir.

## 📡 API Endpointleri
### 1) Chat Gönderme
**POST** `/v1/chat`

Request:
```json
{
  "message": "hello ai",
  "sessionId": "47b2b877-b53d-4bee-877c-585dda3f9e71"
}
```
> `sessionId` opsiyoneldir. Gönderilmezse backend yeni bir session UUID üretir.

Response:
```json
{
  "message": "Hello, how can i help you today?",
  "sessionId": "47b2b877-b53d-4bee-877c-585dda3f9e71"
}
```

### 2) Session Geçmişi
**GET** `/v1/chat/{sessionId}`

Response (örnek):
```json
{
  "sessionId": "47b2b877-b53d-4bee-877c-585dda3f9e71",
  "messages": [
    {
      "id": 1,
      "kind": "USER_PROMPT",
      "message": "hello",
      "timestamp": "1755605815"
    },
    {
      "id": 2,
      "kind": "LLM_OUTPUT",
      "message": "Hello, how can i help you today?",
      "timestamp": "1755605880"
    }
  ]
}
```

---

## 🧪 Testler

- **Handler** ve **Service** katmanları için unit testler yazılmıştır.
- Dış bağımlılıklar (**OpenAI**, **MySQL**) **GoMock** kullanılarak **mock**lanır.
- Çalıştırma:
```bash
make test
# veya
go test ./...
```

---

## 📝 Loglama

- Uygulama, **zap** ile  log üretir.
- İstek/yanıt akışı, hata ve uyarılar detaylı biçimde loglanır.

---

**Not:** Detaylı kabul kriterleri ve senaryolar için `chat-application-backend.md` dosyasına bakabilirsiniz.
