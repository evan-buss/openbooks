import React from 'react';
import Typography from 'antd/es/typography'
import Icon from 'antd/es/icon';
import List from 'antd/es/list';
import Spin from 'antd/es/spin';
import PropTypes from 'prop-types'

const { Title } = Typography;

const listStyle = {
  overflow: 'auto',
  height: '100vh',
  position: 'fixed',
  width: "220px",
  paddingTop: "16px",
  paddingRight: "8px",
  backgroundColor: "#001529",
  right: 0,
};

const spinStyle = {
  display: "flex",
  flexFlow: "row nowrap",
  justifyContent: "center",
  alignItems: "center",
  marginTop: "80px"
};

export default class ServerList extends React.Component {
  render() {
    const { servers } = this.props;
    return (
      <div style={listStyle}>
        <Title level={3} style={{ color: "white", textAlign: "center" }}>Online Servers</Title>

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

ServerList.propTypes = {
  servers: PropTypes.arrayOf(PropTypes.string)
};