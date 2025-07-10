<template>
  <div class="cluster-api-groups-tab">
    <div class="tab-header">
      <h4>API资源</h4>
      <el-select v-model="apiType" @change="fetchApiGroups" style="width: 150px">
        <el-option label="集群级别" value="cluster" />
        <el-option label="命名空间级别" value="namespace" />
      </el-select>
    </div>
    
    <div class="api-groups-list">
      <el-collapse v-model="activeNames">
        <el-collapse-item
          v-for="group in clusterApiGroups"
          :key="group.group"
          :title="group.group || 'core'"
          :name="group.group || 'core'"
        >
          <el-table :data="group.resources" border>
            <el-table-column prop="resource" label="资源名称" />
            <el-table-column prop="verbs" label="操作权限">
              <template #default="{ row }">
                <el-tag
                  v-for="verb in row.verbs"
                  :key="verb"
                  size="small"
                  style="margin-right: 4px"
                >
                  {{ verb }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-collapse-item>
      </el-collapse>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useClusterStore } from '../../../../stores/kubernetes/cluster'

const props = defineProps<{
  clusterId: number
}>()

const clusterStore = useClusterStore()
const apiType = ref('cluster')
const activeNames = ref<string[]>([])

const clusterApiGroups = computed(() => clusterStore.clusterApiGroups)

const fetchApiGroups = async () => {
  await clusterStore.fetchClusterApiGroups(props.clusterId, apiType.value)
}

onMounted(() => {
  fetchApiGroups()
})
</script>

<style scoped>
.tab-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.tab-header h4 {
  margin: 0;
  color: #303133;
}

.api-groups-list {
  margin-top: 20px;
}
</style> 