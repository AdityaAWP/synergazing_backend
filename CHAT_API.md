# Chat API Documentation

## Overview

The Chat API provides real-time messaging functionality between users using WebSockets and REST endpoints.

## WebSocket Connection

### Connect to WebSocket

```
ws://localhost:3002/ws/chat?user_id={user_id}
```

**Parameters:**

- `user_id`: The ID of the user connecting

### WebSocket Message Types

#### Client to Server Messages

1. **Join Chat**

```json
{
  "type": "join_chat",
  "chat_id": 1
}
```

2. **Send Message**

```json
{
  "type": "send_message",
  "chat_id": 1,
  "content": "Hello, how are you?"
}
```

3. **Mark Messages as Read**

```json
{
  "type": "mark_read",
  "chat_id": 1
}
```

#### Server to Client Messages

1. **Connection Confirmation**

```json
{
  "type": "connected",
  "data": {
    "message": "Connected to chat server"
  }
}
```

2. **New Message**

```json
{
  "type": "new_message",
  "data": {
    "id": 1,
    "chat_id": 1,
    "sender_id": 2,
    "content": "Hello!",
    "is_read": false,
    "created_at": "2025-08-07T10:30:00Z",
    "sender": {
      "id": 2,
      "name": "John Doe"
    }
  }
}
```

3. **Chat Joined**

```json
{
  "type": "joined_chat",
  "data": {
    "chat_id": 1,
    "message": "Joined chat successfully"
  }
}
```

4. **Error**

```json
{
  "type": "error",
  "data": {
    "error": "Error message here"
  }
}
```

## REST API Endpoints

All REST endpoints require authentication via JWT token in the Authorization header:

```
Authorization: Bearer {jwt_token}
```

### 1. Get or Create Chat with Another User

```
GET /api/chat/with/{user_id}
```

**Response:**

```json
{
  "success": true,
  "message": "Chat retrieved successfully",
  "data": {
    "id": 1,
    "user1_id": 1,
    "user2_id": 2,
    "user1": {
      "id": 1,
      "name": "Alice",
      "email": "alice@example.com"
    },
    "user2": {
      "id": 2,
      "name": "Bob",
      "email": "bob@example.com"
    },
    "created_at": "2025-08-07T10:00:00Z",
    "updated_at": "2025-08-07T10:30:00Z"
  }
}
```

### 2. Get All User Chats

```
GET /api/chat/
```

**Response:**

```json
{
    "success": true,
    "message": "Chats retrieved successfully",
    "data": [
        {
            "id": 1,
            "user1_id": 1,
            "user2_id": 2,
            "user1": {...},
            "user2": {...},
            "messages": [
                {
                    "id": 5,
                    "content": "Latest message",
                    "created_at": "2025-08-07T10:30:00Z"
                }
            ],
            "created_at": "2025-08-07T10:00:00Z",
            "updated_at": "2025-08-07T10:30:00Z"
        }
    ]
}
```

### 3. Get Chat Messages

```
GET /api/chat/{chat_id}/messages?page=1&limit=50
```

**Query Parameters:**

- `page`: Page number (default: 1)
- `limit`: Messages per page (default: 50, max: 100)

**Response:**

```json
{
  "success": true,
  "message": "Messages retrieved successfully",
  "data": {
    "messages": [
      {
        "id": 1,
        "chat_id": 1,
        "sender_id": 1,
        "content": "Hello!",
        "is_read": true,
        "sender": {
          "id": 1,
          "name": "Alice"
        },
        "created_at": "2025-08-07T10:25:00Z",
        "updated_at": "2025-08-07T10:25:00Z"
      }
    ],
    "page": 1,
    "limit": 50
  }
}
```

### 4. Mark Chat Messages as Read

```
PUT /api/chat/{chat_id}/read
```

**Response:**

```json
{
  "success": true,
  "message": "Messages marked as read",
  "data": null
}
```

## Database Schema

### Chats Table

```sql
CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    user1_id INTEGER NOT NULL REFERENCES users(id),
    user2_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Messages Table

```sql
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    chat_id INTEGER NOT NULL REFERENCES chats(id),
    sender_id INTEGER NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Testing

A test HTML file is provided at `/storage/chat-test.html` that you can open in your browser to test the WebSocket functionality.

1. Run your Go server
2. Open `http://localhost:3002/storage/chat-test.html` in your browser
3. Enter different User IDs for each chat window
4. Use the same Chat ID for both users to test real-time messaging

## Error Handling

Common error scenarios:

- **Unauthorized access**: User trying to access a chat they're not part of
- **Invalid chat ID**: Non-existent chat ID
- **Empty message**: Attempting to send empty message
- **Self-chat**: Trying to create chat with yourself
- **Connection failures**: WebSocket connection issues

## Security Considerations

1. **Authentication**: REST endpoints are protected by JWT middleware
2. **Authorization**: Users can only access chats they're part of
3. **Data validation**: Input validation on all endpoints
4. **WebSocket security**: User ID verification on WebSocket connection

## Performance Considerations

1. **Pagination**: Messages are paginated to avoid large data transfers
2. **Connection management**: WebSocket connections are properly managed and cleaned up
3. **Database indexes**: Consider adding indexes on frequently queried fields
4. **Message limits**: Consider implementing message history limits for performance
