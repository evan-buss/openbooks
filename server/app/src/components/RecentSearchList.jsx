import React from 'react';
import Menu from "antd/es/menu";
import Icon from "antd/es/icon";
import PropTypes from 'prop-types';

import Typography from 'antd/es/typography';

const {Title} = Typography;


const titleStyle = {
  color: "white",
  textAlign: "center"
};

export default class RecentSearchList extends React.Component {
  render() {
    return (
      <div>
        <Title level={3} style={titleStyle}>Recent Searches</Title>
        <Menu
          mode="inline"
          theme="dark"
        >
          {this.props.searches.length === 0 && <Menu.Item disabled={true}>No Recent Searches</Menu.Item>}
          {this.props.searches.length > 0 && this.props.searches.map((search, index) => {
            return (
              <Menu.Item key={"" + index} onClick={() => this.props.clickHandler(index)}>
                <Icon type="file-search"/>
                <span>{search.length > 17 ? search.substring(0, 17) + "..." : search}</span>
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
  clickHandler: PropTypes.func
};