import { Server } from '@styled-icons/feather/Server';
import { Pane, Text } from 'evergreen-ui';
import React from 'react';

type Props = {
    servers: string[]
}

const ServerList: React.FC<Props> = ({ servers }: Props) => {
    console.log(servers);

    return (
        <>
            {servers.map(name =>
                (<Pane border padding={6} elevation={1} margin={8} key={name}>
                    <Server size={24} color="green" title="server icon" />
                    <Text marginLeft={12}>{name}</Text>
                </Pane>)
            )}
        </>
    );
};


export default ServerList;