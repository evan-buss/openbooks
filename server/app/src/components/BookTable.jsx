import React from 'react';

import {Button, Icon, Input, Table, Typography} from 'antd';
import Highlighter from 'react-highlight-words';
import {MessageTypes} from "../messages";
import PropTypes from 'prop-types';

const {Title} = Typography;

class BookTable extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      searchText: '',
    };

    // Create a key for each item using the index
    this.props.items.forEach((item, index) => {
      item.key = index;
    })
  }


  getColumnSearchProps = dataIndex => ({
    filterDropdown: ({setSelectedKeys, selectedKeys, confirm, clearFilters}) => (
      <div style={{padding: 8}}>
        <Input
          ref={node => {
            this.searchInput = node;
          }}
          placeholder={`Search ${dataIndex}`}
          value={selectedKeys[0]}
          onChange={e => setSelectedKeys(e.target.value ? [e.target.value] : [])}
          onPressEnter={() => this.handleSearch(selectedKeys, confirm)}
          style={{width: 188, marginBottom: 8, display: 'block'}}
        />
        <Button
          type="primary"
          onClick={() => this.handleSearch(selectedKeys, confirm)}
          icon="search"
          size="small"
          style={{width: 90, marginRight: 8}}
        >
          Search
        </Button>
        <Button onClick={() => this.handleReset(clearFilters)} size="small" style={{width: 90}}>
          Reset
        </Button>
      </div>
    ),
    filterIcon: filtered => (
      <Icon type="search" style={{color: filtered ? '#1890ff' : undefined}}/>
    ),
    onFilter: (value, record) =>
      record[dataIndex]
        .toString()
        .toLowerCase()
        .includes(value.toLowerCase()),
    onFilterDropdownVisibleChange: visible => {
      if (visible) {
        setTimeout(() => this.searchInput.select());
      }
    },
    render: text => (
      <Highlighter
        highlightStyle={{backgroundColor: '#ffc069', padding: 0}}
        searchWords={[this.state.searchText]}
        autoEscape
        textToHighlight={text.toString()}
      />
    ),
  });


  handleSearch = (selectedKeys, confirm) => {
    confirm();
    this.setState({searchText: selectedKeys[0]});
  };

  handleReset = clearFilters => {
    clearFilters();
    this.setState({searchText: ''});
  };

  handleDownload = (bookString) => {
    console.log("DOWNLOAD CALLBACK")
    console.log(bookString);
    if (this.props.socket != null) {
      this.props.socket.send(JSON.stringify({
        type: MessageTypes.DOWNLOAD,
        payload: {
          book: bookString
        }
      }))
    }
  };

  tableStyle = {
    width: "95%",
    margin: "30px"
  };

  render() {
    const columns = [
      {
        title: 'Server',
        dataIndex: 'server',
        key: 'server',
        ...this.getColumnSearchProps('server'),
      },
      {
        title: 'Title',
        dataIndex: 'title',
        key: 'title',
        ...this.getColumnSearchProps('title'),
      },
      {
        title: 'Author',
        dataIndex: 'author',
        key: 'author',
        ...this.getColumnSearchProps('author'),
      },
      {
        // Note: May enable filters in the future. Right now it is too much of a hassle.
        title: 'Format',
        dataIndex: 'format',
        key: 'format',
        ...this.getColumnSearchProps('format'),
      },
      {
        // NOTE: May enable sorting in the future. Must first convert the string to a number
        title: 'Size',
        dataIndex: 'size',
        key: 'size',
      },
      {
        title: 'Download',
        dataIndex: 'download',
        key: 'download',
        width: "5%",
        render: (text, record) => {
          return this.props.disabled ? (
            <Button disabled>Download</Button>
          ) : (
            <Button onClick={() => this.handleDownload(record.full)}>Download</Button>
          )
        }

      }
    ];
    return (
      <Table title={() => <Title>Search Results</Title>}
             style={this.tableStyle}
             columns={columns}
             dataSource={this.props.items}
             loading={this.props.disabled}
             pagination={{position: "bottom", defaultPageSize: 10}}
      />
    )
  }
}

Table.propTypes = {
  items: PropTypes.arrayOf(PropTypes.object),
  socket: PropTypes.object,
  disabled: PropTypes.bool
};

export default BookTable