import React from 'react';
import Input from "antd/es/input";
import Button from "antd/es/button";
import PropTypes from "prop-types";

export default class Search extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            searchTerm: ""
        }
    }

    // Update the searchTerm state when the input's value changes
    handleChange = (event) => {
        this.setState({
            searchTerm: event.target.value
        })
    };

    // Send the current query to the callback
    handleSearch = () => {
        this.props.searchCallback(this.state.searchTerm);
    };

    // Pressing enter in the search box should send the search query to the server
    _handleEnterPress = (event) => {
        if (event.key === 'Enter') {
            this.handleSearch();
        }
    }

    render() {
        return (
            <div style={searchContainer}>
                <Input style={inputStyle}
                    placeholder="Search for a book" size="large"
                    onChange={this.handleChange}
                    disabled={this.props.disabled}
                    onKeyPress={this._handleEnterPress}
                />
                <Button
                    size="large"
                    type="primary"
                    icon="search"
                    onClick={this.handleSearch}
                    disabled={this.props.disabled}>
                    Search
                </Button>
            </div>
        )
    }
}

const searchContainer = {
    display: "flex",
    flexFlow: "row nowrap",
    width: "80%",
    margin: 20,
};

const inputStyle = {
    marginRight: 50
};


Search.propTypes = {
    searchCallback: PropTypes.func,
    disabled: PropTypes.bool
};