// Prisma Client wrapper with enhanced utilities
import { PrismaClient } from '@prisma/client';
import { generateId } from './utils/ulid';

const prisma = new PrismaClient({
  log: process.env.NODE_ENV === 'development' 
    ? ['query', 'error', 'warn'] 
    : ['error'],
});

// Middleware to automatically inject IDs for create operations
prisma.$use(async (params, next) => {
  // Auto-inject IDs for create operations
  if (params.action === 'create' && !params.args.data?.id) {
    params.args.data.id = generateId();
  }
  
  // Auto-inject IDs for createMany operations
  if (params.action === 'createMany' && params.args.data) {
    params.args.data = (Array.isArray(params.args.data) 
      ? params.args.data 
      : [params.args.data]
    ).map((item: any) => ({
      ...item,
      id: item.id || generateId(),
    }));
  }

  return next(params);
});

// Graceful shutdown
process.on('SIGINT', async () => {
  await prisma.$disconnect();
  process.exit(0);
});

process.on('SIGTERM', async () => {
  await prisma.$disconnect();
  process.exit(0);
});

export default prisma;
export { generateId };
