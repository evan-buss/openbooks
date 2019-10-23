import React from 'react';
import Layout from 'antd/es/layout';
import Typography from 'antd/es/typography';
import Spin from 'antd/es/spin';

import './App.css';
import Search from './components/Search';
import BookTable from './components/BookTable';
import ServerList from './components/ServerList';
import RecentSearchList from './components/RecentSearchList';
import PlaceHolder from './components/PlaceHolder';

import { countdownTimer, messageRouter, MessageTypes, sendNotification } from "./messages"

import { fakeItems } from "./dummyData";


const { Header, Content, Sider } = Layout;
const { Title } = Typography;

export default class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      items: [],
      searchQueries: ["this is a very long string title"], //Holds individual queries when the user clicks search button
      searchResults: [], //Holds past search results.
      servers: [],
      socket: null,
      loading: false,
      timeLeft: 30,
      selected: 0
      // selected: -1
    }
  }

  componentDidMount() {
    // Dummy data for testing
    this.setState({
      items: fakeItems.books,
      // searches: recentSearches,
      // servers: servers.servers
    });

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
        this.reloadServersHandler();
      } else {
        this.setState((state) => messageRouter(JSON.parse(message.data), state))
      }
    }
  }


  searchCallback = (queryString) => {
    this.setState((state) => {
      return {
        searchQueries: [...state.searchQueries, queryString],
        selected: state.searchQueries.length,
        loading: true
      }
    });

    this.state.socket != null &&
      this.state.socket.send(JSON.stringify(
        {
          type: MessageTypes.SEARCH,
          payload: {
            query: queryString
          }
        }
      ))
  };


  pastSearchHandler = (index) => {
    this.setState((state) => {
      return {
        items: state.searchResults[index],
        selected: index
      }
    })
  };

  reloadServersHandler = () => {
    console.log("reload")
    this.state.socket != null && this.state.socket.send(JSON.stringify({
      type: MessageTypes.SERVERS,
      payload: {}
    }))
  }

  render() {
    return (
      <Layout>
        <Sider width={240} style={siderStyle}>
          <RecentSearchList
            searches={this.state.searchQueries}
            selected={this.state.selected}
            clickHandler={this.pastSearchHandler}
          />
          <ServerList servers={this.state.servers} reloadCallback={this.reloadServersHandler} />
        </Sider>
        <Layout>
          <Header style={headerStyle}>
            <Title style={titleStyle}>OpenBooks</Title>
          </Header>
          <Content style={contentStyle}>
            <Search disabled={this.state.timeLeft > 0 || this.state.loading}
              searchCallback={this.searchCallback} />
            {/*Show instructions if the user hasn't yet search for anything*/}
            {this.state.searchQueries.length === 0 &&
              <PlaceHolder />}
            {/*Show loading indicator when waiting for server results*/}
            {(this.state.loading && this.state.items.length === 0) &&
              <Spin size="large" style={spinnerStyle} />}
            {/*Show table when items are received and we aren't loading*/}
            {this.state.items.length > 0 &&
              <BookTable
                items={this.state.items}
                socket={this.state.socket}
                disabled={this.state.loading}
              />}
          </Content>
        </Layout>
      </Layout>
    );
  }
}

const headerStyle = {
  position: "fixed",
  width: "100vw"
}

const titleStyle = {
  textAlign: "center",
  color: "white",
  marginBottom: 0,
  marginTop: 5,
  marginLeft: 220
};

const contentStyle = {
  marginTop: 64,
  overflow: "auto",
  marginLeft: 240,
  height: "calc(100vh - 64px)",
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
  fontSize: "calc(10px + 2vmin)",
  color: "white"
}

const siderStyle = {
  overflow: 'auto',
  height: '100vh',
  position: 'fixed',
  left: 0,
  zIndex: 5
}

const spinnerStyle = {
  marginTop: 40
}