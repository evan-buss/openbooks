import React from 'react';
import './App.css';
import Search from './components/Search';
import Table from './components/Table';
import { MessageTypes, messageRouter } from "./messages"

class App extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      name: "",
      items: [],
      socket: null
    }
  }

  componentDidMount() {
    // Connect websocket here
    let socket = new WebSocket("ws://127.0.0.1:8080/ws");

    socket.onopen = () => {
      console.log("Successfully Connected");
      this.setState({
        socket: socket
      })
    };

    socket.onclose = event => {
      console.log("Socket Closed Connection: ", event);
      this.setState({
        socket: null
      })
    };

    socket.onerror = error => {
      console.log("Socket Error: ", error);
      this.setState({
        socket: null
      })
    };

    socket.onmessage = message => {
      messageRouter(JSON.parse(message.data))
      // this.setState({
      //   name: message.data
      // })
    }

  }

  sendMessage = () => {
    this.state.socket.send(JSON.stringify({
      "type": MessageTypes.CONNECT,
      "payload": {
        "value": 20
      }
    }))
  }

  render() {
    return (
      <div className="app-body">
        <button onClick={this.sendMessage}>Test Button</button>
        <div id="app-container">
          <p>{this.state.name}</p>
          <Search />
          <Table />
        </div>
      </div>
    );
  }
}

export default App;