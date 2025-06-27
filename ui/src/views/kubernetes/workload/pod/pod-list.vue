<template>
  <div class="pod-list-view">
    <h1>Pod List</h1>
    <!-- Basic table structure for listing pods -->
    <el-table :data="pods" style="width: 100%">
      <el-table-column prop="name" label="Name" />
      <el-table-column prop="namespace" label="Namespace" />
      <el-table-column prop="status" label="Status" />
      <el-table-column prop="node" label="Node" />
      <el-table-column prop="ip" label="IP" />
      <el-table-column label="Actions">
        <template #default="scope">
          <el-button size="small" @click="viewPodDetails(scope.row)">Details</el-button>
          <el-button size="small" type="danger" @click="deletePod(scope.row)">Delete</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- Placeholder for pagination if needed -->
    <!-- <el-pagination
      layout="prev, pager, next"
      :total="50">
    </el-pagination> -->
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { ElTable, ElTableColumn, ElButton } from 'element-plus';
// import { getPods } from '@/api/kubernetes/pods'; // Assuming API function

// Define a type for a Pod item for basic display
interface Pod {
  name: string;
  namespace: string;
  status: string;
  node: string;
  ip: string;
  // Add other relevant properties as needed
}

const pods = ref<Pod[]>([]);

// Placeholder data - replace with actual API call
const samplePods: Pod[] = [
  { name: 'nginx-deployment-12345-abcde', namespace: 'default', status: 'Running', node: 'node-1', ip: '10.244.1.2' },
  { name: 'my-app-pod-67890-fghij', namespace: 'kube-system', status: 'Pending', node: 'node-2', ip: '10.244.2.3' },
  { name: 'another-pod-instance-54321-klmno', namespace: 'default', status: 'Succeeded', node: 'node-1', ip: '10.244.1.4' },
];

onMounted(async () => {
  // Simulate API call
  // try {
  //   const response = await getPods({ clusterId: 'current-cluster-id' }); // Pass necessary params
  //   pods.value = response.data.items; // Adjust based on actual API response structure
  // } catch (error) {
  //   console.error('Failed to fetch pods:', error);
  //   // Fallback to sample data or show error message
  pods.value = samplePods;
  // }
});

const viewPodDetails = (pod: Pod) => {
  console.log('View details for:', pod.name);
  // Implement navigation to pod detail page or show a modal
  // router.push(`/homepage/kubernetes/workload/pods/${pod.namespace}/${pod.name}/detail`);
};

const deletePod = (pod: Pod) => {
  console.log('Delete pod:', pod.name);
  // Implement delete functionality, possibly with confirmation
  // Call API and then refresh list or remove from local array
};

</script>

<style scoped>
.pod-list-view {
  padding: 20px;
}
</style>
