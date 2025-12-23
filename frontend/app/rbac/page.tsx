'use client';

import React, { useState, useEffect } from 'react';
import {
  Table,
  Button,
  Modal,
  Form,
  Input,
  Select,
  Checkbox,
  Space,
  Tabs,
  Tag,
  message,
  Spin,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  LockOutlined,
} from '@ant-design/icons';
import api from '@/services/api';

interface Role {
  id: string;
  role_name: string;
  name?: string;  // Support both role_name and name fields
  description: string;
  permissions: string[];
  is_active: boolean;
}

interface Permission {
  id: string;
  permission_name: string;
  description: string;
}

const RBACDashboard = () => {
  const [roles, setRoles] = useState<Role[]>([]);
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedRole, setSelectedRole] = useState<Role | null>(null);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [isPermissionModalVisible, setIsPermissionModalVisible] =
    useState(false);
  const [form] = Form.useForm();
  const [permissionForm] = Form.useForm();

  // Fetch roles on component mount
  useEffect(() => {
    fetchRoles();
    fetchPermissions();
  }, []);

  const fetchRoles = async () => {
    setLoading(true);
    try {
      const roles = (await api.get('/api/v1/rbac/roles')) as Role[];
      setRoles(roles || []);
      message.success('Roles loaded successfully');
    } catch (error) {
      message.error('Failed to fetch roles');
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  const fetchPermissions = async () => {
    try {
      const perms = (await api.get('/api/v1/rbac/permissions')) as Permission[];
      setPermissions(perms || []);
    } catch (error) {
      console.error('Failed to fetch permissions', error);
    }
  };

  const handleCreateRole = async (values: Record<string, any>) => {
    try {
      await api.post('/api/v1/rbac/roles', {
        role_name: values.role_name,
        description: values.description,
        is_active: values.is_active || true,
      });
      message.success('Role created successfully');
      form.resetFields();
      setIsModalVisible(false);
      fetchRoles();
    } catch (error) {
      message.error('Failed to create role');
      console.error(error);
    }
  };

  const handleAssignPermissions = async (values: Record<string, any>) => {
    if (!selectedRole) {
      message.error('Please select a role');
      return;
    }

    try {
      await api.put(`/api/v1/rbac/roles/${selectedRole.id}/permissions`, {
        permission_ids: values.permission_ids || [],
      });
      message.success('Permissions assigned successfully');
      permissionForm.resetFields();
      setIsPermissionModalVisible(false);
      fetchRoles();
    } catch (error) {
      message.error('Failed to assign permissions');
      console.error(error);
    }
  };

  const handleDeleteRole = async (roleId: string) => {
    Modal.confirm({
      title: 'Delete Role',
      content: 'Are you sure you want to delete this role?',
      okText: 'Yes',
      cancelText: 'No',
      onOk: async () => {
        try {
          await api.delete(`/api/v1/rbac/roles/${roleId}`);
          message.success('Role deleted successfully');
          fetchRoles();
        } catch (error) {
          message.error('Failed to delete role');
          console.error(error);
        }
      },
    });
  };

  const rolesColumns = [
    {
      title: 'Role Name',
      dataIndex: 'name',
      key: 'name',
      render: (text: string) => <strong>{text}</strong>,
    },
    {
      title: 'Description',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: 'Permissions',
      dataIndex: 'permissions',
      key: 'permissions',
      render: (permissions: string[]) => (
        <>
          {permissions && permissions.length > 0 ? (
            <>
              {permissions.slice(0, 3).map((perm: string) => (
                <Tag key={perm} color="blue">
                  {perm}
                </Tag>
              ))}
              {permissions.length > 3 && (
                <Tag color="default">+{permissions.length - 3} more</Tag>
              )}
            </>
          ) : (
            <Tag color="red">No permissions</Tag>
          )}
        </>
      ),
    },
    {
      title: 'Status',
      dataIndex: 'is_active',
      key: 'is_active',
      render: (isActive: boolean) => (
        <Tag color={isActive ? 'green' : 'red'}>
          {isActive ? 'Active' : 'Inactive'}
        </Tag>
      ),
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_: any, record: Role) => (
        <Space>
          <Button
            type="primary"
            size="small"
            icon={<LockOutlined />}
            onClick={() => {
              setSelectedRole(record);
              setIsPermissionModalVisible(true);
            }}
          >
            Permissions
          </Button>
          <Button
            danger
            size="small"
            icon={<DeleteOutlined />}
            onClick={() => handleDeleteRole(record.id)}
          >
            Delete
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <Spin spinning={loading}>
      <div style={{ padding: '24px' }}>
        <h1>RBAC Management</h1>

        <Tabs
          items={[
            {
              key: 'roles',
              label: 'Roles',
              children: (
                <div>
                  <Button
                    type="primary"
                    icon={<PlusOutlined />}
                    onClick={() => setIsModalVisible(true)}
                    style={{ marginBottom: '16px' }}
                  >
                    Create Role
                  </Button>

                  <Table
                    columns={rolesColumns}
                    dataSource={roles}
                    rowKey="id"
                    pagination={{
                      pageSize: 10,
                      showSizeChanger: true,
                      showTotal: (total: number) => `Total ${total} roles`,
                    }}
                  />
                </div>
              ),
            },
            {
              key: 'permissions',
              label: 'Permissions',
              children: (
                <Table
                  columns={[
                    {
                      title: 'Permission Code',
                      dataIndex: 'permission_name',
                      key: 'permission_name',
                    },
                    {
                      title: 'Description',
                      dataIndex: 'description',
                      key: 'description',
                    },
                    {
                      title: 'Resource',
                      dataIndex: 'resource',
                      key: 'resource',
                    },
                    {
                      title: 'Action',
                      dataIndex: 'action',
                      key: 'action',
                    },
                  ]}
                  dataSource={permissions}
                  rowKey="id"
                  pagination={{
                    pageSize: 10,
                    showSizeChanger: true,
                  }}
                />
              ),
            },
          ]}
        />

        {/* Create Role Modal */}
        <Modal
          title="Create New Role"
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
            onFinish={handleCreateRole}
          >
            <Form.Item
              name="role_name"
              label="Role Name"
              rules={[
                { required: true, message: 'Please enter role name' },
              ]}
            >
              <Input placeholder="e.g., Manager, Admin, Agent" />
            </Form.Item>

            <Form.Item
              name="description"
              label="Description"
              rules={[
                {
                  required: true,
                  message: 'Please enter role description',
                },
              ]}
            >
              <Input.TextArea placeholder="Describe the role's purpose" />
            </Form.Item>

            <Form.Item
              name="is_active"
              valuePropName="checked"
              initialValue={true}
            >
              <Checkbox>Active</Checkbox>
            </Form.Item>
          </Form>
        </Modal>

        {/* Assign Permissions Modal */}
        <Modal
          title={`Assign Permissions - ${selectedRole?.name || ''}`}
          visible={isPermissionModalVisible}
          onOk={() => permissionForm.submit()}
          onCancel={() => {
            setIsPermissionModalVisible(false);
            permissionForm.resetFields();
            setSelectedRole(null);
          }}
        >
          <Form
            form={permissionForm}
            layout="vertical"
            onFinish={handleAssignPermissions}
          >
            <Form.Item
              name="permission_ids"
              label="Select Permissions"
              rules={[
                {
                  required: true,
                  message: 'Please select at least one permission',
                },
              ]}
            >
              <Checkbox.Group
                options={permissions.map((perm: Permission) => ({
                  label: `${perm.permission_name} - ${perm.description}`,
                  value: perm.id,
                }))}
              />
            </Form.Item>
          </Form>
        </Modal>
      </div>
    </Spin>
  );
};

export default RBACDashboard;
