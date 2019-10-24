import React from 'react';
import Button from 'antd/es/button';
import Icon from 'antd/es/icon';
import Table from 'antd/es/table';
import Title from 'antd/es/typography/Title';
import Input from 'antd/es/input';
import Highlighter from 'react-highlight-words';
import PropTypes from 'prop-types';

export default class BookTable extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      searchText: '',
      onlyRecommended: false,
      height: 0
    };

    // Create a key for each item using the index
    this.props.items.forEach((item, index) => {
      item.key = index;
    })

    this.updateWindowDimensions = this.updateWindowDimensions.bind(this);
  }

  // My personal favorite servers. Can change over time. I have just 
  // found these to be the most reliable and they are always my first choice
  recommendedServers = [
    "oatmeal",
    "pondering42"
  ];


  componentWillMount() {
    this.updateWindowDimensions();
    window.addEventListener('resize', this.updateWindowDimensions);
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.updateWindowDimensions);
  }

  // Set the table size according to certain screen resolutions
  updateWindowDimensions() {
    let rows;
    let height = window.innerHeight;
    if (height >= 1300) {
      rows = 18;
    } else if (height >= 1100) {
      rows = 16;
    } else if (height >= 1000) {
      rows = 14;
    } else if (height >= 910) {
      rows = 12;
    } else if (height >= 800) {
      rows = 10;
    } else if (height >= 720) {
      rows = 8;
    } else {
      rows = 6;
    }
    console.log(rows)
    this.setState({ height: rows });
  }

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

  // Modify the list of items based on the user's preferences
  getItems() {
    if (this.state.onlyRecommended) {
      return this.props.items.filter((item) => {
        return this.recommendedServers.includes(item.server.toLowerCase());
      })
    }
    return this.props.items
  }

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
              <Button onClick={() => this.props.downloadCallback(record.full)}>Download</Button>
            )
        }
      }
    ];

    const rows = this.state.height;

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

          pagination={{
            position: "bottom",
            showSizeChanger: true,
            defaultPageSize: rows,
            pageSizeOptions: ['6', '8', '10', '12', '14', '16', '18', '20']
          }}
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
  downloadCallback: PropTypes.func,
  disabled: PropTypes.bool
};