import React, { useEffect, useState } from 'react';
import { chat } from './pb/message_pb'; // Assuming chat is imported from somewhere

function App() {
  const [ws, setWs] = useState(null);
  const [messages, setMessages] = useState([]);
  const [user, setUser] = useState("");
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [msg, setMsg] = useState("");
  const [onlineCount, setOnlineCount] = useState(0);
  const [onlineMembers, setOnlineMembers] = useState([]);

  useEffect(() => {
    let heartbeatInterval;
    if (ws) {
      heartbeatInterval = setInterval(() => {
        if (ws.readyState === WebSocket.OPEN) {
          ws.send(JSON.stringify({ type: 'heartbeat', user_id: user }));
        }
      }, 3000); 
    }
    return () => clearInterval(heartbeatInterval);
  }, [ws, user]);

  const handleLogin = (e) => {
    e.preventDefault();
    setIsLoggedIn(true);

    if ("WebSocket" in window) {
      const websocket = new WebSocket(`ws://${window.location.hostname}:8080/ws?user_id=${user}`);
      websocket.binaryType = "arraybuffer";

      setWs(websocket);

      websocket.onclose = () => {
        console.log("WebSocket connection closed.");
      };

      websocket.onopen = () => {
        console.log("WebSocket connection established.");
        fetch(`http://${window.location.hostname}:8080/online`)
          .then((response) => response.json())
          .then((data) => {
            setOnlineCount(data.count);
            setOnlineMembers(data.members);
          });
      };

      websocket.onmessage = (evt) => {
        const buffer = new Uint8Array(evt.data);
        const msgObj = chat.ChatMessage.decode(buffer);
        setMessages((prevMessages) => [...prevMessages, `${msgObj.user}: ${msgObj.message}`]);
      };

      websocket.onerror = (error) => {
        console.error("WebSocket error:", error);
      };
    }
  };

  const sendMessage = (msg) => {
    const sender = user;
    const messageObj = chat.ChatMessage.create({ user: sender, message: msg });
    const msgBuffer = chat.ChatMessage.encode(messageObj).finish();
    ws.send(msgBuffer);
    setMsg("");
  };

  return (
    <div>
      {!isLoggedIn ? (
        <form onSubmit={handleLogin}>
          <input
            type="text"
            placeholder="Enter your username"
            value={user}
            onChange={(e) => setUser(e.target.value)}
            required
          />
          <button type="submit">Login</button>
        </form>
      ) : (
        <div>
          <h1>WebSocket Chat</h1>
          <div>
            <h2>Online Users: {onlineCount}</h2>
            <ul>
              {onlineMembers.map((member, index) => (
                <li key={index}>{member.nickname}</li>
              ))}
            </ul>
          </div>
          <div>
            <input
              type="text"
              value={msg}
              onChange={(e) => setMsg(e.target.value)}
              placeholder="Enter your message"
            />
            <button onClick={() => sendMessage(msg)}>Send</button>
          </div>
          <div>
            <h2>Messages</h2>
            <ul>
              {messages.map((message, index) => (
                <li key={index}>{message}</li>
              ))}
            </ul>
          </div>
        </div>
      )}
    </div>
  );
}

export default App;