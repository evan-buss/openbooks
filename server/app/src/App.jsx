import React from 'react';
import Layout from 'antd/es/layout';
import Title from 'antd/es/typography/Title';
import Spin from 'antd/es/spin';

import './App.css';
import Search from './components/Search';
import BookTable from './components/BookTable';
import ServerList from './components/ServerList';
import RecentSearchList from './components/RecentSearchList';
import PlaceHolder from './components/PlaceHolder';


import { countdownTimer, messageRouter, MessageTypes, sendNotification } from "./messages"

const { Header, Content, Sider } = Layout;

export default class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      items: [],
      searchQueries: [], //Holds individual queries when the user clicks search button
      searchResults: [], //Holds past search results.
      servers: [],
      socket: null,
      loading: false,
      timeLeft: 30,
      selected: -1
    }
  }

  componentDidMount() {
    let localQueries = JSON.parse(window.localStorage.getItem("queries"))
    if (localQueries !== null && localQueries.length !== 0) {
      this.setState({
        searchQueries: localQueries,
        selected: 0
      })
    }

    let localResults = JSON.parse(window.localStorage.getItem("results"));
    console.log(localResults);
    if (localResults !== null && localResults.length !== 0) {
      this.setState({
        searchResults: localResults,
        items: localResults[0]
      })
    }

    let href = window.location.href;

    let addr = `ws://${href.substring(href.indexOf("://") + 3)}ws`
    console.log(addr)
    let socket = new WebSocket(addr);

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

  // This is called when the enters a search query
  searchCallback = (queryString) => {
    console.log("search callback");
    this.setState((state) => {
      let queries = [...state.searchQueries, queryString]
      window.localStorage.setItem("queries", JSON.stringify(queries))
      return {
        searchQueries: queries,
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

  // This is called when a user clicks an item in the search history sidebar
  loadPastSearchHandler = (index) => {
    console.log("past search handler")
    this.setState((state) => {
      return {
        items: state.searchResults[index],
        selected: index
      }
    })
  };

  // This is called when a user clicks the delete button on a search history item
  deletePastSearchHandler = (index) => {
    this.setState((state) => {
      // Remove the item from both the query and searchresults lists
      state.searchQueries.splice(index, 1);
      state.searchResults.splice(index, 1);

      // Update the local storage data
      window.localStorage.setItem("queries", JSON.stringify(state.searchQueries));
      window.localStorage.setItem("results", JSON.stringify(state.searchResults));

      let selected = -1;
      let items = [];
      if (state.searchQueries.length > 0) { // Set active to first entry if list isn't empty
        selected = 0;
        items = state.searchResults[selected]
      }

      return {
        searchQueries: state.searchQueries,
        searchResults: state.searchResults,
        selected: selected,
        items: items
      }
    });
  }

  // This is called when the user clicks on the reload servers button
  reloadServersHandler = () => {
    this.state.socket != null && this.state.socket.send(JSON.stringify({
      type: MessageTypes.SERVERS,
      payload: {}
    }))
  }

  // This is called when the user clicks a download link on a book
  bookDownloadHandler = (bookString) => {
    if (this.state.socket != null) {
      this.state.socket.send(JSON.stringify({
        type: MessageTypes.DOWNLOAD,
        payload: {
          book: bookString
        }
      }))

      this.setState({ loading: true })
    }
  }

  render() {
    return (
      <Layout>
        <Sider width={240} className="sider-style">
          <RecentSearchList
            disabled={this.state.loading}
            searches={this.state.searchQueries}
            selected={this.state.selected}
            clickHandler={this.loadPastSearchHandler}
            deleteHandler={this.deletePastSearchHandler}
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
                disabled={this.state.loading || this.state.timeLeft > 0}
                downloadCallback={this.bookDownloadHandler}
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

const spinnerStyle = {
  marginTop: 40
}