// src/App.js
import React, { useState, useEffect, useRef } from "react";
import "./App.css";

const App = () => {
  const [conn, setConn] = useState(null);
  const [msg, setMsg] = useState("");
  const logRef = useRef();

  const senders = ["Alice", "Bob", "Charlie", "Dave", "Eve", "Grace", "Mallory"];

  const getRandomSender = () => {
    const randomIndex = Math.floor(Math.random() * senders.length);
    return senders[randomIndex];
  };

  const appendLog = (text) => {
    const log = logRef.current;
    if (log) {
      const doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
      const item = document.createElement("div");
      item.innerText = text;
      log.appendChild(item);
      if (doScroll) {
        log.scrollTop = log.scrollHeight - log.clientHeight;
      }
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!conn || !msg) return;

    const sender = getRandomSender();
    const messageObj = {
      sender: sender,
      content: msg,
    };

    conn.send(JSON.stringify(messageObj));
    setMsg("");
  };

  useEffect(() => {
    if ("WebSocket" in window) {
      const websocket = new WebSocket(`ws://${window.location.hostname}:8080/ws`);
      setConn(websocket);

      websocket.onclose = () => appendLog("Connection closed.");
      websocket.onmessage = (evt) => {
        const messages = evt.data.split("\n");
        messages.forEach((message) => {
          const msgObj = JSON.parse(message);
          appendLog(`${msgObj.sender}: ${msgObj.content}`);
        });
      };

      return () => websocket.close();
    } else {
      appendLog("Your browser does not support WebSockets.");
    }
  }, []);

  return (
    <div>
      <div id="log" ref={logRef} style={logStyle}></div>
      <form id="form" onSubmit={handleSubmit} style={formStyle}>
        <input
          type="text"
          id="msg"
          value={msg}
          onChange={(e) => setMsg(e.target.value)}
          size="64"
          autoFocus
        />
        <input type="submit" value="Send" />
      </form>
    </div>
  );
};

const logStyle = {
  background: "white",
  margin: "0",
  padding: "0.5em",
  position: "absolute",
  top: "0.5em",
  left: "0.5em",
  right: "0.5em",
  bottom: "3em",
  overflow: "auto",
};

const formStyle = {
  padding: "0 0.5em",
  margin: "0",
  position: "absolute",
  bottom: "1em",
  left: "0px",
  width: "100%",
  overflow: "hidden",
};

export default App;