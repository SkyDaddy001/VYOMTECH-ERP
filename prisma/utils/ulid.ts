// ULID generation utility for Prisma models
// Install: npm install ulid

import { ulid } from 'ulid';

/**
 * Generate a new ULID (Universally Unique Lexicographically Sortable Identifier)
 * ULID format: 26 alphanumeric characters (48-bit timestamp + 80-bit randomness)
 * Better for database indexes than UUIDs or sequential integers
 * 
 * @returns {string} A new ULID
 */
export function generateId(): string {
  return ulid();
}

/**
 * Generate a ULID with optional timestamp
 * Useful for backdating records
 * 
 * @param timestamp - Optional timestamp in milliseconds
 * @returns {string} A new ULID
 */
export function generateIdWithTimestamp(timestamp?: number): string {
  return ulid(timestamp);
}

/**
 * Batch generate multiple ULIDs
 * 
 * @param count - Number of ULIDs to generate
 * @returns {string[]} Array of ULIDs
 */
export function generateIds(count: number): string[] {
  return Array.from({ length: count }, () => ulid());
}

/**
 * Validate if a string is a valid ULID format
 * 
 * @param id - String to validate
 * @returns {boolean} True if valid ULID format
 */
export function isValidUlid(id: string): boolean {
  return /^[0-7][0-9A-HJKMNP-TV-Z]{25}$/.test(id);
}

/**
 * Get timestamp from ULID (milliseconds since epoch)
 * 
 * @param id - ULID string
 * @returns {number} Milliseconds since epoch
 */
export function getTimestampFromUlid(id: string): number {
  const timeChars = id.substring(0, 10);
  const timeValue = timeChars.split('').reduce((acc, char, i) => {
    const charSet = '0123456789ABCDEFGHJKMNPQRSTVWXYZ';
    return acc * 32 + charSet.indexOf(char);
  }, 0);
  return timeValue;
}

export default {
  generateId,
  generateIdWithTimestamp,
  generateIds,
  isValidUlid,
  getTimestampFromUlid,
};
