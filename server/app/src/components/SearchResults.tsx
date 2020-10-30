// import { Button, Table } from 'evergreen-ui';
import { Table } from 'antd';
import { ColumnType } from 'antd/lib/table';
import React from 'react';
import { useDispatch } from 'react-redux';
import { sendMessage } from '../state/stateSlice';

type Props = {

}

const SearchResults: React.FC<Props> = (props) => {
    // const download = (query: string, event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    //     event.stopPropagation();
    //     console.log(query);
    // }


    const columns: ColumnType<any>[] = [
        {
            title: 'Name (all screens)',
            dataIndex: 'name',
            key: 'name',
            render: (text: string) => <a>{text}</a>,
        },
        {
            title: 'Age (medium screen or bigger)',
            dataIndex: 'age',
            key: 'age',
            responsive: ['md'],
        },
        {
            title: 'Address (large screen or bigger)',
            dataIndex: 'address',
            key: 'address',
            responsive: ['lg'],
        },
    ];

    const data = [
        {
            key: '1',
            name: 'John Brown',
            age: 32,
            address: 'New York No. 1 Lake Park',
        },
    ];

    return <Table columns={columns} dataSource={data}></Table>


    // const dispatch = useDispatch()
    // return <Button onClick={() => dispatch(sendMessage("HELLO WORLD"))}>Hello World</Button>


    // return (
    //     <Table width="100%">
    //         <Table.Head>
    //             <Table.SearchHeaderCell placeholder="Server"></Table.SearchHeaderCell>
    //             <Table.TextHeaderCell>Author</Table.TextHeaderCell>
    //             <Table.TextHeaderCell>Title</Table.TextHeaderCell>
    //             <Table.TextHeaderCell>Format</Table.TextHeaderCell>
    //             <Table.TextHeaderCell>Size</Table.TextHeaderCell>
    //             <Table.TextHeaderCell>Full</Table.TextHeaderCell>
    //             <Table.TextHeaderCell>Download?</Table.TextHeaderCell>
    //         </Table.Head>
    //         <Table.VirtualBody height={800}>
    {/* {['evan', 'adam', 'ian'].map(name => (
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
                ))} */}
    {/* </Table.VirtualBody>
        </Table> */}
    // )
}


export default SearchResults;
