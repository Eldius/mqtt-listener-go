
import React from "react"

import './App.css';

import NetworkMonitor from "./components/network/NetworkMonitor"

function App() {
  return (
    <div className="App">
      <h1>MQTT Listener</h1>
      <NetworkMonitor />
    </div>
  );
}

export default App;
