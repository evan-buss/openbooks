import React from 'react';
import './App.css';
import Search from './components/Search';
import BookTable from './components/BookTable';
import {messageRouter, MessageTypes} from "./messages"
import {Divider, Layout, notification, Result, Typography} from 'antd';
// import fakeItems from "./dummyData";

const {Header, Sider, Content} = Layout;
const {Title} = Typography;

class App extends React.Component {
    constructor(props) {
        super(props);

        // NOTE: Set state via function if you need the existing state to calculate new state
        this.state = {
            status: "",
            searchString: "",
            items: [],
            servers: [],
            socket: null,
            loading: false,
            timeLeft: 30
        }
    }

    loadingCallback = (bool) => {
        this.setState({loading: bool})
    };

    componentDidMount() {
        // this.setState({
        //     items: fakeItems.books
        // });

        let socket = new WebSocket("ws://127.0.0.1:8080/ws");

        socket.onopen = () => {
            console.log("Successfully Connected");
            this.setState({
                socket: socket
            });

            this.state.socket.send(JSON.stringify({
                type: MessageTypes.CONNECT,
                payload: {
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
            this.setState(messageRouter(JSON.parse(message.data)))
        }
    }

    style = {textAlign: "center", color: "white", marginBottom: 0, marginTop: "5px"};
    sidebarStyle = {marginTop: "15px"};

    render() {
        return (
            <Layout className="full-size">
                <Sider>
                    <Title level={3} style={{...this.style, ...this.sidebarStyle}}>Recent Searches</Title>
                    <Divider/>
                </Sider>
                <Layout>
                    <Header><Title style={this.style}>OpenBooks</Title></Header>
                    <Content>
                        <div className="app-body">
                            <Search socket={this.state.socket} timeLeft={this.state.timeLeft}
                                    loadingCallback={this.loadingCallback}/>
                            {this.state.loading && <Result title={"Loading your results. Please wait."}/>}
                            {this.state.items.length > 0 && !this.state.loading ?
                                (<BookTable items={this.state.items} socket={this.state.socket}/>) : (
                                    <Result title={"Search a book to get started"}/>)}
                        </div>
                    </Content>
                </Layout>
            </Layout>
        );
    }
}

export default App;