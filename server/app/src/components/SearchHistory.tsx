import { Search } from "@styled-icons/feather/Search";
import { Badge, Pane, Spinner, Text } from 'evergreen-ui';
import React from 'react';
import { useHistory } from '../models/HistoryProvider';


const SearchHistory: React.FC = () => {
    const { history } = useHistory()!;


    const historyList = () => {
        history.map(item => (
            <Pane border padding={6} elevation={1} margin={8} key={item.time} display="flex">
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
        ))
    }


    return (
        <>
            {
                history.length > 0 ?
                    historyList() :
                    <Pane display="flex" justifyContent="center">
                        <Text marginX="auto" color="muted">History is a mystery.</Text>
                    </Pane>
            }
        </>
    );
}

export default SearchHistory;