import { Button, Form, Input, InputNumber, type FormProps } from 'antd';
import { useMutation, useQueryClient } from '@tanstack/react-query'
import './App.css'

type Item = {
  id?: string;
  name: string;
  quantity: number;
  price: number;
}

type FieldType = {
  name?: string;
  quantity?: string;
  price: number;
};

export default function AddForm() {

  const queryClient = useQueryClient();

  const createItem = async (data: Item) => {
    const response = await fetch('http://localhost:1323/item', {
      method: "POST",
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

  const createMutation = useMutation({
    mutationFn: (data: Item) => createItem(data),
    onSuccess: (data) => {
      console.log("Item created:", data);
      queryClient.invalidateQueries({ queryKey: ["items"] });
      queryClient.invalidateQueries({ queryKey: ["listdetails"] });
    },
    onError: (error) => {
      console.error("Error creating item:", error);
    },
  });

  const onFinish: FormProps<FieldType>['onFinish'] = (values) => {
    console.log('Success:', JSON.stringify(values));
    createMutation.mutate({
      name: values.name ?? '',
      quantity: (values.quantity ?? 0) as number,
      price: (values.price ?? 0) as number,
    });
  };

  const onFinishFailed: FormProps<FieldType>['onFinishFailed'] = (errorInfo) => {
    console.log('Failed:', errorInfo);
  };

  return (
    <Form
      name="add"
      initialValues={{
        quantity: 1
      }}
      layout='inline'
      onFinish={onFinish}
      onFinishFailed={onFinishFailed}
    >
      <Form.Item<FieldType>
        label="Name"
        name="name"
        rules={[{ required: true, message: 'Please input item name!' }]}
      >
        <Input />
      </Form.Item>
      <Form.Item<FieldType>
        label="Quantity"
        name="quantity"
        rules={[{ required: true, message: 'Please input item quantity!' }]}
      >
        <InputNumber />
      </Form.Item>
      <Form.Item<FieldType>
        label="Price"
        name="price"
        rules={[{ required: true, message: 'Please input item price!' }]}
      >
        <InputNumber />
      </Form.Item>
      <Form.Item label={null}>
        <Button disabled={createMutation.isPending} type="primary" htmlType="submit" >
          Add
        </Button>
      </Form.Item>
      {createMutation.isError && <p>Something went wrong with create</p>}
      {createMutation.isSuccess && <p>Item created successfully!</p>}
    </Form >
  )
}