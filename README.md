# go-chat

## Overview

`go-chat` is a simple chat application built with Go for the backend and React for the frontend. It uses WebSockets for real-time communication between clients.

## Getting Started

### Prerequisites

- Go 1.22.5 or later
- Node.js 14.x or later
- npm 6.x or later

### Backend Setup

1. Clone the repository:

    ```sh
    git clone https://github.com/ycj3/go-chat.git
    cd go-chat
    ```

2. Install Go dependencies:

    ```sh
    go get
    ```

3. Run the backend server:

    ```sh
    go run main.go
    ```

    The backend server will start on `http://localhost:8080`.

### Frontend Setup

1. Navigate to the frontend directory:

    ```sh
    cd frontend
    ```

2. Install Node.js dependencies:

    ```sh
    npm install
    ```

3. Start the frontend development server:

    ```sh
    npm start
    ```

    The frontend server will start on `http://localhost:3000`.

## Usage

1. Open your browser and navigate to `http://localhost:3000`.
2. Enter your username and click "Login".
3. You can now send and receive messages in real-time.

## API Endpoints

### WebSocket Endpoint

- : `ws://localhost8080/ws?user=<username>`

### REST Endpoint

- `GET /online`: Returns the number of online users and the list of online members.

Example response:

```json
{
  "count": 2,
  "members": ["user1", "user2"]
}
```

## Project Structure

### Backend

- main.go: Entry point for the backend server.
- websocket: Contains WebSocket-related code.
  - `client.go`: Handles WebSocket client connections.
  - `hub.go`: Manages the set of active clients and broadcasts messages.
- proto: Contains Protocol Buffers definitions.
  - `message.proto`: Protocol Buffers schema for chat messages.
- pb: Contains generated Protocol Buffers code.
  - `message.pb.go`: Generated Go code from `message.proto`.

### Frontend

- `src/`: Contains React application code.
  - `App.js`: Main React component.
  - `message_pb.js`: Generated JavaScript code from `message.proto`.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
