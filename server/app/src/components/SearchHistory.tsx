import { Search } from "@styled-icons/feather/Search";
import { Button, Pane, Text } from 'evergreen-ui';
import React from 'react';
import { useHistory } from '../models/HistoryProvider';


const SearchHistory: React.FC = () => {
    const { history, addItem } = useHistory()!;

    return (
        <>
            <Button onClick={() => addItem({ query: "hello world", time: new Date().getTime() })}>Add</Button>
            {
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
                            <Text color="muted" size={300}>{!item.results?.length ? 'In Progress' : item.results?.length + ' Results'}</Text>
                        </Pane>
                    </Pane>
                ))
            }
        </>
    );
}

export default SearchHistory;