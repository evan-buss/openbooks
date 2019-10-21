import React from 'react';
import './App.css';
import Countdown from './components/Countdown'
import Search from './components/Search';
import Table from './components/Table';
import {messageRouter, MessageTypes} from "./messages"
import GridLoader from 'react-spinners/GridLoader'

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
    }

    componentDidMount() {
        let socket = new WebSocket("ws://127.0.0.1:8080/ws");

        socket.onopen = () => {
            console.log("Successfully Connected");
            this.setState({
                socket: socket
            })

            this.state.socket.send(JSON.stringify({
                type: MessageTypes.CONNECT,
                payload: {
                    name: "bot"
                }
            }))

            let downloadTimer = setInterval(() => {
                this.setState((state, props) => {
                    return {timeLeft: state.timeLeft -= 1}
                })
                if (this.state.timeLeft <= 0) {
                    clearInterval(downloadTimer);
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

    render() {
        return (
            <div className="app-body">
                <div id="app-container">
                    <h1>{this.state.status}</h1>
                    <Search socket={this.state.socket}
                            loadingCallback={this.loadingCallback}/>
                    {this.state.timeLeft > 0 && <Countdown time={this.state.timeLeft}/>}
                    {this.state.loading && <GridLoader
                        color="#09d3ac"
                        css={
                            `size: 100;`
                        }
                    />}
                    {this.state.items.length > 0 && <Table items={this.state.items} socket={this.state.socket} />}
                </div>
            </div>
        );
    }
}

export default App;