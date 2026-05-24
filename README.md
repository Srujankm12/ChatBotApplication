# LLM Chatbot with Inference Logging

A full-stack chatbot application built using Next.js, Go (Echo), MongoDB, and Google Gemini, with a custom SDK layer for inference logging and observability.

---

## Overview

This project demonstrates how to build a lightweight LLM system with:

* Multi-turn chatbot
* SDK wrapper for LLM calls
* Real-time inference logging
* Ingestion pipeline for metadata storage
* Clean and scalable backend architecture

---

## Architecture

```
Frontend (Next.js)
    ↓
Backend API (Go + Echo)
    ↓
SDK Layer (LLM Wrapper)
    ↓
Gemini API

Parallel flow:
SDK → /logs → MongoDB
```

---

## Features

* Multi-turn conversations (last 5 messages as context)
* Session-based chat using UUID
* Chat persistence (restores after refresh)
* Conversation history (switch between sessions)
* Cancel / New chat support
* Inference logging (latency, token usage, status)
* Metrics endpoint for observability

---

## SDK Layer

A lightweight wrapper around LLM calls.

### Responsibilities

* Calls Gemini API
* Builds prompt from chat history
* Measures latency
* Handles retry (one retry with delay)
* Sends logs asynchronously

### Captured Metadata

* model
* provider
* latency
* token_usage
* input/output preview
* timestamp
* session_id
* status

---

## Ingestion Pipeline

### Endpoint

POST /logs

### Behavior

* Receives logs from SDK
* Validates payload
* Stores metadata in MongoDB

---

## Database Design

### chats

```json
{
  "session_id": "string",
  "messages": [
    {
      "role": "user | assistant",
      "content": "string",
      "timestamp": "ISODate"
    }
  ]
}
```

### logs

```json
{
  "model": "string",
  "provider": "string",
  "latency_ms": "number",
  "input_preview": "string",
  "output_preview": "string",
  "token_usage": "number",
  "timestamp": "ISODate",
  "session_id": "string",
  "status": "success | error"
}
```

---

## Tech Stack

* Frontend: Next.js, Tailwind CSS
* Backend: Go (Echo)
* Database: MongoDB
* LLM: Google Gemini

---

## Running Locally

### Backend

```bash
cd backend
go run main.go
```

Runs on:
https://chatbotbackendd.onrender.com

---

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Runs on:
http://localhost:3000

---

## Docker Setup

```bash
docker-compose up --build
```

* Frontend: http://localhost:3000
* Backend: https://chatbotbackendd.onrender.com

---

## API Endpoints

### POST /chat

```json
{
  "message": "Hello",
  "session_id": "uuid"
}
```

### GET /chat/:session_id

Fetch conversation history.

### GET /sessions

List recent sessions.

### GET /metrics

```json
{
  "total_requests": number,
  "avg_latency": number,
  "error_rate": number
}
```

---

## Demo

Suggested demonstration flow:

* Chat interaction with multiple turns
* Refresh to show persistence
* Switching between sessions
* Viewing logs in MongoDB
* Calling /metrics endpoint

---

## Key Design Decisions

* Last 5 messages context
  Keeps token usage controlled and improves performance

* Asynchronous logging
  Prevents logging from affecting response latency

* Separate ingestion endpoint
  Improves scalability and separation of concerns

* Token estimation instead of exact usage
  Avoids dependency on provider-specific APIs

---

## Tradeoffs

* No streaming responses implemented
  Keeps system simpler and easier to maintain

* No dashboard UI
  Focuses on backend observability

* Token usage is estimated
  Faster but less accurate

---

## Schema Design Decisions

* Chats stored per session_id
  Simplifies conversation retrieval

* Messages stored as an array in a single document
  Efficient for fetching recent context

* Logs stored in a separate collection
  Separates operational data from user data

* Only preview data stored
  Reduces storage size and improves performance

---

## Architecture Notes

### Ingestion Flow

SDK → /logs → MongoDB

### Logging Strategy

Logs are sent asynchronously using goroutines to avoid blocking chat responses.

### Scaling Considerations

* Backend is stateless
* Can scale horizontally
* MongoDB handles persistence

### Failure Handling

* Retry logic for LLM calls
* Graceful error responses
* Errors stored in logs collection

---

## What I Would Improve With More Time

* Add streaming responses for better user experience
* Build a dashboard for metrics visualization
* Support multiple LLM providers beyond Gemini
* Add rate limiting and validation
* Improve token usage accuracy using provider APIs
* Add authentication and user-based sessions
