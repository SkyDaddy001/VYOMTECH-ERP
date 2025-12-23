'use client';

import React, { useState, useEffect } from 'react';
import {
  Table,
  Button,
  Card,
  Modal,
  Form,
  Input,
  InputNumber,
  Select,
  DatePicker,
  Space,
  Tag,
  message,
  Spin,
  Statistic,
  Row,
  Col,
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined,
  DollarOutlined,
  BankOutlined,
} from '@ant-design/icons';
import dayjs from 'dayjs';
import api from '@/services/api';

interface Financing {
  id: number;
  booking_id: number;
  bank_id?: number;
  loan_amount: number;
  sanctioned_amount: number;
  disbursed_amount: number;
  loan_type: string;
  interest_rate?: number;
  tenure_months?: number;
  emi_amount?: number;
  status: string;
  application_ref_no?: string;
  sanction_date?: string;
  created_at: string;
  bank?: any;
}

interface CreateFinancingRequest {
  booking_id: number;
  bank_id?: number;
  loan_amount: number;
  sanctioned_amount: number;
  loan_type: string;
  interest_rate?: number;
  tenure_months?: number;
  application_ref_no?: string;
}

const statusColors: Record<string, string> = {
  pending: 'orange',
  approved: 'blue',
  sanctioned: 'cyan',
  disbursing: 'processing',
  completed: 'success',
  rejected: 'error',
};

const loanTypes = [
  { label: 'Home Loan', value: 'Home Loan' },
  { label: 'Construction Loan', value: 'Construction Loan' },
  { label: 'Bridge Loan', value: 'Bridge Loan' },
];

