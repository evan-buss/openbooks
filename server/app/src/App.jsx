import React from 'react';
import { Layout, Result, Typography, Spin } from 'antd';

import './App.css';
import Search from './components/Search';
import BookTable from './components/BookTable';
import ServerList from './components/ServerList';
import RecentSearchList from './components/RecentSearchList';
import { countdownTimer, messageRouter, MessageTypes, sendNotification } from "./messages"

const { Header, Sider, Content } = Layout;
const { Title } = Typography;

const titleStyle = { textAlign: "center", color: "white", marginBottom: 0, marginTop: "5px" };

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      items: [],
      searchQueries: [], //Holds individual queries when the user clicks search button
      searchResults: [], //Holds past search results.
      servers: [],
      socket: null,
      loading: false,
      timeLeft: 30
    }
  }

  searchCallback = (loading, queryString) => {
    this.setState((state) => {
      return {
        searchQueries: [...state.searchQueries, queryString],
      }
    });
    this.setState({ loading: true })
  };

  componentDidMount() {
    // Dummy data for testing
    // this.setState({
    //   // items: fakeItems.books,
    //   // searches: recentSearches,
    //   servers: servers.servers
    // });

    //TODO: How do I pass a variable to this...
    let socket = new WebSocket("ws://127.0.0.1:8080/ws");

    socket.onopen = () => {
      console.log("Successfully Connected");
      this.setState({
        socket: socket
      });

      // Send connection request once the socket opens
      this.state.socket.send(JSON.stringify({
        type: MessageTypes.CONNECT,
        payload: {
          //TODO: Make server use the connection name...
          name: "bot"
        }
      }));
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
      let response = JSON.parse(message.data);
      if (response.type === MessageTypes.CONNECT) {
        sendNotification("success", "Successfully Connected", response.status);
        this.setState({ timeLeft: response.wait });
        if (response.wait !== 0) {
          countdownTimer(response.wait, (time) => {
            this.setState({ timeLeft: time })
          });
        }
        this.state.socket.send(JSON.stringify({
          type: MessageTypes.SERVERS,
          payload: {
            data: ""
          }
        }))
      } else {
        this.setState((state) => messageRouter(JSON.parse(message.data), state))
      }
    }
  }

  sidebarStyle = {
    paddingTop: "16px",
    overflow: 'auto',
    height: '100vh',
    position: 'fixed',
    left: 0,
  };

  spinnerStyle = {
    marginTop: "40px"
  }

  contentStyle = {
    marginLeft: "220px",
    marginRight: "220px"
  }

  pastSearchHandler = (index) => {
    this.setState((state) => {
      return {
        items: state.searchResults[index]
      }
    })
  };

  serverButtonHandler = () => {
    console.log("server button handler")
    // Always send a server list request after connection
    this.state.socket.send(JSON.stringify({
      type: MessageTypes.SERVERS,
      payload: {
        data: ""
      }
    }))
  }

  render() {
    return (
      <Layout className="full-size">
        <Sider width={220} style={this.sidebarStyle}>
          <RecentSearchList searches={this.state.searchQueries} clickHandler={this.pastSearchHandler} />
        </Sider>
        <Layout>
          <Header><Title style={titleStyle}>OpenBooks</Title></Header>
          <Content style={this.contentStyle}>
            <div className="app-body">
              <Search socket={this.state.socket}
                disabled={this.state.timeLeft > 0 || this.state.loading}
                searchCallback={this.searchCallback} />
              <button onClick={this.serverButtonHandler}>Refresh Servers</button>
              {/*Show instructions if the user hasn't yet search for anything*/}
              {this.state.searchQueries.length === 0 && <Result title={"Search a book to get started"} />}
              {/*Show loading indicator when waiting for server results*/}
              {(this.state.loading && this.state.items.length > 0) &&
                <Spin size="large" style={this.spinnerStyle} />}
              {/*Show table when items are received and we aren't loading*/}
              {this.state.items.length > 0 &&
                <BookTable
                  items={this.state.items}
                  socket={this.state.socket}
                  disabled={this.state.loading}
                />}
            </div>
          </Content>
        </Layout>
        <ServerList servers={this.state.servers} />
      </Layout>
    );
  }
}

export default App;