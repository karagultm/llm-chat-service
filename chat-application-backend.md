# LLM Chat Service

## 1. GOAL
You will create an LLM Chat REST service which uses openai chat completions api with keeping history support. 

## 2. ENDPOINTS

### 2.1. Chat
 - Clients can ask questions to LLM and get response with this endpoint.
 - Request has two parameters, `message` and `sessionId`. 
 - `message` parameter contains user's prompt. It's mandatory and should be minimum 3 characters, maximum 2048  characters long. 
 - `sessionId` is being used to keep chat history. This parameter is optional, if it's not given in the request body, backend assumes that it's a new chat session, generates a uuid as session id and sends it in the response.  If it's given in the request, backend assumes that it's an existing chat session, then it stores new message under this session.
 - Response has same parameters, `message` in the response payload contains LLM's response.

**URL:** `POST v1/chat`
**Request Body:** 
```json
{
  "message": "hello ai",
  "sessionId": "47b2b877-b53d-4bee-877c-585dda3f9e71"
}
```

**Response Body:**
```json
{
  "message": "Hello, how can i help you today?",
  "sessionId": "47b2b877-b53d-4bee-877c-585dda3f9e71"
}
```


### 2.2. Session History
 - This endpoint allows clients to retrieve history of a chat session including user and llm messages ordered by timestamp.
 - id parameter is database id of the message object.
 - kind parameter is an Enum, contains two values , `USER_PROMPT` and `LLM_OUTPUT`. 
 - message parameter contains user prompt or llm's response.
 - timestamp parameter is unix timestamp of the message.
**URL:** `GET v1/chat/{sessionId}`

**Response Body:** 
```json
{
  "sessionId": "47b2b877-b53d-4bee-877c-585dda3f9e71",
  "messages": [
    {
      "id": 1,
      "kind": "USER_PROMPT",
      "message": "hello",
      "timestamp": 1755605815
    },
    {
      "id": 1,
      "kind": "LLM_OUTPUT",
      "message": "Hello, how can i help you today",
      "timestamp": 1755605880
    }
  ]
}
```

## 3. Required Tools and Libraries

 - Use go 1.22 or higher version.
 - Use mysql 8 or higher verison as database.
 - Use Echo latest stable version as rest library
 - Use official openai go library: https://github.com/openai/openai-go
 - Use zap library for logging ***
 - Use gomock as mock library https://github.com/uber-go/mock
 - Use dotenv library and .env files to keep&read db and openai parameters.


## 4. Project structure and others
 - Use cmd/pkg/internal project structure
 - Use panic recovery middleware ***
 - Return proper json objects in case of non-successful responses (e.g. 5xx , 4xx) 
 - Add unit tests


## 5. Chat Endpoint Acceptance Scenarios

### 5.1. Happy Path 1 , Without SessionId
1. User sends request with a valid message, without sessionId.
2. Backend generates a uuid, saves incoming message and uuid into messages table with the current timestamp.
3. Backend sends received user message to OpenAI Chat Completions API.
4. Backend receives LLM's response from Chat Completions API, saves it into messages table with same sessionId and current timestamp.
5. Backend sends LLM's response to the user.

### 5.2. Happy Path 2, With SessionId
1. User sends request with a valid message and a valid sessionId.
2. Backend fetches all messages for this session id ordered by timestamp.
3. Backend sends received user message and all previous messages to OpenAI Chat Completions API.
4. Backend receives LLM's response from Chat Completions API, saves it into messages table with same sessionId and current timestamp.
5. Backend sends LLM's response to the user.

### 5.3. Invalid message
1. User sends request with 1 character long message parameter, with or without sessionId.
2. Backend returns error like this `{"error": "invalid message"}` with HTTP 400 status code.

### 5.4. Invalid sessionId
1. Users sends request with a valid message but with an invalid sessionId which doesn't exists in db.
2. Backend returns error like this `{"error": "sessionId not found}`  with HTTP 404 status code.

### 5.5. Other Errors
1. User sends request with a valid message and a valid sessionId.
2. Something goes wrong like OpenAI API connection or database operations.
3. Backend returns error like this: `{"error": "unable to perform chat completion"}` with HTTP 500 status code.


## 6. Session History Endpoint Acceptance Scenarios

### 6.1. Happy Path
1. User sends request with a valid sessionId parameter.
2. Backend fetches all messages for given session id ordered by timestamp.
3. Backend sends all messages to the user.


### 6.2. Invalid sessionId
1. Users sends request with a valid message but with an invalid sessionId which doesn't exists in db.
2. Backend returns error like this `{"error": "sessionId not found}`  with HTTP 404 status code.

### 6.3. Other Errors
1. User sends request with a valid sessionId parameter.
2. Something goes wrong like in the database operations.
3. Backend returns error like this: `{"error": "unable to respond session history"}` with HTTP 500 status code.
