import { Badge, DatabaseIcon, Pane, Text } from 'evergreen-ui';
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
                        display="flex"
                        alignItems="center"
                        padding={6} elevation={1} margin={6}>
                        <DatabaseIcon size={15} color="#234361" />
                        <Text width="100%" marginLeft={24} display="flex" alignItems="center" justifyContent="space-between">
                            {name}
                            {filters.has(name) &&
                                <Badge color="blue">
                                    Active
                            </Badge>
                            }
                        </Text>
                    </Pane>
                )
            )}
        </>
    );
};

export default ServerList;
