import axios from 'axios';

// Use relative path for API when in production (embedded mode)
// In development, use the localhost URL
const API_URL = process.env.NODE_ENV === 'production'
  ? '/api'
  : 'http://localhost:8080/api';

const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  },
  timeout: 10000
});

export default {
  // Get history entries with optional filters
  getHistoryEntries(filters = {}) {
    return apiClient.get('/history', { params: filters });
  },

  // Batch delete history entries
  batchDeleteHistoryEntries(ids) {
    return apiClient.delete('/history', { data: { ids } });
  }
};
