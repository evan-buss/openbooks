import { Badge, DatabaseIcon, Pane, Text } from 'evergreen-ui';
import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useGetServersQuery } from "../state/api";
import { toggleServerFilter } from '../state/stateSlice';
import { RootState } from '../state/store';

const ServerList: React.FC = () => {
    const dispatch = useDispatch();
    const { data: servers, error, isLoading } = useGetServersQuery(null);
    const filters = useSelector((store: RootState) => new Set<string>(store.state.serverFilters));

    const stateText = (): string => {
        if (isLoading) {
            return 'Loading Server List';
        } else if (error) {
            return 'Error Loading Server List'
        } else if (servers?.elevatedUsers?.length === 0) {
            return 'Servers not found.';
        } else {
            return 'Click to filter.';
        }
    }

    return (
        <>
            <Pane display="flex" justifyContent="center">
                <Text marginX="auto" marginY={4} color="muted">
                    {stateText()}
                </Text>
            </Pane >
            {servers?.elevatedUsers?.map(name =>
            (
                <Pane
                    borderLeft={filters.has(name) ? '3px solid #1070CA' : '1px solid #E4E7EB;'}
                    cursor="pointer"
                    padding={6}
                    elevation={1}
                    margin={8}
                    display="flex"
                    justifyContent="space-between"
                    alignItems="center"
                    border
                    onClick={() => dispatch(toggleServerFilter(name))}
                    key={name}>
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
