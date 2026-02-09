import { Table } from 'antd';
import { useQuery } from '@tanstack/react-query'
import './App.css'

type ListDetails = {
    totalprice: number;
    spendingLimit: number;
}

const getListDetails = async (): Promise<ListDetails> => {
    const res = await fetch(`http://localhost:1323/list`, {
        method: 'GET',
    });
    if (!res.ok) {
        console.log('Failed to calculate total list price');
        throw new Error('to calculate total list price');
    }
    return res.json();
};

export default function TotalPriceTableSummary() {
    const { error, data } = useQuery({
        queryKey: ['listdetails'],
        queryFn: () => {
            return getListDetails();
        }
    });

    if (error) {
        console.log(error);
    }

    const formatTotal = (details: ListDetails | undefined): string => {
        return `${details?.totalprice ?? 'Calculating ...'}`;
    };

    return (
        <Table.Summary.Row>
            <Table.Summary.Cell index={0}>Total</Table.Summary.Cell>
            <Table.Summary.Cell index={1} />
            <Table.Summary.Cell index={2} colSpan={3}>{formatTotal(data)}</Table.Summary.Cell>
        </Table.Summary.Row>
    );
}

