// src/App.js
import React, { useState, useEffect, useRef } from "react";
import "./App.css";
import { chat } from './message_pb'; // Path to your generated file

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
    const messageObj = chat.ChatMessage.create({user: sender, message: msg})
    const msgBuffer = chat.ChatMessage.encode(messageObj).finish();
    conn.send(msgBuffer);
    setMsg("");
  };

  useEffect(() => {
    let isConnected = false;

    if ("WebSocket" in window) {
      const websocket = new WebSocket(`ws://${window.location.hostname}:8080/ws`);
      websocket.binaryType = "arraybuffer";

      setConn(websocket);

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
        appendLog(`${msgObj.user}: ${msgObj.message}`);
      };

      return () => { 
        if (isConnected && websocket.readyState === WebSocket.OPEN) {
          websocket.close();
        }
      };
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