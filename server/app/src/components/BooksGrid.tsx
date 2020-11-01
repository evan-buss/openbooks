import { Button, Pane, Spinner, Table, Text } from 'evergreen-ui';
import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { BookDetail, MessageType } from '../models/messages';
import { sendMessage } from '../state/stateSlice';
import { RootState } from '../state/store';


export const BooksGrid: React.FC = () => {
    const activeItem = useSelector((store: RootState) => store.state.activeItem);
    const servers = useSelector((store: RootState) => new Set<string>(store.state.serverFilters));
    const dispatch = useDispatch();

    const download = (book: BookDetail) => {
        dispatch(sendMessage({ type: MessageType.DOWNLOAD, payload: { book: book.full } }));
    }

    const filtered = (books: BookDetail[], servers: Set<string>) => {
        if (servers.size === 0) return books;
        return books.filter(x => servers.has(x.server));
    }
    const filteredRows = filtered(activeItem?.results ?? [], servers);

    if (activeItem === null) {
        return (<Text fontSize="2em" margin="3em">Search a book to get started.</Text>);
    }

    // Active item selected, but the results are null
    // The item's results haven't been loadeded. Show 
    // loading indicator
    if (activeItem !== null && !activeItem.results) {
        return (
            <Pane>
                <Spinner marginX="auto" marginY={120} />
            </Pane>
        );
    }

    // Results length is zero. No results for query
    if (activeItem.results?.length === 0) {
        return (<Text fontSize="2em" margin="3em">No results.</Text>);
    }

    // FilteredRows length is zero. Filter not valid.
    if (filteredRows.length === 0) {
        return (<Text fontSize="2em" margin="3em">No books matching filter.</Text>);
    }

    return (
        <Table width="100%" height="100%">
            <Table.Head>
                <Table.TextHeaderCell flexBasis={200} flexGrow={0} flexShrink={0}></Table.TextHeaderCell>
                <Table.TextHeaderCell flexBasis={250} flexGrow={0} flexShrink={0} >Author</Table.TextHeaderCell>
                <Table.TextHeaderCell>Title</Table.TextHeaderCell>
                <Table.TextHeaderCell flexBasis={100} flexGrow={0} flexShrink={0}>Format</Table.TextHeaderCell>
                <Table.TextHeaderCell flexBasis={100} flexGrow={0} flexShrink={0}>Size</Table.TextHeaderCell>
                <Table.TextHeaderCell flexBasis={150} flexGrow={0} flexShrink={0} >Download?</Table.TextHeaderCell>
            </Table.Head>
            <Table.VirtualBody height={900}>
                {
                    filteredRows.map((book, i) => (
                        <Table.Row key={book.full + i}>
                            <Table.TextCell flexBasis={200} flexGrow={0} flexShrink={0}>{book.server}</Table.TextCell>
                            <Table.TextCell flexBasis={250} flexGrow={0} flexShrink={0}>{book.author}</Table.TextCell>
                            <Table.TextCell>{book.title}</Table.TextCell>
                            <Table.TextCell flexBasis={100} flexGrow={0} flexShrink={0}>{book.format}</Table.TextCell>
                            <Table.TextCell flexBasis={100} flexGrow={0} flexShrink={0}>{book.size}</Table.TextCell>
                            <Table.Cell flexBasis={150} flexGrow={0} flexShrink={0}>
                                <Button appearance="primary"
                                    onClick={(e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => download(book)}>
                                    Download
                            </Button>
                            </Table.Cell>
                        </Table.Row>
                    ))
                }
            </Table.VirtualBody >
        </Table >
    )
}
