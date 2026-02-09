import { Button, Form, InputNumber, type FormProps } from 'antd';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import './App.css'

type ListDetails = {
  spendingLimit: number;
}

type FieldType = {
  spendingLimit: number;
};

const getListDetails = async () => {
  const res = await fetch("http://localhost:1323/list");
  if (!res.ok) {
    console.log('Failed to fetch list details');
    throw new Error('Failed to fetch list details');
  }
  return res.json();
};

const updateListDetails = async (data: ListDetails) => {
  const response = await fetch('http://localhost:1323/list', {
    method: "PUT",
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  if (!response.ok) {
    throw new Error('Failed to create item');
  }
  return response.json();
};

export default function SpendingLimitForm() {
  const queryClient = useQueryClient();

  const query = useQuery({
    queryKey: ['list'],
    queryFn: () => {
      return getListDetails();
    }
  });

  const updateMutation = useMutation({
    mutationFn: (data: ListDetails) => updateListDetails(data),
    onSuccess: (data) => {
      console.log("Item created:", data);
      queryClient.invalidateQueries({ queryKey: ["listdetails"] });
    },
    onError: (error) => {
      console.error("Error creating item:", error);
    },
  });

  const onFinish: FormProps<FieldType>['onFinish'] = (values) => {
    console.log('Success:', JSON.stringify(values));
    updateMutation.mutate({
      spendingLimit: values.spendingLimit ?? 0,
    });
  };

  const onFinishFailed: FormProps<FieldType>['onFinishFailed'] = (errorInfo) => {
    console.log('Failed:', errorInfo);
  };

  if (query.isPending) {
    return 'Loading...';
  }
  if (query.error) {
    return 'An error has occurred: ' + query.error.message;
  }

  return (
    <Form
      name="spendingLimit"
      initialValues={{
        spendingLimit: query.data.spendingLimit
      }}
      layout='inline'
      onFinish={onFinish}
      onFinishFailed={onFinishFailed}
    >
      <Form.Item<FieldType>
        label="Spending Limit"
        name="spendingLimit"
        rules={[{ required: true, message: 'Please input spending limit!' }]}
      >
        <InputNumber />
      </Form.Item>
      <Form.Item label={null}>
        <Button disabled={updateMutation.isPending} type="primary" htmlType="submit" >
          Set Spending Limit
        </Button>
      </Form.Item>
      {updateMutation.isError && <p>Something went wrong with update of list details.</p>}
      {updateMutation.isSuccess && <p>List details updated successfully!</p>}
    </Form >
  )
}