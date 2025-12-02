import { apiClient } from './api'
import { Booking, BookingPayment, BookingMetrics, BookingTimeline } from '@/types/bookings'

export const bookingsService = {
  // Bookings
  async getBookings(): Promise<Booking[]> {
    return apiClient.get<Booking[]>('/api/v1/bookings')
  },

  async getBooking(id: string): Promise<Booking> {
    return apiClient.get<Booking>(`/api/v1/bookings/${id}`)
  },

  async createBooking(data: Partial<Booking>): Promise<Booking> {
    return apiClient.post<Booking>('/api/v1/bookings', data)
  },

  async updateBooking(id: string, data: Partial<Booking>): Promise<Booking> {
    return apiClient.put<Booking>(`/api/v1/bookings/${id}`, data)
  },

  async deleteBooking(id: string): Promise<void> {
    return apiClient.delete(`/api/v1/bookings/${id}`)
  },

  async updateBookingStatus(id: string, status: string): Promise<Booking> {
    return apiClient.put<Booking>(`/api/v1/bookings/${id}/status`, { status })
  },

  // Payments
  async getBookingPayments(bookingId: string): Promise<BookingPayment[]> {
    return apiClient.get<BookingPayment[]>(`/api/v1/bookings/${bookingId}/payments`)
  },

  async addPayment(bookingId: string, data: Partial<BookingPayment>): Promise<BookingPayment> {
    return apiClient.post<BookingPayment>(`/api/v1/bookings/${bookingId}/payments`, data)
  },

  async updatePayment(bookingId: string, paymentId: string, data: Partial<BookingPayment>): Promise<BookingPayment> {
    return apiClient.put<BookingPayment>(`/api/v1/bookings/${bookingId}/payments/${paymentId}`, data)
  },

  async deletePayment(bookingId: string, paymentId: string): Promise<void> {
    return apiClient.delete(`/api/v1/bookings/${bookingId}/payments/${paymentId}`)
  },

  // Metrics
  async getMetrics(): Promise<BookingMetrics> {
    return apiClient.get<BookingMetrics>('/api/v1/bookings/metrics/overview')
  },

  async getTimeline(): Promise<BookingTimeline[]> {
    return apiClient.get<BookingTimeline[]>('/api/v1/bookings/timeline')
  },
}
