import { Button } from 'antd';
import { useMutation, useQueryClient } from '@tanstack/react-query'
import './App.css'
import { DeleteOutlined } from '@ant-design/icons';

const deleteItem = async (id: string) => {
    const response = await fetch(`http://localhost:1323/item/${id}`, {
        method: 'DELETE',
    });
    if (!response.ok) {
        console.log('Failed to delete item');
        throw new Error('Failed to delete item');
    }
};

type DeleteActionProps = {
  id: string;
};

export default function DeleteAction({ id }: DeleteActionProps) {
    const queryClient = useQueryClient();

    const deleteMutation = useMutation({
        mutationFn: (id: string) => deleteItem(id),
        onSuccess: (id) => {
            console.log('Item deleted:', id);
            queryClient.invalidateQueries({ queryKey: ["items"] });
        },
        onError: (error) => {
            console.error('Error deleting item:', error);
        },
    });

    const onClick = (id: undefined | string) => {
        if (!id) {
            console.log('id is undefined');
            return;
        }
        deleteMutation.mutate(id ?? '');
    };

    return (
        <Button danger type="link" onClick={() => onClick(id)}>
            <DeleteOutlined />
        </Button>
    );
}
