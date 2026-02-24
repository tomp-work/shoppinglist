import { Button, Form, Input, type FormProps } from 'antd';
import { useMutation } from '@tanstack/react-query'
import './App.css'

type Email = {
  emailAddress: string;
}

type FieldType = {
  emailAddress: string;
};

const sendEmail = async (data: Email) => {
  const response = await fetch('http://localhost:1323/list/send', {
    method: "POST",
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  if (!response.ok) {
    throw new Error('Failed to send email');
  }
  return;
};

export default function SendEmailForm() {

  const sendEmailMutation = useMutation({
    mutationFn: (data: Email) => sendEmail(data),
    onSuccess: () => {
      console.log("Email sent");
    },
    onError: (error) => {
      console.error("Error sending email:", error);
    },
  });

  const onFinish: FormProps<FieldType>['onFinish'] = (values) => {
    console.log('Success:', JSON.stringify(values));
    sendEmailMutation.mutate({
      emailAddress: values.emailAddress ?? 0,
    });
  };

  const onFinishFailed: FormProps<FieldType>['onFinishFailed'] = (errorInfo) => {
    console.log('Failed:', errorInfo);
  };

  return (
    <Form
      name="sendEmail"
      layout='inline'
      onFinish={onFinish}
      onFinishFailed={onFinishFailed}
    >
      <Form.Item<FieldType>
        label="Email Address"
        name="emailAddress"
        rules={[
          { required: true, message: 'Please input email address!' },
          { type: "email", message: "Please enter a valid email" }
        ]}
      >
        <Input />
      </Form.Item>
      <Form.Item label={null}>
        <Button disabled={sendEmailMutation.isPending} type="primary" htmlType="submit" >
          Send Email
        </Button>
      </Form.Item>
      { sendEmailMutation.isError && <p>Something went wrong with sending email.</p> }
  { sendEmailMutation.isSuccess && <p>Email sent successfully!</p> }
    </Form >
  )
}