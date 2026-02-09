import { Checkbox, type CheckboxProps } from 'antd';
import { useMutation, useQueryClient } from '@tanstack/react-query'
import './App.css';

type ItemUpdate = {
    id: string;
    picked: boolean;
}

const updateItem = async (update: ItemUpdate) => {
    const response = await fetch(`http://localhost:1323/item/${update.id}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({picked: update.picked}),
    });
    if (!response.ok) {
        throw new Error('Failed to create item');
    }
    return response.json();
};

type PickedCheckboxProps = {
    id: string;
    picked: boolean;
};

export default function PickedCheckbox({ id, picked }: PickedCheckboxProps) {
    const queryClient = useQueryClient();

    const updateMutation = useMutation({
        mutationFn: (update: ItemUpdate) => updateItem(update),
        onSuccess: (id) => {
            console.log('Item updated:', id);
            queryClient.invalidateQueries({ queryKey: ["items"] });
            queryClient.invalidateQueries({ queryKey: ["listdetails"] });
        },
        onError: (error) => {
            console.error('Error deleting item:', error);
        },
    });

    const onChange: CheckboxProps['onChange'] = (e) => {
        console.log(`checked = ${e.target.checked} id=${id}`);
        updateMutation.mutate({ id: id ?? '', picked: e.target.checked });
    };

    return (
        <Checkbox checked={picked} onChange={onChange} />
    )
}
