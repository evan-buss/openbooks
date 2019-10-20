import React from 'react';
import './App.css';
import Search from './components/Search';
import Table from './components/Table';
import { messageRouter } from "./messages"
import GridLoader from 'react-spinners/GridLoader'

class App extends React.Component {
  constructor(props) {
    super(props);

    // NOTE: Set state via function if you need the existing state to calculate new state
    this.state = {
      username: "",
      searchString: "",
      items: [],
      servers: [],
      socket: null,
      loading: false
    }
  }

  loadingCallback = (bool) => {
    this.setState({ loading: bool })
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
      this.setState(messageRouter(JSON.parse(message.data)))
    }

  }

  render() {
    return (
      <div className="app-body">
        <div id="app-container">
          <p>{this.state.username}</p>
          <Search socket={this.state.socket}
            loadingCallback={this.loadingCallback} />
          {this.state.loading ? (
            <GridLoader
              color="#09d3ac"
              css={
                `size: 100;`
              }
            />
          ) : (
              <Table items={this.state.items} socket={this.state.socket} />
            )}
        </div>
      </div>
    );
  }
}

export default App;