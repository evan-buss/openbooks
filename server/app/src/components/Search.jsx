import React from 'react';
import './Search.css'
import {MessageTypes} from '../messages';
import Input from "antd/es/input";
import Button from "antd/es/button";


class Search extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            searchTerm: ""
        }
    }

    handleChange = (event) => {
        console.log(event.target.value);
        this.setState({
            searchTerm: event.target.value
        })
    }

    handleSearch = () => {
        this.props.loadingCallback(true);
        this.props.socket != null && this.props.socket.send(JSON.stringify({
            type: MessageTypes.SEARCH,
            payload: {
                query: this.state.searchTerm
            }
        }))
    }

    render() {
        return (
            <div id="search-container">
                <Input placeholder="Search for a book" size="large" onChange={this.handleChange}
                       disabled={this.props.timeLeft > 0}/>
                <Button size="large" type="primary" icon="search" onClick={this.handleSearch}
                        disabled={this.props.timeLeft > 0}>Search</Button>
            </div>
        )
    }
}

export default Search;