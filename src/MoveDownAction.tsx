import { Button } from 'antd';
import { DownCircleOutlined } from '@ant-design/icons';
import { useMutation, useQueryClient } from '@tanstack/react-query'
import './App.css'

const moveItemDown = async (id: string) => {
    const response = await fetch(`http://localhost:1323/item/${id}/down`, {
        method: 'POST',
    });
    if (!response.ok) {
        console.log('Failed to move item down');
        throw new Error('Failed to move item down');
    }
};

type MoveDownActionProps = {
  id: string;
};

export default function MoveDownAction({ id }: MoveDownActionProps) {
    const queryClient = useQueryClient();

    const moveDownMutation = useMutation({
        mutationFn: (id: string) => moveItemDown(id),
        onSuccess: (id) => {
            console.log('Item moved down:', id);
            queryClient.invalidateQueries({ queryKey: ["items"] });
            queryClient.invalidateQueries({ queryKey: ["listdetails"] });
        },
        onError: (error) => {
            console.error('Error moving item down:', error);
        },
    });

    const onClick = (id: undefined | string) => {
        if (!id) {
            console.log('id is undefined');
            return;
        }
        moveDownMutation.mutate(id ?? '');
    };

    return (
        <Button type="link" onClick={() => onClick(id)}>
            <DownCircleOutlined />
        </Button>
    );
}
