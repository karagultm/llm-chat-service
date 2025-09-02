# LLM Chat Backend Service
## Project Description
In this project, a backend service based on OpenAI's LLM (Large Language Model) was developed. The service provides two main functions:
- Sending a message to the LLM (chat)
- Viewing the conversation history with the LLM on a session basis

Users can send messages to the LLM by providing a sessionId or by starting a new session. All conversation history is stored in a MySQL database and can be retrieved when needed.

## 📂 Project Structure

The project is organized according to the **cmd / pkg / internal** standard:

```
.
├── cmd/
│   └── myapp/             # main.go (entrypoint)
├── internal/chat/         # domain layer
│   ├── handler.go         # HTTP handlers
│   ├── handler_test.go    # handler unit tests
│   ├── service.go         # business logic
│   ├── service_test.go    # service unit tests
│   ├── repository.go      # MySQL repository (GORM)
│   ├── model.go           # data models
│   ├── client.go          # OpenAI API client
│   └── mock_*             # gomock generated mocks
├── pkg/
│   ├── config/            # env & config (dotenv)
│   ├── database/          # MySQL connection with GORM
│   └── logger/            # zap logging
├── .env.example           # sample environment variables
├── Makefile               # build & test & run commands
└── go.mod / go.sum
```

## ⚙️ Technologies and Libraries Used
- **Go** (>=1.22)
- **Echo**: for REST API
- **GORM**: for MySQL database operations
- **MySQL** (>=8): for conversation history and session management
- **OpenAI Go SDK**: for LLM API integration
- **Zap**: for logging
- **Gomock (uber-go/mock)**: for mocking external dependencies in unit tests
- **dotenv**: for reading configuration from `.env`

## ⚡ Setup and Run

1) Clone the repo:
    ```bash
    git clone https://github.com/karagultm/llm-chat-service
    cd llm-chat-service
    ```
2) Install dependencies:
    ```bash
    go mod tidy
    ```
3) Add the required environment variables to the `.env` file. Use `.env.example` as a reference:
	```env
	APP_ENV=dev
	APP_PORT=3000
	DATABASE_URL=your-database-url
	OPENAI_API_KEY=your-api-key-here
	```
4) Create your MySQL database and enter the connection details in the `.env` file.
5) Load dependencies:
	```sh
	make mocks
	```
6) Run the tests:
	```sh
	make test
	```
7) Start the application:
	```sh
	make run
	```
> On startup, the application connects to MySQL using GORM and tries to create the table (if it doesn’t exist) via `AutoMigrate`. For this to work, the **DATABASE_URL** in `.env` must be correct and the target database (e.g., `chatdb`) must be available.

## 📡 API Endpoints
### 1) Send Chat
**POST** `/v1/chat`

Request:
```json
{
  "message": "hello ai",
  "sessionId": "47b2b877-b53d-4bee-877c-585dda3f9e71"
}
```
> `sessionId` is optional. If not provided, the backend generates a new session UUID.

Response:
```json
{
  "message": "Hello, how can i help you today?",
  "sessionId": "47b2b877-b53d-4bee-877c-585dda3f9e71"
}
```

### 2) Session History
**GET** `/v1/chat/{sessionId}`

Response (example):
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

## 🧪 Tests

- Unit tests are written for both **Handler** and **Service** layers.
- External dependencies (**OpenAI**, **MySQL**) are mocked using **GoMock**.
- Run tests:
```bash
make test
# or
go test ./...
```

---

## 📝 Logging

- The application logs with **zap**.
- Request/response flow, errors, and warnings are logged in detail.

---

**Note:** For detailed acceptance criteria and scenarios, see the `chat-application-backend.md` file.

---

👉 For the Turkish version of this README, please check [README.tr.md](README.tr.md).
