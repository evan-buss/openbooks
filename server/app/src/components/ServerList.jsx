import React from 'react';
import Typography from 'antd/es/typography'
import Icon from 'antd/es/icon';
import Menu from "antd/es/menu";
import PropTypes from 'prop-types'

const {Title} = Typography;

export default class ServerList extends React.Component {
  render() {
    return (
      <div>
        <Title level={3} style={{color: "white"}}>Online Servers</Title>
        <Menu
          mode="inline"
          theme="dark"
        >
          {this.props.servers.length === 0 && <Menu.Item disabled={true}>Loading Available</Menu.Item>}
          {this.props.servers.length > 0 && this.props.servers.map((server, index) => {
            return (
              <Menu.Item key={"" + index} disabled={true}>
                <Icon type="check-circle" theme="twoTone" twoToneColor="#52c41a"
                      style={{fontSize: "20px"}}/>
                <span style={{color: "white", fontWeight: "bold"}}>{server}</span>
              </Menu.Item>
            );
          })}
        </Menu>
      </div>

    )
  }
}

ServerList.propTypes = {
  servers: PropTypes.arrayOf(PropTypes.string)
};