export default function BankFinancingPage() {
  const [financings, setFinancings] = useState<Financing[]>([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  const [form] = Form.useForm();
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [selectedFinancing, setSelectedFinancing] = useState<Financing | null>(null);
  const [detailsVisible, setDetailsVisible] = useState(false);
  const [summaryStats, setSummaryStats] = useState<any>(null);

  useEffect(() => {
    fetchFinancings();
  }, [pagination.current, pagination.pageSize]);

  const fetchFinancings = async () => {
    setLoading(true);
    try {
      const response = (await api.financing.list({
        limit: pagination.pageSize,
        offset: (pagination.current - 1) * pagination.pageSize,
      })) as any;
      
      setFinancings((response?.data as Financing[]) || []);
      setPagination(prev => ({ ...prev, total: (response?.total as number) || 0 }));
    } catch (error) {
      message.error('Failed to load financing records');
      console.error('Error:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchFinancingSummary = async (id: number) => {
    try {
      const response = await api.financing.getSummary(id);
      setSummaryStats(response);
    } catch (error) {
      console.error('Error fetching summary:', error);
    }
  };

  const showModal = (financing?: Financing) => {
    if (financing) {
      setEditingId(financing.id);
      form.setFieldsValue({
        booking_id: financing.booking_id,
        loan_amount: financing.loan_amount,
        sanctioned_amount: financing.sanctioned_amount,
        loan_type: financing.loan_type,
        interest_rate: financing.interest_rate,
        tenure_months: financing.tenure_months,
        application_ref_no: financing.application_ref_no,
      });
    } else {
      setEditingId(null);
      form.resetFields();
    }
    setIsModalVisible(true);
  };

  const handleSave = async (values: CreateFinancingRequest) => {
    try {
      if (editingId) {
        await api.financing.update(editingId, values);
        message.success('Financing updated successfully');
      } else {
        await api.financing.create(values);
        message.success('Financing created successfully');
      }
      setIsModalVisible(false);
      form.resetFields();
      fetchFinancings();
    } catch (error) {
      message.error('Failed to save financing');
      console.error('Error:', error);
    }
  };

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: 'Delete Financing',
      content: 'Are you sure you want to delete this financing record?',
      okText: 'Delete',
      okType: 'danger',
      onOk: async () => {
        try {
          await api.financing.delete(id);
          message.success('Financing deleted successfully');
          fetchFinancings();
        } catch (error) {
          message.error('Failed to delete financing');
        }
      },
    });
  };

  const showDetails = (financing: Financing) => {
    setSelectedFinancing(financing);
    setDetailsVisible(true);
    fetchFinancingSummary(financing.id);
  };

  const columns = [
    {
      title: 'Booking ID',
      dataIndex: 'booking_id',
      key: 'booking_id',
      width: 100,
    },
    {
      title: 'Loan Type',
      dataIndex: 'loan_type',
      key: 'loan_type',
      width: 120,
    },
    {
      title: 'Loan Amount',
      dataIndex: 'loan_amount',
      key: 'loan_amount',
      render: (text: number) => `₹${text?.toLocaleString() || 0}`,
      width: 130,
    },
    {
      title: 'Sanctioned',
      dataIndex: 'sanctioned_amount',
      key: 'sanctioned_amount',
      render: (text: number) => `₹${text?.toLocaleString() || 0}`,
      width: 130,
    },
    {
      title: 'Disbursed',
      dataIndex: 'disbursed_amount',
      key: 'disbursed_amount',
      render: (text: number) => `₹${text?.toLocaleString() || 0}`,
      width: 130,
    },
    {
      title: 'Interest Rate',
      dataIndex: 'interest_rate',
      key: 'interest_rate',
      render: (text?: number) => text ? `${text}%` : '-',
      width: 100,
    },
    {
      title: 'EMI',
      dataIndex: 'emi_amount',
      key: 'emi_amount',
      render: (text?: number) => text ? `₹${text?.toLocaleString()}` : '-',
      width: 120,
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (text: string) => <Tag color={statusColors[text]}>{text}</Tag>,
      width: 100,
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_: any, record: Financing) => (
        <Space size="small">
          <Button
            type="text"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => showDetails(record)}
            title="View Details"
          />
          <Button
            type="text"
            size="small"
            icon={<EditOutlined />}
            onClick={() => showModal(record)}
            title="Edit"
          />
          <Button
            type="text"
            size="small"
            danger
            icon={<DeleteOutlined />}
            onClick={() => handleDelete(record.id)}
            title="Delete"
          />
        </Space>
      ),
      width: 100,
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Card>
        {/* Header */}
        <div style={{ marginBottom: '24px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <h1 style={{ margin: 0 }}>
              <BankOutlined /> Bank Financing Management
            </h1>
            <p style={{ margin: '8px 0 0 0', color: '#666' }}>
              Manage bank loans, disbursements, and NOC tracking
            </p>
          </div>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => showModal()}>
            Add Financing
          </Button>
        </div>

        {/* Summary Stats */}
        {summaryStats && (
          <Row gutter={16} style={{ marginBottom: '24px' }}>
            <Col xs={24} sm={12} lg={6}>
              <Statistic
                title="Total Financed"
                value={summaryStats.loan_amount}
                prefix="₹"
                precision={0}
              />
            </Col>
            <Col xs={24} sm={12} lg={6}>
              <Statistic
                title="Disbursed"
                value={summaryStats.total_disbursed}
                prefix="₹"
                precision={0}
              />
            </Col>
            <Col xs={24} sm={12} lg={6}>
              <Statistic
                title="Outstanding"
                value={summaryStats.outstanding_amount}
                prefix="₹"
                precision={0}
              />
            </Col>
            <Col xs={24} sm={12} lg={6}>
              <Statistic
                title="Disbursement %"
                value={summaryStats.disbursement_percentage}
                suffix="%"
                precision={1}
              />
            </Col>
          </Row>
        )}

        {/* Table */}
        <Spin spinning={loading}>
          <Table
            columns={columns}
            dataSource={financings.map(f => ({ ...f, key: f.id }))}
            pagination={{
              current: pagination.current,
              pageSize: pagination.pageSize,
              total: pagination.total,
              onChange: (page, pageSize) => {
                setPagination(prev => ({ ...prev, current: page, pageSize }));
              },
            }}
            scroll={{ x: 1200 }}
          />
        </Spin>
      </Card>

      {/* Create/Edit Modal */}
      <Modal
        title={editingId ? 'Edit Financing' : 'Create Financing'}
        open={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        footer={null}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSave}
        >
          <Form.Item
            label="Booking ID"
            name="booking_id"
            rules={[{ required: true, message: 'Booking ID is required' }]}
          >
            <InputNumber min={1} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            label="Loan Type"
            name="loan_type"
            rules={[{ required: true, message: 'Loan Type is required' }]}
          >
            <Select options={loanTypes} />
          </Form.Item>

          <Form.Item
            label="Loan Amount"
            name="loan_amount"
            rules={[{ required: true, message: 'Loan Amount is required' }]}
          >
            <InputNumber prefix="₹" style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            label="Sanctioned Amount"
            name="sanctioned_amount"
            rules={[{ required: true, message: 'Sanctioned Amount is required' }]}
          >
            <InputNumber prefix="₹" style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            label="Interest Rate (%)"
            name="interest_rate"
          >
            <InputNumber min={0} max={100} step={0.1} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            label="Tenure (Months)"
            name="tenure_months"
          >
            <InputNumber min={1} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            label="Application Ref No"
            name="application_ref_no"
          >
            <Input placeholder="e.g., BANKREF001" />
          </Form.Item>

          <Form.Item>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => setIsModalVisible(false)}>Cancel</Button>
              <Button type="primary" htmlType="submit">
                {editingId ? 'Update' : 'Create'}
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>

      {/* Details Drawer */}
      <Modal
        title="Financing Details"
        open={detailsVisible}
        onCancel={() => setDetailsVisible(false)}
        footer={null}
        width={600}
      >
        {selectedFinancing && (
          <div>
            <p><strong>Booking ID:</strong> {selectedFinancing.booking_id}</p>
            <p><strong>Loan Type:</strong> {selectedFinancing.loan_type}</p>
            <p><strong>Loan Amount:</strong> ₹{selectedFinancing.loan_amount?.toLocaleString()}</p>
            <p><strong>Sanctioned Amount:</strong> ₹{selectedFinancing.sanctioned_amount?.toLocaleString()}</p>
            <p><strong>Disbursed Amount:</strong> ₹{selectedFinancing.disbursed_amount?.toLocaleString()}</p>
            <p><strong>Interest Rate:</strong> {selectedFinancing.interest_rate}%</p>
            <p><strong>Tenure:</strong> {selectedFinancing.tenure_months} months</p>
            <p><strong>EMI:</strong> ₹{selectedFinancing.emi_amount?.toLocaleString()}</p>
            <p><strong>Status:</strong> <Tag color={statusColors[selectedFinancing.status]}>{selectedFinancing.status}</Tag></p>
            <p><strong>Application Ref:</strong> {selectedFinancing.application_ref_no || '-'}</p>
            <p><strong>Sanction Date:</strong> {selectedFinancing.sanction_date ? dayjs(selectedFinancing.sanction_date).format('DD-MMM-YYYY') : '-'}</p>
          </div>
        )}
      </Modal>
    </div>
  );
}
