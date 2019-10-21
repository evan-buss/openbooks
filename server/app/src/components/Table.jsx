import React from 'react';
import PropTypes from 'prop-types'
import Item from './Item';
import './Table.css';
import { MessageTypes } from '../messages';

class Table extends React.Component {
  downloadCallback = (downloadString) => {
    console.log("DOWNLOAD CALLBACK")
    console.log(downloadString);
    if (this.props.socket != null) {
      this.props.socket.send(JSON.stringify({
        type: MessageTypes.DOWNLOAD,
        payload: {
          book: downloadString
        }
      }))
    }
  }

  render() {
    return (
      <div id="container">
        <h1>Search Results</h1>
        {this.props.items.length > 0 ? (
          <table>
            <thead>
              <tr>
                <th>Server</th>
                <th>Title</th>
                <th>Author</th>
                <th>Format</th>
                <th>Size</th>
                <th>Download</th>
              </tr>
            </thead>
            <tbody>
              {this.props.items.map((item) =>
                <Item key={item.title + item.size + item.server}
                  data={item}
                  callback={this.downloadCallback} />
              )}
            </tbody>
          </table>
        ) : (
            <h2>Search a book to view the results.</h2>
          )
        }
      </div>
    )
  }
}

Table.propTypes = {
  items: PropTypes.arrayOf(PropTypes.object).isRequired,
  socket: PropTypes.object
}

export default Table