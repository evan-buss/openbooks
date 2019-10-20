import React from 'react';
import './Search.css'
import { MessageTypes } from '../messages';


class Search extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      searchTerm: ""
    }
  }

  handleChange = (event) => {
    this.setState({
      searchTerm: event.target.value
    })
  }

  handleSearch = () => {
    this.props.loadingCallback(true)
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
        <input id="search" type="text" onChange={this.handleChange}
          placeholder="Book Search" />
        <input id="search-button" type="button" value="Search"
          onClick={this.handleSearch} />
      </div>
    )
  }

}

export default Search;