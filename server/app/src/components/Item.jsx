import React from 'react';
import PropTypes from 'prop-types'
import './Item.css';


class Item extends React.Component {

    getBook = () => {
        this.props.callback(this.props.data.full)
    }

    render() {
        return (
            <tr id="item">
                <td> {this.props.data.server} </td>
                <td> {this.props.data.title} </td>
                <td> {this.props.data.author} </td>
                <td> {this.props.data.format} </td>
                <td> {this.props.data.size} </td>
                <td>
                    <button onClick={this.getBook}>Download</button>
                </td>
            </tr>
        )
    }
}

Item.propTypes = {
    data: PropTypes.object.isRequired
}

export default Item