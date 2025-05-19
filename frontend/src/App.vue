<template>
  <div id="app">
    <div v-if="loading" class="loading-container">
      <i class="el-icon-loading"></i>
      <p>Loading application...</p>
    </div>
    <HistoryManagement v-else-if="databaseStatus === 'ok'" />
    <DatabaseNotFound v-else-if="databaseStatus === 'not_found'" :dbPath="dbPath" />
    <DatabaseNotFound v-else-if="databaseStatus === 'no_path'" />
    <div v-else class="error-container">
      <i class="el-icon-warning-outline"></i>
      <h1>Database Error</h1>
      <p>{{ errorMessage }}</p>
      <el-button type="primary" @click="refreshPage">
        <i class="el-icon-refresh"></i> Refresh Page
      </el-button>
    </div>
  </div>
</template>

<script>
import HistoryManagement from './components/HistoryManagement.vue'
import DatabaseNotFound from './components/DatabaseNotFound.vue'
import axios from 'axios'

export default {
  name: 'App',
  components: {
    HistoryManagement,
    DatabaseNotFound
  },
  data() {
    return {
      loading: true,
      databaseStatus: null,
      dbPath: '',
      errorMessage: ''
    }
  },
  async created() {
    try {
      // Use relative path for API when in production (embedded mode)
      const apiUrl = process.env.NODE_ENV === 'production'
        ? '/api/config/db-status'
        : 'http://localhost:8080/api/config/db-status';

      // Check the database status
      const response = await axios.get(apiUrl)
      this.databaseStatus = response.data.status
      this.dbPath = response.data.path || ''
    } catch (error) {
      console.error('Failed to check database status:', error)

      if (error.response) {
        // The request was made and the server responded with a status code
        // that falls out of the range of 2xx
        const { status, data } = error.response

        if (status === 404) {
          // Database not found or no path
          this.databaseStatus = data.status || 'not_found'
          this.dbPath = data.path || ''
        } else {
          // Other error
          this.databaseStatus = 'error'
          this.errorMessage = data.error || 'An unknown error occurred'
        }
      } else if (error.request) {
        // The request was made but no response was received
        this.databaseStatus = 'connection_error'
        this.errorMessage = 'Failed to connect to the backend server. Please make sure it is running.'
      } else {
        // Something happened in setting up the request that triggered an Error
        this.databaseStatus = 'error'
        this.errorMessage = error.message || 'An unknown error occurred'
      }
    } finally {
      this.loading = false
    }
  },
  methods: {
    refreshPage() {
      window.location.reload()
    }
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  margin: 0;
  padding: 0;
}

body {
  margin: 0;
  padding: 0;
}

.loading-container, .error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  text-align: center;
  padding: 20px;
}

.error-container {
  color: #606266;
}

.error-container i {
  font-size: 80px;
  color: #f56c6c;
  margin-bottom: 20px;
}

.error-container h1 {
  font-size: 32px;
  color: #303133;
  margin-bottom: 20px;
}

.error-container p {
  font-size: 18px;
  margin-bottom: 30px;
}

.error-container .el-button {
  padding: 12px 20px;
  font-size: 16px;
}
</style>