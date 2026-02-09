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
    const query = useQuery({
        queryKey: ['listdetails'],
        queryFn: () => {
            return getListDetails();
        }
    });

    if (query.error) {
        console.log(query.error);
    }

    const formatTotal = (details: ListDetails | undefined): string => {
        return `${details?.totalprice ?? 'Calculating ...'}`;
    };
    const formatAlertMsg = (details: ListDetails | undefined): string => {
        if (!details) {
            return 'Calculating ...';
        } else if (details && details.totalprice <= details.spendingLimit) {
            return '(total price is within spending limit)';
        }
        return '(WARNING: total price is greater than spending limit)';
    }

    return (
        <Table.Summary.Row>
            <Table.Summary.Cell index={0}>Total</Table.Summary.Cell>
            <Table.Summary.Cell index={1} />
            <Table.Summary.Cell index={2} >{formatTotal(query.data)}</Table.Summary.Cell>
            <Table.Summary.Cell index={3} colSpan={2}>{formatAlertMsg(query.data)}</Table.Summary.Cell>
        </Table.Summary.Row>
    );
}

