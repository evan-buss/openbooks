import { Server } from '@styled-icons/feather/Server';
import { Pane, Text } from 'evergreen-ui';
import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { selectServers } from "../state/serverSlice";
import { toggleServerFilter } from '../state/stateSlice';
import { RootState } from '../state/store';

const ServerList: React.FC = () => {
    const dispatch = useDispatch();
    const servers = useSelector(selectServers);
    const filters = useSelector((store: RootState) => new Set<string>(store.state.serverFilters));

    return (
        <>
            <Pane display="flex" justifyContent="center">
                <Text marginX="auto" marginY={4} color="muted">Click to filter.</Text>
            </Pane >
            {servers.map(name =>
                (
                    <Pane cursor="pointer"
                        onClick={() => dispatch(toggleServerFilter(name))}
                        key={name}
                        border
                        background={filters.has(name) ? "greenTint" : ""}
                        padding={6} elevation={1} margin={8}>
                        <Server size={24} color="green" title="server icon" />
                        <Text marginLeft={24}>{name}</Text>
                    </Pane>
                )
            )}
        </>
    );
};

export default ServerList;
