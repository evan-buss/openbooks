import React from 'react';
import Item from './Item';
import './Table.css';

class Table extends React.Component {
  render() {
    return (
      <div id="container">
        <h1>Search Results</h1>
        <div id="table">
          <Item />
          <Item />
          <Item />
        </div>
      </div>
    )
  }
}

export default Table