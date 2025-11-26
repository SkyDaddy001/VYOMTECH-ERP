/**
 * Date formatting utility for consistent DD-MMM-YYYY format across the application
 * Example: 01-JAN-2025
 */

export const formatDateToDDMMMYYYY = (date: string | Date | null | undefined): string => {
  if (!date) return '-'

  try {
    const dateObj = typeof date === 'string' ? new Date(date) : date

    if (isNaN(dateObj.getTime())) {
      return '-'
    }

    const day = String(dateObj.getDate()).padStart(2, '0')
    const month = dateObj.toLocaleString('en-US', { month: 'short' }).toUpperCase()
    const year = dateObj.getFullYear()

    return `${day}-${month}-${year}`
  } catch (error) {
    return '-'
  }
}

export const formatDateToInput = (date: string | Date | null | undefined): string => {
  if (!date) return ''

  try {
    const dateObj = typeof date === 'string' ? new Date(date) : date

    if (isNaN(dateObj.getTime())) {
      return ''
    }

    const year = dateObj.getFullYear()
    const month = String(dateObj.getMonth() + 1).padStart(2, '0')
    const day = String(dateObj.getDate()).padStart(2, '0')

    return `${year}-${month}-${day}`
  } catch (error) {
    return ''
  }
}
