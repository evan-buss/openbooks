import React from 'react';
import {MessageTypes} from '../messages';
import Input from "antd/es/input";
import Button from "antd/es/button";
import PropTypes from "prop-types";

const searchContainer = {
    display: "flex",
    flexFlow: "row nowrap",
    alignItems: "center",
    justifyContent: "space-between",
    width: "50%",
    margin: "20px"
};

const inputStyle = {
    marginRight: "50px"
};


class Search extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            searchTerm: ""
        }
    }

    handleChange = (event) => {
        this.setState({
            searchTerm: event.target.value
        })
    };

    handleSearch = () => {
        this.props.searchCallback(true, this.state.searchTerm);
        this.props.socket != null && this.props.socket.send(JSON.stringify({
            type: MessageTypes.SEARCH,
            payload: {
                query: this.state.searchTerm
            }
        }))
    };


    render() {
        return (
            <div id="search-container" style={searchContainer}>
                <Input style={inputStyle} placeholder="Search for a book" size="large" onChange={this.handleChange}
                       disabled={this.props.disabled}/>
                <Button size="large" type="primary" icon="search" onClick={this.handleSearch}
                        disabled={this.props.disabled}>Search</Button>
            </div>
        )
    }
}

Search.propTypes = {
    searchCallback: PropTypes.func,
    socket: PropTypes.object,
    disabled: PropTypes.bool
};

export default Search;