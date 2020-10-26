import { Button, Table } from 'evergreen-ui';
import React from 'react';

type Props = {

}

const SearchResults: React.FC<Props> = (props) => {
    const download = (query: string, event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
        event.stopPropagation();
        console.log(query);
    }

    return (
        <Table width="100%" margin={28}>
            <Table.Head>
                <Table.SearchHeaderCell placeholder="Server"></Table.SearchHeaderCell>
                <Table.TextHeaderCell>Author</Table.TextHeaderCell>
                <Table.TextHeaderCell>Title</Table.TextHeaderCell>
                <Table.TextHeaderCell>Format</Table.TextHeaderCell>
                <Table.TextHeaderCell>Size</Table.TextHeaderCell>
                <Table.TextHeaderCell>Full</Table.TextHeaderCell>
                <Table.TextHeaderCell>Download?</Table.TextHeaderCell>
            </Table.Head>
            <Table.VirtualBody height={800}>
                {['evan', 'adam', 'ian'].map(name => (
                    <Table.Row key={name} isSelectable onSelect={() => alert(name)}>
                        <Table.TextCell>{name}</Table.TextCell>
                        <Table.TextCell>{name}</Table.TextCell>
                        <Table.TextCell>{name}</Table.TextCell>
                        <Table.TextCell>{name}</Table.TextCell>
                        <Table.TextCell>{name}</Table.TextCell>
                        <Table.TextCell>{name}</Table.TextCell>

                        <Table.Cell>
                            <Button appearance="primary"
                                onClick={(e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => download(name, e)}>
                                Download
                            </Button>
                        </Table.Cell>
                    </Table.Row>
                ))}
            </Table.VirtualBody>
        </Table>
    )
}


export default SearchResults;