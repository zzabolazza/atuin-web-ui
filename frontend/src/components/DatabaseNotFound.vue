<template>
  <div class="database-not-found">
    <div class="error-container">
      <div class="error-icon">
        <i class="el-icon-warning-outline"></i>
      </div>
      <h1>Database Not Found</h1>
      <p class="error-message">
        Could not locate your Atuin database file. Please make sure Atuin is installed and has been used to record command history.
      </p>
      <p v-if="dbPath" class="db-path">
        Looking for database at: <code>{{ dbPath }}</code>
      </p>
      <div class="instructions">
        <h2>Troubleshooting Steps:</h2>
        <ol>
          <li>
            <strong>Verify Atuin is installed:</strong>
            <pre>atuin --version</pre>
          </li>
          <li>
            <strong>Check your Atuin database location:</strong>
            <pre>atuin info</pre>
            Look for the line that says <code>client db path: "/path/to/history.db"</code>
          </li>
          <li>
            <strong>Make sure you've recorded some history:</strong>
            <pre>atuin history list</pre>
          </li>
          <li>
            <strong>Restart the application</strong> after confirming your Atuin database exists.
          </li>
        </ol>
      </div>
      <div class="actions">
        <el-button type="primary" @click="refreshPage">
          <i class="el-icon-refresh"></i> Refresh Page
        </el-button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'DatabaseNotFound',
  props: {
    dbPath: {
      type: String,
      default: ''
    }
  },
  methods: {
    refreshPage() {
      window.location.reload();
    }
  }
}
</script>

<style scoped>
.database-not-found {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f7fa;
  padding: 20px;
}

.error-container {
  max-width: 800px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 40px;
  text-align: center;
}

.error-icon {
  font-size: 80px;
  color: #e6a23c;
  margin-bottom: 20px;
}

h1 {
  font-size: 32px;
  color: #303133;
  margin-bottom: 20px;
}

.error-message {
  font-size: 18px;
  color: #606266;
  margin-bottom: 15px;
}

.db-path {
  font-size: 16px;
  color: #606266;
  margin-bottom: 30px;
  background-color: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  word-break: break-all;
}

.instructions {
  text-align: left;
  margin-bottom: 30px;
}

h2 {
  font-size: 22px;
  color: #303133;
  margin-bottom: 15px;
}

ol {
  padding-left: 20px;
}

li {
  margin-bottom: 15px;
  line-height: 1.5;
}

pre {
  background-color: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  margin: 8px 0;
  font-family: monospace;
  overflow-x: auto;
}

code {
  background-color: #f5f7fa;
  padding: 2px 4px;
  border-radius: 3px;
  font-family: monospace;
}

.actions {
  margin-top: 30px;
}

.el-button {
  padding: 12px 20px;
  font-size: 16px;
}

.el-button i {
  margin-right: 5px;
}
</style>
