import { Server } from '@styled-icons/feather/Server';
import { Pane, Text } from 'evergreen-ui';
import React from 'react';
import { useSelector } from 'react-redux';
import { selectServers } from "../state/serverSlice";


const ServerList: React.FC = () => {
    const servers = useSelector(selectServers);

    return (
        <>
            <Pane display="flex" justifyContent="center">
                <Text marginX="auto" marginY={4} color="muted">Click to filter.</Text>
            </Pane >
            {servers.map(name =>
                (
                    <Pane cursor="pointer" key={name} border padding={6} elevation={1} margin={8}>
                        <Server size={24} color="green" title="server icon" />
                        <Text marginLeft={24}>{name}</Text>
                    </Pane>
                )
            )}
        </>
    );
};

export default ServerList;