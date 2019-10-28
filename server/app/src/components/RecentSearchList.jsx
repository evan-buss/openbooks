import React from 'react';
import Menu from "antd/es/menu";
import Icon from "antd/es/icon";
import PropTypes from 'prop-types';
import Title from 'antd/es/typography/Title';

import "../App.css"

const titleStyle = {
  color: "white",
  paddingTop: 16,
  paddingLeft: 8,
};

// Use material design menu item style
const menuItemStyle = {
  borderRadius: "0 25px 25px 0",
  display: "flex",
  flexFlow: "row nowrap",
  justifyContent: "space-between",
  alignItems: "center"
}

export default class RecentSearchList extends React.Component {
  render() {
    return (
      <div>
        <Title level={3} style={titleStyle}>Recent Searches</Title>
        <Menu
          selectedKeys={"" + this.props.selected}
          mode="inline"
          theme="dark"
          style={{ width: "220px" }}
        >
          {this.props.searches.length === 0 &&
            <Menu.Item disabled={true}>No Recent Searches</Menu.Item>}
          {this.props.searches.length > 0 && this.props.searches.map((search, index) => {
            return (
              <Menu.Item
                disabled={this.props.disabled}
                style={menuItemStyle}
                key={"" + index}
                onClick={() => this.props.clickHandler(index)}>
                <div>
                  <Icon type="file-search" />
                  <span>{search}</span>
                </div>
                <Icon type="close" className="close-button" onClick={() => this.props.deleteHandler(search, index)}/>
              </Menu.Item>
            );
          })}
        </Menu>
      </div>
    )
  }
}

RecentSearchList.propTypes = {
  searches: PropTypes.arrayOf(PropTypes.string),
  clickHandler: PropTypes.func,
  selected: PropTypes.number
};