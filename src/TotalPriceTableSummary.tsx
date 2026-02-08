import { Table } from 'antd';
import { useQuery } from '@tanstack/react-query'
import './App.css'

type PriceReport = {
    totalprice: number;
}

const calcTotalPrice = async (): Promise<PriceReport> => {
    const res = await fetch(`http://localhost:1323/list/total`, {
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
        queryKey: ['totalprice'],
        queryFn: () => {
            return calcTotalPrice();
        }
    });

    if (error) {
        console.log(error);
    }

    const formatTotal = (r: PriceReport | undefined): string => {
        return `${r?.totalprice ?? 'Calculating ...'}`;
    };

    return (
        <Table.Summary.Row>
            <Table.Summary.Cell index={0}>Total</Table.Summary.Cell>
            <Table.Summary.Cell index={1} />
            <Table.Summary.Cell index={2} colSpan={3}>{formatTotal(data)}</Table.Summary.Cell>
        </Table.Summary.Row>
    );
}

