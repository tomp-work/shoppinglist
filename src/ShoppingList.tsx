import { Card, Divider, Table } from 'antd';
import { useQuery } from '@tanstack/react-query'
import AddForm from './AddForm.tsx'
import './App.css'

const getItemList = async () => {
    const res = await fetch("http://localhost:1323/item");
    if (!res.ok) {
        console.log('Failed to fetch list items');
        throw new Error('Failed to fetch list items');
    }
    return res.json();
};

export default function ShoppingList() {
    const { isPending, error, data } = useQuery({
        queryKey: ['items'],
        queryFn: () => {
            return getItemList();
        }
    });

    // TODO: Improve FE status/error handling.
    if (isPending) {
        return 'Loading...';
    }
    if (error) {
        return 'An error has occurred: ' + error.message;
    }

    // AntD table column definition for shopping list table.
    const columns = [
        {
            title: 'Name',
            dataIndex: 'name',
            key: 'name',
        },
        {
            title: 'Quantity',
            dataIndex: 'quantity',
            key: 'quantity',
        },
    ];

    return (
        <Card title="Shopping List">
            <AddForm />
            <Divider />
            <Table bordered rowKey="id" dataSource={data} columns={columns} />
        </Card>
    )
}
