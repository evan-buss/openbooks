import React from 'react';
import './Search.css'

const Search = (props) => {
  return (
    <div id="search-container">
      <input id="search" type="text" />
      <input id="search-button" type="button" value="Search" />
    </div>
  )
}

export default Search;