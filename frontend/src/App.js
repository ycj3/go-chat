import React, { useEffect, useState } from 'react';
import { chat } from './message_pb'; // Assuming chat is imported from somewhere

function App() {
  const [ws, setWs] = useState(null);
  const [messages, setMessages] = useState([]);
  const [user, setUser] = useState("");
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [msg, setMsg] = useState("");
  const [onlineCount, setOnlineCount] = useState(0);
  const [onlineMembers, setOnlineMembers] = useState([]);

  const handleLogin = (e) => {
    e.preventDefault();
    setIsLoggedIn(true);
  };

  const sendMessage = (msg) => {
    const sender = user;
    const messageObj = chat.ChatMessage.create({ user: sender, message: msg });
    const msgBuffer = chat.ChatMessage.encode(messageObj).finish();
    ws.send(msgBuffer);
    setMsg("");
  };

  useEffect(() => {
    let isConnected = false;

    if (isLoggedIn && "WebSocket" in window) {
      const websocket = new WebSocket(`ws://${window.location.hostname}:8080/ws?user=${user}`);
      websocket.binaryType = "arraybuffer";

      setWs(websocket);

      websocket.onclose = () => {
        console.log("WebSocket connection closed.");
      };

      websocket.onopen = () => {
        console.log("WebSocket connection established.");
        isConnected = true;
      };

      websocket.onmessage = (evt) => {
        const buffer = new Uint8Array(evt.data);
        const msgObj = chat.ChatMessage.decode(buffer);
        setMessages((prevMessages) => [...prevMessages, `${msgObj.user}: ${msgObj.message}`]);
      };

      return () => {
        if (isConnected && websocket.readyState === WebSocket.OPEN) {
          websocket.close();
        }
      };
    } else if (!isLoggedIn) {
      console.log("User is not logged in.");
    } else {
      console.log("Your browser does not support WebSockets.");
    }
  }, [isLoggedIn, user]);

  useEffect(() => {
    if (isLoggedIn) {
      fetch(`http://${window.location.hostname}:8080/online`)
        .then((response) => response.json())
        .then((data) => {
          setOnlineCount(data.count);
          setOnlineMembers(data.members);
        });
    }
  }, [isLoggedIn]);

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
                <li key={index}>{member}</li>
              ))}
            </ul>
          </div>
          <div>
            {messages.map((msg, index) => (
              <div key={index}>{msg}</div>
            ))}
          </div>
          <input
            type="text"
            value={msg}
            onChange={(e) => setMsg(e.target.value)}
            placeholder="Type a message"
          />
          <button onClick={() => sendMessage(msg)}>Send Message</button>
        </div>
      )}
    </div>
  );
}

export default App;