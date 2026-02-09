import { Card, Divider, Table } from 'antd';
import { useQuery } from '@tanstack/react-query';
import AddForm from './AddForm.tsx';
import DeleteAction from './DeleteAction.tsx';
import PickedCheckbox from './PickedCheckbox.tsx';
import MoveUpAction from './MoveUpAction.tsx';
import MoveDownAction from './MoveDownAction.tsx';
import SpendingLimitForm from './SpendingLimitForm.tsx';
import TotalPriceTableSummary from './TotalPriceTableSummary.tsx';
import './App.css';

type Item = {
    id?: string;
    name: string;
    quantity: number;
    picked: boolean;
    price: number;
}

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
        {
            title: 'Price',
            dataIndex: 'price',
            key: 'price',
        },
        {
            title: 'Picked',
            key: 'picked',
            render: (_: any, item: Item) => (<PickedCheckbox id={item.id ?? ''} picked={item.picked} />),
        },
        {
            title: "Actions",
            key: "actions",
            render: (_: any, item: Item) => (
                <>
                    <DeleteAction id={item.id ?? ''} />
                    <MoveUpAction id={item.id ?? ''} />
                    <MoveDownAction id={item.id ?? ''} />
                </>
            ),
        },
    ];

    return (
        <Card title="Shopping List">
            <AddForm />
            <Divider />
            <Table
                bordered
                rowKey="id"
                dataSource={data}
                columns={columns}
                summary={() => <TotalPriceTableSummary />}
            />
            <Divider />
            <SpendingLimitForm />
        </Card >
    )
}
