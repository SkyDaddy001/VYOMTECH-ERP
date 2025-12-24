// Prisma-generated types re-export for type safety
// NOTE: Uncomment after running 'npm run prisma:generate'
/*
export type {
  Tenant,
  User,
  Team,
  SalesLead,
  Booking,
  Property,
  Lead,
  Applicant,
  Partner,
  Broker,
  JointApplicant,
  Document,
  DocumentCategory,
  DocumentType,
  SiteVisit,
  Possession,
  TitleClearance,
  Loan,
  BankFinancing,
  Project,
  Employee,
  Leave,
  Salary,
  Attendance,
  Payroll,
  Account,
  Voucher,
  CallLog,
  Integration,
  PaymentTransaction,
  LeadActivity,
  PasswordResetToken,
  AuthToken,
  AuditLog,
  SystemConfig,
} from '@prisma/client';
*/

// Union types for common operations
export type CreatedRecord = {
  id: string;
  tenantId: string;
  createdAt: Date;
  updatedAt?: Date;
};

export type PaginationParams = {
  page: number;
  limit: number;
};

export type SoftDeleteFilter = {
  deletedAt: null;
};

export type CRUDResult<T> = {
  success: boolean;
  data?: T;
  error?: string;
};
