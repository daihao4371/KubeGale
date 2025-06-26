<template>
  <div>
    <el-table ref="multipleTableRef" :data="tableData" style="width: 100%" tooltip-effect="dark" row-key="metadata.uid">
      <el-table-column type="selection" width="55" />
      <el-table-column align="left" label="名称" prop="metadata.name" min-width="150">
        <template #default="scope">
          <span class="operate-span-primary" @click="() => handleDetail(scope.row)">{{ scope.row.metadata.name }}</span>
          <el-tooltip placement="top">
            <template #content>
              <div v-if="scope.row.metadata.labels && Object.keys(scope.row.metadata.labels).length > 0">
                <div v-for="(v, k, i) in scope.row.metadata.labels" :key="i">
                  <span> {{ k }}: {{ v }}</span>
                </div>
              </div>
              <div v-else>无标签</div>
            </template>
            <el-button size="small" type="primary" link icon="PriceTag" />
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column label="Internal IP" prop="status.addresses" min-width="120">
        <template #default="scope">
          <span>{{ getInternalIP(scope.row) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" prop="status.conditions" min-width="150">
        <template #default="scope">
          <el-tag :type="getNodeReadyStatus(scope.row) === 'Ready' ? 'success' : 'danger'">
            {{ getNodeDisplayStatus(scope.row) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="角色" prop="metadata.labels" min-width="100">
        <template #default="scope">
          <span>{{ getNodeRole(scope.row.metadata.labels) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="CPU用量" min-width="150">
        <template #default="scope">
          <div v-if="scope.row.cpuUsagePersent !== undefined">
            <div><span>{{ scope.row.cpuUsagePersent }}%</span></div>
            <div>
              <span>{{ CpuFormat(scope.row.cpuUsage) }} / {{ CpuFormat(scope.row.status.allocatable?.cpu) }} 核</span>
            </div>
          </div>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="内存用量" min-width="150">
        <template #default="scope">
          <div v-if="scope.row.memoryUsagePersent !== undefined">
            <div><span>{{ scope.row.memoryUsagePersent }}%</span></div>
            <div>
              <span>{{ giMemoryFormat(scope.row.memoryUsage) }} / {{ giMemoryFormat(scope.row.status.allocatable?.memory) }} Gi</span>
            </div>
          </div>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="版本" prop="status.nodeInfo.kubeletVersion" min-width="180" align="center">
        <template #default="scope">
          <div style="text-align: center" v-if="scope.row.status?.nodeInfo">
            <span>{{ scope.row.status.nodeInfo.kubeletVersion }}</span><br>
            <span>{{ scope.row.status.nodeInfo.containerRuntimeVersion }}</span><br>
            <span>{{ scope.row.status.nodeInfo.osImage }}</span>
          </div>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column align="left" label="存活时间" prop="metadata.creationTimestamp" width="180">
        <template #default="scope">
          <span>{{ AgeFormat(scope.row.metadata.creationTimestamp) }}</span>
        </template>
      </el-table-column>
      <el-table-column align="left" label="操作" width="180" fixed="right">
        <template #default="scope">
          <el-button size="small" type="primary" link icon="edit" @click="() => handleUpdate(scope.row)">编辑</el-button>
          <el-dropdown>
            <el-button type="primary" link icon="MoreFilled" size="small">更多</el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item icon="InfoFilled" @click="() => handleDetail(scope.row)">详情</el-dropdown-item>
                <el-dropdown-item icon="Expand" @click="() => handleShell(scope.row)">Shell</el-dropdown-item>
                <el-dropdown-item icon="Monitor" @click="() => handleMonitor(scope.row)">监控</el-dropdown-item>
                <el-dropdown-item icon="edit" @click="() => handleYAML(scope.row)">YAML</el-dropdown-item>
                <el-dropdown-item icon="Crop" @click="() => handleDrain(scope.row)">节点排空</el-dropdown-item>
                <el-dropdown-item icon="Setting" @click="() => handleSchedule(scope.row)">调度设置</el-dropdown-item>
                <!-- <el-dropdown-item icon="delete" @click="() => handleDelete(scope.row)">删除</el-dropdown-item> -->
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import type { ElTable } from 'element-plus';
import { AgeFormat } from '@/utils/age'; // Assuming path is correct relative to this file or aliased
import { giMemoryFormat, CpuFormat } from '@/utils/unitConvert'; // Assuming path

// Define Node type - ideally imported from a central types file
// This is a simplified version for the table component
interface NodeAddress {
  type: string;
  address: string;
}

interface NodeCondition {
  type: string;
  status: string;
  lastHeartbeatTime?: string;
  lastTransitionTime?: string;
  message?: string;
  reason?: string;
}

interface Node {
  metadata: {
    name: string;
    uid: string;
    labels?: Record<string, string>;
    creationTimestamp: string;
    [key: string]: any;
  };
  spec: {
    unschedulable?: boolean;
    [key: string]: any;
  };
  status: {
    addresses?: NodeAddress[];
    conditions?: NodeCondition[];
    nodeInfo?: {
      kubeletVersion: string;
      containerRuntimeVersion: string;
      osImage: string;
      [key: string]: any;
    };
    allocatable?: {
      cpu?: string;
      memory?: string;
      [key: string]: any;
    };
    [key: string]: any;
  };
  // Properties added by index.vue for metrics
  cpuUsagePersent?: number;
  cpuUsage?: string;
  memoryUsagePersent?: number;
  memoryUsage?: string;
}

interface Props {
  tableData: Node[];
}

const props = defineProps<Props>();

const emit = defineEmits([
  'update',
  'delete',
  'search',
  'detail',
  'schedule',
  'drain',
  'yaml',
  'monitor',
  'shell'
]);

const multipleTableRef = ref<InstanceType<typeof ElTable>>();

// Methods are now functions within setup
const handleUpdate = (value: Node) => emit('update', value);
// const handleDelete = (value: Node) => emit('delete', value); // Delete is commented out in template
const handleDrain = (value: Node) => emit('drain', value);
const handleSchedule = (value: Node) => emit('schedule', value);
const handleDetail = (value: Node) => emit('detail', value);
const handleYAML = (value: Node) => emit('yaml', value);
const handleMonitor = (value: Node) => emit('monitor', value);
const handleShell = (value: Node) => emit('shell', value);

const getInternalIP = (node: Node): string => {
  if (!node.status?.addresses) return '-';
  const internalIP = node.status.addresses.find(addr => addr.type === 'InternalIP');
  return internalIP ? internalIP.address : '-';
};

const getNodeReadyStatus = (node: Node): string => {
  if (!node.status?.conditions) return 'Unknown';
  const readyCondition = node.status.conditions.find(c => c.type === 'Ready');
  return readyCondition ? (readyCondition.status === 'True' ? 'Ready' : readyCondition.reason || 'NotReady') : 'Unknown';
};

const getNodeDisplayStatus = (node: Node): string => {
  let status = getNodeReadyStatus(node);
  if (node.spec?.unschedulable) {
    status += ', SchedulingDisabled';
  }
  return status;
};

const getNodeRole = (labels?: Record<string, string>): string => {
  if (!labels || typeof labels !== 'object') {
    return '<none>';
  }
  const roleKeys = Object.keys(labels).filter(key => key.startsWith('node-role.kubernetes.io/'));
  if (roleKeys.length > 0) {
    return roleKeys.map(key => key.substring(key.indexOf('/') + 1)).join(', ');
  }
  if (labels['kubernetes.io/role']) {
     return labels['kubernetes.io/role'];
  }
  // Default to Worker or check for common master/control-plane labels if needed
  if (labels['node-role.kubernetes.io/master'] !== undefined || labels['node-role.kubernetes.io/control-plane'] !== undefined) {
    return 'Control-Plane'; // Or 'Master'
  }
  return 'Worker'; // A common default, though nodes can have no explicit role
};

</script>

<style scoped>
.el-dropdown-link {
  cursor: pointer;
  color: var(--el-color-primary);
  display: flex;
  align-items: center;
}
.operate-span-primary {
  color: var(--el-color-primary);
  cursor: pointer;
  font-weight: bold;
}
</style>
