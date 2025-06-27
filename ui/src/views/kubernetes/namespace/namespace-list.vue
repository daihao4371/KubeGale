<template>
  <div class="namespace-list-view">
    <h1>Namespace List</h1>
    <!-- Button to create a new namespace -->
    <el-button type="primary" @click="openCreateNamespaceDialog" style="margin-bottom: 20px;">
      Create Namespace
    </el-button>

    <!-- Basic table structure for listing namespaces -->
    <el-table :data="namespaces" style="width: 100%">
      <el-table-column prop="name" label="Name" />
      <el-table-column prop="status" label="Status" />
      <el-table-column prop="age" label="Age" />
      <el-table-column label="Actions">
        <template #default="scope">
          <el-button size="small" @click="viewNamespaceDetails(scope.row)">Details</el-button>
          <el-button size="small" @click="editNamespace(scope.row)">Edit</el-button>
          <el-button size="small" type="danger" @click="deleteNamespace(scope.row)">Delete</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- Placeholder for pagination -->
    <!-- <el-pagination
      layout="prev, pager, next"
      :total="50">
    </el-pagination> -->

    <!-- Dialog for creating/editing a namespace -->
    <!-- <el-dialog v-model="dialogVisible" :title="dialogTitle">
      <namespace-form ref="namespaceForm" :namespace="selectedNamespace" @submit="handleFormSubmit" />
    </el-dialog> -->
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { ElTable, ElTableColumn, ElButton, ElDialog } from 'element-plus';
// import NamespaceForm from './namespace-form.vue'; // Assuming a form component
// import { getNamespaces, createNamespace, updateNamespace, deleteNamespace } from '@/api/kubernetes/namespace'; // Assuming API functions

// Define a type for a Namespace item
interface Namespace {
  name: string;
  status: string;
  age: string; // Or Date object, depending on how you want to handle it
  // Add other relevant properties
}

const namespaces = ref<Namespace[]>([]);
const dialogVisible = ref(false);
const dialogTitle = ref('');
const selectedNamespace = ref<Namespace | null>(null);
// const namespaceForm = ref(null); // Ref for the form component

// Placeholder data - replace with actual API call
const sampleNamespaces: Namespace[] = [
  { name: 'default', status: 'Active', age: '365d' },
  { name: 'kube-system', status: 'Active', age: '365d' },
  { name: 'kube-public', status: 'Active', age: '365d' },
  { name: 'my-namespace', status: 'Active', age: '10d' },
];

onMounted(async () => {
  // Simulate API call
  // try {
  //   const response = await getNamespaces({ clusterId: 'current-cluster-id' }); // Pass cluster context
  //   namespaces.value = response.data.items; // Adjust based on API
  // } catch (error) {
  //   console.error('Failed to fetch namespaces:', error);
  namespaces.value = sampleNamespaces;
  // }
});

const openCreateNamespaceDialog = () => {
  selectedNamespace.value = null;
  dialogTitle.value = 'Create Namespace';
  dialogVisible.value = true;
  // if (namespaceForm.value) {
  //   (namespaceForm.value as any).resetForm(); // Assuming form has a reset method
  // }
};

const editNamespace = (namespace: Namespace) => {
  selectedNamespace.value = { ...namespace }; // Create a copy to avoid modifying the original directly
  dialogTitle.value = 'Edit Namespace';
  dialogVisible.value = true;
};

const viewNamespaceDetails = (namespace: Namespace) => {
  console.log('View details for:', namespace.name);
  // Implement navigation or modal display for details
  // router.push(`/homepage/kubernetes/namespace/${namespace.name}/detail`);
};

const deleteNamespaceConfirmation = (namespace: Namespace) => {
  // ElMessageBox.confirm(`Are you sure you want to delete namespace "${namespace.name}"?`, 'Warning', {
  //   confirmButtonText: 'OK',
  //   cancelButtonText: 'Cancel',
  //   type: 'warning',
  // }).then(async () => {
  //   try {
  //     // await deleteNamespace(namespace.name, { clusterId: 'current-cluster-id' });
  //     // ElMessage.success('Namespace deleted successfully');
  //     // Refresh list
  //     // onMounted(); // or remove from local array
  namespaces.value = namespaces.value.filter(ns => ns.name !== namespace.name); // Optimistic update for now
  console.log('Deleted namespace (simulated):', namespace.name);
  //   } catch (error) {
  //     // ElMessage.error('Failed to delete namespace');
  //     console.error('Failed to delete namespace:', error);
  //   }
  // }).catch(() => {
  //   // ElMessage.info('Delete canceled');
  // });
};

const deleteNamespace = (namespace: Namespace) => {
  console.log('Attempting to delete namespace:', namespace.name);
  deleteNamespaceConfirmation(namespace); // Directly call for now without ElMessageBox
};

const handleFormSubmit = async (formData: Namespace) => {
  // try {
  //   if (selectedNamespace.value && selectedNamespace.value.name) { // Editing existing
  //     // await updateNamespace(selectedNamespace.value.name, formData, { clusterId: 'current-cluster-id' });
  //     // ElMessage.success('Namespace updated successfully');
  console.log('Simulated update for namespace:', formData);
  //   } else { // Creating new
  //     // await createNamespace(formData, { clusterId: 'current-cluster-id' });
  //     // ElMessage.success('Namespace created successfully');
  console.log('Simulated create for namespace:', formData);
  //   }
  //   dialogVisible.value = false;
  //   // Refresh list
  //   // onMounted();
  // } catch (error) {
  //   // ElMessage.error(`Failed to ${selectedNamespace.value ? 'update' : 'create'} namespace`);
  //   console.error('Namespace form submission error:', error);
  // }
};

</script>

<style scoped>
.namespace-list-view {
  padding: 20px;
}
</style>
