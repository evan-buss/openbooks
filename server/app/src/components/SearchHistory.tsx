import { Dispatch } from "@reduxjs/toolkit";
import { Search } from "@styled-icons/feather/Search";
import { Badge, Pane, Position, Spinner, Text, Tooltip } from 'evergreen-ui';
import React from 'react';
import { useDispatch, useSelector } from "react-redux";
import { deleteHistoryItem, HistoryItem, selectHistory } from "../state/historySlice";


const SearchHistory: React.FC = () => {
    const history = useSelector(selectHistory);
    const dispatch = useDispatch();

    return (
        <>
            {
                history.length > 0 ?
                    history.map((item: HistoryItem) => <HistoryCard key={item.timestamp.toString()} item={item} dispatch={dispatch} />)
                    : <Pane display="flex" justifyContent="center">
                        <Text marginX="auto" marginY={16} color="muted">History is a mystery.</Text>
                    </Pane >
            }
        </>
    );
}

type Props = {
    item: HistoryItem;
    dispatch: Dispatch<any>
}

const HistoryCard: React.FC<Props> = ({ item, dispatch }: Props) => {
    return (
        <Tooltip position={Position.RIGHT} content="Click to delete.">
            <Pane
                cursor="pointer"
                onClick={() => dispatch(deleteHistoryItem(item.timestamp))}
                border
                padding={6}
                elevation={1}
                margin={8}
                display="flex">
                <Search size={24} title="search history icon" />
                <Pane marginLeft={12}
                    display="flex"
                    justifyContent="space-between"
                    alignItems="center"
                    width="100%">
                    <Text width={140}
                        whiteSpace="nowrap"
                        display="block"
                        fontWeight={500}
                        overflow="hidden"
                        textOverflow="ellipsis">
                        {item.query}
                    </Text>
                    {!item.results?.length ?
                        <Spinner size={24} /> :
                        <Badge color="green" size={300}>
                            {item.results?.length + ' Results'}
                        </Badge>
                    }
                </Pane>
            </Pane>
        </Tooltip>
    )
}

export default SearchHistory;