import React from 'react';

import {Button, Icon, Input, Table, Typography} from 'antd';
import Highlighter from 'react-highlight-words';
import {MessageTypes} from "../messages";

const {Title} = Typography;

class BookTable extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            searchText: '',
        };


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
                    return <Button onClick={() => this.handleDownload(record.full)}>Download</Button>
                }
            }
        ];
        return (
            <div style={this.tableStyle}>
                <Title>Search Results</Title>
                <Table columns={columns} dataSource={this.props.items}
                       pagination={{position: "top", defaultPageSize: 10}}/>
            </div>
        )
    }
}


// import PropTypes from 'prop-types'
// import Item from './Item';
// import './BookTable.css';
// import { MessageTypes } from '../messages';

// class Table extends React.Component {
//   downloadCallback = (downloadString) => {
//     console.log("DOWNLOAD CALLBACK")
//     console.log(downloadString);
//     if (this.props.socket != null) {
//       this.props.socket.send(JSON.stringify({
//         type: MessageTypes.DOWNLOAD,
//         payload: {
//           book: downloadString
//         }
//       }))
//     }
//   }
//
//   render() {
//     return (
//       <div id="container">
//         <h1>Search Results</h1>
//         {this.props.items.length > 0 ? (
//           <table>
//             <thead>
//               <tr>
//                 <th>Server</th>
//                 <th>Title</th>
//                 <th>Author</th>
//                 <th>Format</th>
//                 <th>Size</th>
//                 <th>Download</th>
//               </tr>
//             </thead>
//             <tbody>
//               {this.props.items.map((item) =>
//                 <Item key={item.title + item.size + item.server}
//                   data={item}
//                   callback={this.downloadCallback} />
//               )}
//             </tbody>
//           </table>
//         ) : (
//             <h2>Search a book to view the results.</h2>
//           )
//         }
//       </div>
//     )
//   }
// }
//
// Table.propTypes = {
//   items: PropTypes.arrayOf(PropTypes.object).isRequired,
//   socket: PropTypes.object
// }

export default BookTable