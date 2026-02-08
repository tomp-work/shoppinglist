import { Button } from 'antd';
import { UpCircleOutlined } from '@ant-design/icons';
import { useMutation, useQueryClient } from '@tanstack/react-query'
import './App.css'

const moveItemUp = async (id: string) => {
    const response = await fetch(`http://localhost:1323/item/${id}/up`, {
        method: 'POST',
    });
    if (!response.ok) {
        console.log('Failed to move item up');
        throw new Error('Failed to move item up');
    }
};

type MoveUpActionProps = {
    id: string;
};

export default function MoveUpAction({ id }: MoveUpActionProps) {
    const queryClient = useQueryClient();

    const moveUpMutation = useMutation({
        mutationFn: (id: string) => moveItemUp(id),
        onSuccess: (id) => {
            console.log('Item moved up:', id);
            queryClient.invalidateQueries({ queryKey: ["items"] });
            queryClient.invalidateQueries({ queryKey: ["totalprice"] });
        },
        onError: (error) => {
            console.error('Error moving item up:', error);
        },
    });

    const onClick = (id: undefined | string) => {
        if (!id) {
            console.log('id is undefined');
            return;
        }
        moveUpMutation.mutate(id ?? '');
    };

    return (
        <Button type="link" onClick={() => onClick(id)}>
            <UpCircleOutlined />
        </Button>
    );
}
