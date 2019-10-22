import React from 'react';
import {messageRouter, MessageTypes} from "./messages"
import {Layout, notification, Result, Typography} from 'antd';

import './App.css';
import Search from './components/Search';
import BookTable from './components/BookTable';
import ServerList from './components/ServerList';
import RecentSearchList from './components/RecentSearchList';
import {servers} from "./dummyData";

const {Header, Sider, Content} = Layout;
const {Title} = Typography;

const titleStyle = {textAlign: "center", color: "white", marginBottom: 0, marginTop: "5px"};

class App extends React.Component {
    constructor(props) {
        super(props);

        // NOTE: Set state via function if you need the existing state to calculate new state
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
        this.setState((state, props) => {
            return {
                searchQueries: [...state.searchQueries, queryString],
            }
        });
        this.setState({loading: true})
    };

    componentDidMount() {
        this.setState({
            // items: fakeItems.books,
            // searches: recentSearches,
            servers: servers.servers
        });

        //TODO: How do I pass a variable to this...
        let socket = new WebSocket("ws://127.0.0.1:8080/ws");

        socket.onopen = () => {
            console.log("Successfully Connected");
            this.setState({
                socket: socket
            });

            this.state.socket.send(JSON.stringify({
                type: MessageTypes.CONNECT,
                payload: {
                    //TODO: Make server use the connection name...
                    name: "bot"
                }
            }));

            notification.open({
                key: "timer",
                type: "info",
                message: "Please wait before searching",
                description: this.state.timeLeft + " seconds",
                duration: 0
            });

            let downloadTimer = setInterval(() => {
                this.setState((state, props) => {
                    return {timeLeft: state.timeLeft -= 1}
                });
                notification.open({
                    key: "timer",
                    type: "info",
                    message: "Please wait before searching",
                    description: this.state.timeLeft + " seconds",
                    duration: 0
                });
                if (this.state.timeLeft <= 0) {
                    clearInterval(downloadTimer);
                    notification.open({
                        key: "timer",
                        type: "success",
                        message: "Server is ready",
                        description: "Enter a search query to get started"
                    });
                }
            }, 1000);
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
            this.setState((state, props) => messageRouter(JSON.parse(message.data), state))
        }
    }

    sidebarStyle = {
        paddingTop: "16px",
        paddingLeft: "8px",
    };

    pastSearchHandler = (index) => {
        this.setState((state, props) => {
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
                                    disabled={this.state.timeLeft !== 0 || this.state.loading}
                                    searchCallback={this.searchCallback}/>
                            {/*Show instructions if the user hasn't yet search for anything*/}
                            {this.state.searchQueries.length === 0 && <Result title={"Search a book to get started"}/>}
                            {/*Show loading indicator when waiting for server results*/}
                            {this.state.loading && <Result title={"Loading your results. Please wait."}/>}
                            {/*Show table when items are received and we aren't loading*/}
                            {this.state.items.length > 0 &&
                            <BookTable
                                items={this.state.items}
                                socket={this.state.socket}
                                disabled={this.state.loading}/>}
                        </div>
                    </Content>
                </Layout>
            </Layout>
        );
    }
}

export default App;