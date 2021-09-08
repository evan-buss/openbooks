import { Button, majorScale, Pane, SearchInput } from 'evergreen-ui';
import React, { FormEvent, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import styled from 'styled-components';
import { BooksGrid } from '../components/BooksGrid';
import { sendSearch } from '../state/stateSlice';
import { RootState } from '../state/store';

const SearchPage: React.FC = () => {
    const dispatch = useDispatch();
    const [searchQuery, setSearchQuery] = useState("");
    const activeItem = useSelector((store: RootState) => store.state.activeItem);

    const searchHandler = (event: FormEvent) => {
        event.preventDefault();
        if (searchQuery === "") return;
        dispatch(sendSearch(searchQuery));
        setSearchQuery("");
    }

    return (
        <Pane margin={majorScale(4)} width="100%" display="flex" flexDirection="column" alignItems="center">
            <Form onSubmit={(e) => searchHandler(e)}>
                <SearchInput
                    disabled={activeItem !== null && !activeItem.results}
                    value={searchQuery}
                    onChange={(e: any) => setSearchQuery(e.target.value)}
                    placeholder="Search for a book."
                    height={majorScale(5)}
                    marginRight={majorScale(4)}
                    width="100%">
                </SearchInput>
                <Button type="submit" height={40} appearance="primary" disabled={searchQuery === ""}>Search</Button>
            </Form>
            <BooksGrid />
        </Pane>
    )
}

const Form = styled.form`
    width: 100%;
    display: flex;
    flex-flow: row nowrap;
    justify-content: center;
    align-items: center;
    margin-bottom: 32px;
`

export default SearchPage;
