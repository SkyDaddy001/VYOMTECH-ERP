'use client';

import React, { useState, useEffect } from 'react';
import {
  Table,
  Button,
  Modal,
  Form,
  Select,
  Space,
  Tag,
  message,
  Spin,
  DatePicker,
} from 'antd';
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons';
import api from '@/services/api';
import dayjs from 'dayjs';

interface User {
  id: string;
  name: string;
  email: string;
}

interface Role {
  id: string;
  role_name: string;
  description: string;
}

interface UserRole {
  user_id: string;
  role_id: string;
  assigned_at: string;
  expires_at: string | null;
}

const UserRoleManagement = () => {
  const [userRoles, setUserRoles] = useState<UserRole[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [roles, setRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(false);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchUserRoles();
    fetchUsers();
    fetchRoles();
  }, []);

  const fetchUserRoles = async () => {
    setLoading(true);
    try {
      // This would be an endpoint that lists all user-role assignments
      // For now, we'll show the structure
      setUserRoles([]);
      message.info('User roles endpoint not yet implemented');
    } catch (error) {
      message.error('Failed to fetch user roles');
    } finally {
      setLoading(false);
    }
  };

  const fetchUsers = async () => {
    try {
      // Fetch users from the system
      const users = (await api.get('/api/v1/users')) as User[];
      setUsers(users || []);
    } catch (error) {
      console.error('Failed to fetch users', error);
    }
  };

  const fetchRoles = async () => {
    try {
      const roles = (await api.get('/api/v1/rbac/roles')) as Role[];
      setRoles(roles || []);
    } catch (error) {
      console.error('Failed to fetch roles', error);
    }
  };

  const handleAssignRole = async (values: Record<string, any>) => {
    try {
      await api.post(`/api/v1/rbac/users/${values.user_id}/roles`, {
        user_id: values.user_id,
        role_id: values.role_id,
        expires_at: values.expires_at
          ? values.expires_at.format('YYYY-MM-DD HH:mm:ss')
          : null,
      });
      message.success('Role assigned successfully');
      form.resetFields();
      setIsModalVisible(false);
      fetchUserRoles();
    } catch (error) {
      message.error('Failed to assign role');
      console.error(error);
    }
  };

  const handleRemoveRole = async (userId: string, roleId: string) => {
    Modal.confirm({
      title: 'Remove Role',
      content: 'Are you sure you want to remove this role from the user?',
      okText: 'Yes',
      cancelText: 'No',
      onOk: async () => {
        try {
          await api.delete(`/api/v1/rbac/users/${userId}/roles/${roleId}`);
          message.success('Role removed successfully');
          fetchUserRoles();
        } catch (error) {
          message.error('Failed to remove role');
          console.error(error);
        }
      },
    });
  };

  const columns = [
    {
      title: 'User ID',
      dataIndex: 'user_id',
      key: 'user_id',
    },
    {
      title: 'Role',
      dataIndex: 'role_name',
      key: 'role_name',
      render: (text: string) => <Tag color="blue">{text}</Tag>,
    },
    {
      title: 'Assigned At',
      dataIndex: 'assigned_at',
      key: 'assigned_at',
      render: (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: 'Expires At',
      dataIndex: 'expires_at',
      key: 'expires_at',
      render: (date: string) =>
        date ? dayjs(date).format('YYYY-MM-DD HH:mm:ss') : 'Never',
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_: any, record: UserRole) => (
        <Space>
          <Button
            danger
            size="small"
            icon={<DeleteOutlined />}
            onClick={() =>
              handleRemoveRole(record.user_id, record.role_id)
            }
          >
            Remove
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <Spin spinning={loading}>
      <div style={{ padding: '24px' }}>
        <h1>User Role Assignment</h1>

        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => setIsModalVisible(true)}
          style={{ marginBottom: '16px' }}
        >
          Assign Role to User
        </Button>

        <Table
          columns={columns}
          dataSource={userRoles}
          rowKey={(record: UserRole) => `${record.user_id}-${record.role_id}`}
          pagination={{
            pageSize: 20,
            showSizeChanger: true,
            showTotal: (total: number) => `Total ${total} assignments`,
          }}
        />

        {/* Assign Role Modal */}
        <Modal
          title="Assign Role to User"
          visible={isModalVisible}
          onOk={() => form.submit()}
          onCancel={() => {
            setIsModalVisible(false);
            form.resetFields();
          }}
        >
          <Form
            form={form}
            layout="vertical"
            onFinish={handleAssignRole}
          >
            <Form.Item
              name="user_id"
              label="User"
              rules={[{ required: true, message: 'Please select a user' }]}
            >
              <Select
                placeholder="Select a user"
                options={users.map((user: User) => ({
                  label: user.name || `User ${user.id}`,
                  value: user.id,
                }))}
              />
            </Form.Item>

            <Form.Item
              name="role_id"
              label="Role"
              rules={[{ required: true, message: 'Please select a role' }]}
            >
              <Select
                placeholder="Select a role"
                options={roles.map((role: Role) => ({
                  label: role.role_name,
                  value: role.id,
                }))}
              />
            </Form.Item>

            <Form.Item name="expires_at" label="Expiration Date (Optional)">
              <DatePicker showTime placeholder="Select expiration date" />
            </Form.Item>
          </Form>
        </Modal>
      </div>
    </Spin>
  );
};

export default UserRoleManagement;
