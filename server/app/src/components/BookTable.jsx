import React from 'react';
import Button from 'antd/es/button';
import Icon from 'antd/es/icon';
import Table from 'antd/es/table';
import Title from 'antd/es/typography/Title';
import Input from 'antd/es/input';
import Highlighter from 'react-highlight-words';
import PropTypes from 'prop-types';

import { MessageTypes } from "../messages";

export default class BookTable extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      searchText: '',
      onlyRecommended: false,
    };

    // Create a key for each item using the index
    this.props.items.forEach((item, index) => {
      item.key = index;
    })
  }

  recommendedServers = [
    "oatmeal",
    "pondering42"
  ];


  getColumnSearchProps = dataIndex => ({
    filterDropdown: ({ setSelectedKeys, selectedKeys, confirm, clearFilters }) => (
      <div style={{ padding: 8 }}>
        <Input
          ref={node => {
            this.searchInput = node;
          }}
          placeholder={`Search ${dataIndex}`}
          value={selectedKeys[0]}
          onChange={e => setSelectedKeys(e.target.value ? [e.target.value] : [])}
          onPressEnter={() => this.handleSearch(selectedKeys, confirm)}
          style={{ width: 188, marginBottom: 8, display: 'block' }}
        />
        <Button
          type="primary"
          onClick={() => this.handleSearch(selectedKeys, confirm)}
          icon="search"
          size="small"
          style={{ width: 90, marginRight: 8 }}
        >
          Search
        </Button>
        <Button onClick={() => this.handleReset(clearFilters)} size="small" style={{ width: 90 }}>
          Reset
        </Button>
      </div>
    ),
    filterIcon: filtered => (
      <Icon type="search" style={{ color: filtered ? '#1890ff' : undefined }} />
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
        highlightStyle={{ backgroundColor: '#ffc069', padding: 0 }}
        searchWords={[this.state.searchText]}
        autoEscape
        textToHighlight={text.toString()}
      />
    ),
  });


  handleSearch = (selectedKeys, confirm) => {
    confirm();
    this.setState({ searchText: selectedKeys[0] });
  };

  handleReset = clearFilters => {
    clearFilters();
    this.setState({ searchText: '' });
  };

  toggleRecommended = () => {
    this.setState((state) => {
      return { onlyRecommended: !state.onlyRecommended }
    });
  }

  getItems() {
    if (this.state.onlyRecommended) {
      return this.props.items.filter((item) => {
        return this.recommendedServers.includes(item.server.toLowerCase());
      })
    }
    return this.props.items
  }

  handleDownload = (bookString) => {
    this.props.socket && this.props.socket.send(JSON.stringify({
      type: MessageTypes.DOWNLOAD,
      payload: {
        book: bookString
      }
    }))
  };

  render() {
    const columns = [
      {
        title: 'Server',
        dataIndex: 'server',
        key: 'server',
        width: "10%",
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
        // NOTE: May enable filters in the future. Right now it is too much of a hassle.
        title: 'Format',
        dataIndex: 'format',
        key: 'format',
        width: "10%",
        ...this.getColumnSearchProps('format'),
      },
      {
        // NOTE: May enable sorting in the future. Must first convert the string to a number
        title: 'Size',
        dataIndex: 'size',
        key: 'size',
        width: '10%'
      },
      {
        title: 'Download',
        dataIndex: 'download',
        key: 'download',
        width: "5%",
        render: (_text, record) => {
          return this.props.disabled ? (
            <Button disabled>Download</Button>
          ) : (
              <Button onClick={() => this.handleDownload(record.full)}>Download</Button>
            )
        }
      }
    ];
    return (
      <div style={containerStyle}>
        <Title style={titleStyle}>Search Results</Title>
        <div style={{ marginBottom: 16, alignSelf: "start" }}>
          <Button
            onClick={this.toggleRecommended}
            type={this.state.onlyRecommended ? 'primary' : ''}>
            Recommended Servers
            </Button>
        </div>
        <Table style={tableStyle}
          size="small"
          columns={columns}
          dataSource={this.getItems()}
          loading={this.props.disabled}
          pagination={{ position: "bottom", defaultPageSize: 12 }}
        />
      </div>
    )
  }
}

const titleStyle = { marginBottom: "0", alignSelf: "start" }

const containerStyle = {
  width: "95%",
  display: "flex",
  flexFlow: "column nowrap",
  justifyContent: "start",
  alignItems: "center"
}

const tableStyle = {
  width: "100%",
};

Table.propTypes = {
  items: PropTypes.arrayOf(PropTypes.object),
  socket: PropTypes.object,
  disabled: PropTypes.bool
};