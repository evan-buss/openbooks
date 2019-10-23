import React from 'react';
import Title from 'antd/es/typography/Title'
import Icon from 'antd/es/icon';
import List from 'antd/es/list';
import Spin from 'antd/es/spin';
import Tooltip from 'antd/es/tooltip';
import PropTypes from 'prop-types'

export default class ServerList extends React.Component {
  render() {
    const { servers } = this.props;
    return (
      <div>
        <Tooltip placement="right" title="Click to refresh">
          <Title
            level={3}
            style={titleStyle}
            onClick={this.props.reloadCallback}>
            Online Servers
          </Title>
        </Tooltip>


        {servers.length !== 0 ? (<List
          size="small"
          itemLayout="horizontal"
          split={false}
          dataSource={servers}
          renderItem={item => (
            <List.Item>
              <Icon type="cloud-server"
                style={{
                  paddingLeft: "20px",
                  color: "rgb(132, 240, 79)",
                  fontSize: "20px",
                  marginRight: "10px"
                }} />
              <span style={{ color: "white", fontWeight: "bold" }}>{item}</span>
            </List.Item>
          )} />
        ) : (<Spin size="large" style={spinStyle} />)}
      </div>

    )
  }
}

const spinStyle = {
  display: "flex",
  flexFlow: "row nowrap",
  justifyContent: "center",
  alignItems: "center",
  marginTop: 80
};

const titleStyle = {
  color: "white",
  paddingLeft: 8,
  cursor: "pointer"
}

ServerList.propTypes = {
  servers: PropTypes.arrayOf(PropTypes.string)
};