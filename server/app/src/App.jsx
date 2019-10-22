import React from 'react';
import {countdownTimer, messageRouter, MessageTypes, sendNotification} from "./messages"
import {Layout, Result, Typography} from 'antd';

import './App.css';
import Search from './components/Search';
import BookTable from './components/BookTable';
import ServerList from './components/ServerList';
import RecentSearchList from './components/RecentSearchList';

const {Header, Sider, Content} = Layout;
const {Title} = Typography;

const titleStyle = {textAlign: "center", color: "white", marginBottom: 0, marginTop: "5px"};

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
    this.setState({loading: true})
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
        this.setState({timeLeft: response.wait});
        if (response.wait !== 0) {
          countdownTimer(response.wait, (time) => {
            this.setState({timeLeft: time})
          });
        } else {
          this.state.socket.send(JSON.stringify({
            type: MessageTypes.SERVERS,
            payload: {
              data: ""
            }
          }))
        }
      } else {
        this.setState((state) => messageRouter(JSON.parse(message.data), state))
      }
    }
  }

  sidebarStyle = {
    paddingTop: "16px",
    paddingLeft: "8px",
  };

  pastSearchHandler = (index) => {
    this.setState((state) => {
      return {
        items: state.searchResults[index]
      }
    })
  };

  render() {
    return (
      <Layout className="full-size">
        <Sider style={this.sidebarStyle}>
          <RecentSearchList searches={this.state.searchQueries} clickHandler={this.pastSearchHandler}/>
          <ServerList servers={this.state.servers}/>
        </Sider>
        <Layout>
          <Header><Title style={titleStyle}>OpenBooks</Title></Header>
          <Content>
            <div className="app-body">
              <Search socket={this.state.socket}
                      disabled={this.state.timeLeft > 0 || this.state.loading}
                      searchCallback={this.searchCallback}/>
              {/*Show instructions if the user hasn't yet search for anything*/}
              {this.state.searchQueries.length === 0 && <Result title={"Search a book to get started"}/>}
              {/*Show loading indicator when waiting for server results*/}
              {(this.state.loading && this.state.items.length > 0) &&
              <Result title={"Loading your results. Please wait."}/>}
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
      </Layout>
    );
  }
}

export default App;