<template>
  <div class="history-management">
    <!-- Filter Form -->
    <el-form :inline="true" :model="filters" class="filter-form" @submit.native.prevent>
      <el-form-item label="Command">
        <el-input v-model="filters.command" placeholder="Filter by command"></el-input>
      </el-form-item>

      <el-form-item label="Directory">
        <el-input v-model="filters.cwd" placeholder="Filter by directory"></el-input>
      </el-form-item>

      <el-form-item label="Hostname">
        <el-input v-model="filters.hostname" placeholder="Filter by hostname"></el-input>
      </el-form-item>

      <el-form-item label="Exit Code">
        <el-select v-model="filters.exit" placeholder="Filter by exit code" clearable>
          <el-option label="Success (0)" :value="0"></el-option>
          <el-option label="Error (non-zero)" :value="1"></el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="Time Range">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="to"
          start-placeholder="Start date"
          end-placeholder="End date"
          value-format="timestamp"
          :picker-options="{ firstDayOfWeek: 1 }"
        ></el-date-picker>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="fetchHistoryEntries">Search</el-button>
        <el-button @click="resetFilters">Reset</el-button>
      </el-form-item>
    </el-form>

    <!-- Action Buttons -->
    <div class="action-buttons">
      <el-button
        type="danger"
        :disabled="selectedEntries.length === 0"
        @click="confirmDelete"
      >
        Delete Selected ({{ selectedEntries.length }})
      </el-button>
    </div>

    <!-- Data Table -->
    <el-table
      v-loading="loading"
      :data="historyEntries"
      style="width: 100%"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55"></el-table-column>

      <el-table-column prop="timestamp" label="Timestamp" width="180">
        <template slot-scope="scope">
          {{ formatTimestamp(scope.row.timestamp) }}
        </template>
      </el-table-column>

      <el-table-column prop="command" label="Command" show-overflow-tooltip></el-table-column>

      <el-table-column prop="cwd" label="Directory" show-overflow-tooltip width="200"></el-table-column>

      <el-table-column prop="duration" label="Duration" width="100">
        <template slot-scope="scope">
          {{ formatDuration(scope.row.duration) }}
        </template>
      </el-table-column>

      <el-table-column prop="exit" label="Exit Code" width="100">
        <template slot-scope="scope">
          <el-tag :type="scope.row.exit === 0 ? 'success' : 'danger'">
            {{ scope.row.exit }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column prop="hostname" label="Hostname" width="270"></el-table-column>
    </el-table>

    <!-- Pagination -->
    <div class="pagination-container">
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="currentPage"
        :page-sizes="[10, 20, 50, 100]"
        :page-size="pageSize"
        layout="total, sizes, prev, pager, next, jumper"
        :total="totalEntries"
      ></el-pagination>
    </div>

    <!-- Delete Confirmation Dialog -->
    <el-dialog
      title="Confirm Deletion"
      :visible.sync="deleteDialogVisible"
      width="30%"
    >
      <span>Are you sure you want to delete {{ selectedEntries.length }} selected entries?</span>
      <span slot="footer" class="dialog-footer">
        <el-button @click="deleteDialogVisible = false">Cancel</el-button>
        <el-button type="danger" @click="deleteSelectedEntries">Delete</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import moment from 'moment';
import api from '@/services/api';

export default {
  name: 'HistoryManagement',
  data() {
    return {
      historyEntries: [],
      selectedEntries: [],
      loading: false,
      deleteDialogVisible: false,
      totalEntries: 0,
      currentPage: 1,
      pageSize: 20,
      dateRange: null,
      filters: {
        command: '',
        cwd: '',
        hostname: '',
        exit: null,
        start_time: null,
        end_time: null
      }
    };
  },
  created() {
    this.fetchHistoryEntries();
  },
  methods: {
    async fetchHistoryEntries() {
      this.loading = true;

      // Update time filters from date range
      if (this.dateRange && this.dateRange.length === 2) {
        // Convert milliseconds to nanoseconds (multiply by 1,000,000)
        this.filters.start_time = this.dateRange[0] * 1000000;
        this.filters.end_time = this.dateRange[1] * 1000000;
      } else {
        this.filters.start_time = null;
        this.filters.end_time = null;
      }

      // Add pagination
      const params = {
        ...this.filters,
        limit: this.pageSize,
        offset: (this.currentPage - 1) * this.pageSize
      };

      // Remove null/empty values
      Object.keys(params).forEach(key => {
        if (params[key] === null || params[key] === '') {
          delete params[key];
        }
      });

      try {
        const response = await api.getHistoryEntries(params);
        this.historyEntries = response.data;

        // For simplicity, we're assuming the total count is the number of returned entries
        // In a real application, you might want to return the total count from the API
        this.totalEntries = response.data.length >= this.pageSize
          ? this.currentPage * this.pageSize + 1 // Indicate there are more
          : (this.currentPage - 1) * this.pageSize + response.data.length;
      } catch (error) {
        console.error('Error fetching history entries:', error);
        this.$message.error('Failed to load history entries');
      } finally {
        this.loading = false;
      }
    },

    async deleteSelectedEntries() {
      this.loading = true;
      const ids = this.selectedEntries.map(entry => entry.id);

      try {
        await api.batchDeleteHistoryEntries(ids);
        this.$message.success('Selected entries deleted successfully');
        this.fetchHistoryEntries();
        this.deleteDialogVisible = false;
      } catch (error) {
        console.error('Error deleting history entries:', error);
        this.$message.error('Failed to delete history entries');
      } finally {
        this.loading = false;
      }
    },

    confirmDelete() {
      this.deleteDialogVisible = true;
    },

    handleSelectionChange(val) {
      this.selectedEntries = val;
    },

    handleSizeChange(size) {
      this.pageSize = size;
      this.fetchHistoryEntries();
    },

    handleCurrentChange(page) {
      this.currentPage = page;
      this.fetchHistoryEntries();
    },

    resetFilters() {
      this.filters = {
        command: '',
        cwd: '',
        hostname: '',
        exit: null,
        start_time: null,
        end_time: null
      };
      this.dateRange = null;
      this.fetchHistoryEntries();
    },

    formatTimestamp(timestamp) {
      // Convert nanoseconds to milliseconds (divide by 1,000,000)
      return moment(timestamp / 1000000).format('YYYY-MM-DD HH:mm:ss');
    },

    formatDuration(duration) {
      // Convert nanoseconds to a readable format
      const ns = duration % 1000000;
      const ms = Math.floor(duration / 1000000) % 1000;
      const seconds = Math.floor(duration / 1000000000) % 60;
      const minutes = Math.floor(duration / (1000000000 * 60)) % 60;
      const hours = Math.floor(duration / (1000000000 * 60 * 60));

      if (hours > 0) {
        return `${hours}h ${minutes}m ${seconds}s`;
      } else if (minutes > 0) {
        return `${minutes}m ${seconds}s`;
      } else if (seconds > 0) {
        return `${seconds}.${ms.toString().padStart(3, '0')}s`;
      } else if (ms > 0) {
        return `${ms}.${Math.floor(ns / 1000).toString().padStart(3, '0')}ms`;
      } else {
        return `${ns}ns`;
      }
    }
  }
};
</script>

<style scoped>
.history-management {
  padding: 20px;
}

.filter-form {
  margin-bottom: 20px;
  background-color: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
}

.action-buttons {
  margin-bottom: 15px;
  display: flex;
  justify-content: flex-end;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
