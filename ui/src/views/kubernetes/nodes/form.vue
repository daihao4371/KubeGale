<template>
  <div style="overflow: hidden; width: 100%">
    <el-form ref="nodeFormRef" :model="localFormData" :rules="rules" label-width="100px">
      <el-form-item label="名称" prop="metadata.name">
        <el-input v-model="localFormData.metadata.name" :disabled="isUpdateMode" autocomplete="off" />
      </el-form-item>

      <el-form-item label="标签" prop="metadata.labels">
        <VueJsoneditor v-model="localFormData.metadata.labels" :options="jsonEditorOptions" height="200px" />
        <div class="form-item-help">Node labels in JSON format.</div>
      </el-form-item>

      <el-form-item label="注解" prop="metadata.annotations">
         <VueJsoneditor v-model="localFormData.metadata.annotations" :options="jsonEditorOptions" height="200px" />
         <div class="form-item-help">Node annotations in JSON format.</div>
      </el-form-item>
    </el-form>
    <div class="dialog-footer" style="text-align: right; margin-top: 20px;">
      <el-button size="small" @click="closeDialog">取 消</el-button>
      <el-button size="small" type="primary" @click="enterDialog">确 定</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue';
import type { FormInstance, FormRules } from 'element-plus';
import VueJsoneditor from 'vue3-ts-jsoneditor';
// Assuming a simplified Node type for the form. Ideally, this would come from your k8s types.
interface NodeFormData {
  metadata: {
    name: string;
    labels?: Record<string, string>;
    annotations?: Record<string, string>;
    [key: string]: any;
  };
  // Include other fields if this form were to edit more aspects of a Node
  [key: string]: any;
}

interface Props {
  form: NodeFormData; // Expect the full node object or relevant parts
  optype: 'create' | 'update'; // Operation type
}

const props = defineProps<Props>();
const emit = defineEmits(['close', 'enter']);

const nodeFormRef = ref<FormInstance>();

// Use a local reactive copy of the prop to allow modification by VueJsoneditor
const localFormData = reactive<NodeFormData>({
  metadata: { name: '' , labels: {}, annotations: {}},
});

const jsonEditorOptions = ref({
  mode: 'code' as const, // 'tree', 'form', 'code', 'text'
  mainMenuBar: false,
  // navigationBar: false, // Hides navigation bar
  // statusBar: false,     // Hides status bar
});


const isUpdateMode = computed(() => props.optype === 'update');

// Initialize localFormData when props.form changes
watch(() => props.form, (newForm) => {
  if (newForm) {
    // Deep copy to avoid modifying the prop directly, especially for nested objects
    localFormData.metadata.name = newForm.metadata?.name || '';
    localFormData.metadata.labels = JSON.parse(JSON.stringify(newForm.metadata?.labels || {}));
    localFormData.metadata.annotations = JSON.parse(JSON.stringify(newForm.metadata?.annotations || {}));
    // Copy other parts of the form if necessary
  } else {
    // Reset if form is null/undefined (e.g. for create mode if not pre-filled)
    localFormData.metadata.name = '';
    localFormData.metadata.labels = {};
    localFormData.metadata.annotations = {};
  }
}, { immediate: true, deep: true });


const rules = reactive<FormRules<NodeFormData>>({
  'metadata.name': [{ required: true, message: '请输入名称', trigger: 'blur' }],
  // Rules for labels and annotations are not typically needed if using JSON editor,
  // unless specific validation of keys/values is required.
});

const closeDialog = () => {
  if (nodeFormRef.value) {
    nodeFormRef.value.resetFields(); // Resets to initial values if prop was set, or clears if not
  }
  // Reset local state manually if resetFields isn't sufficient due to deep objects
  localFormData.metadata.labels = JSON.parse(JSON.stringify(props.form?.metadata?.labels || {}));
  localFormData.metadata.annotations = JSON.parse(JSON.stringify(props.form?.metadata?.annotations || {}));
  emit('close');
};

const enterDialog = async () => {
  if (!nodeFormRef.value) return;
  try {
    const valid = await nodeFormRef.value.validate();
    if (valid) {
      // Emit a deep copy of the local form data
      emit('enter', JSON.parse(JSON.stringify(localFormData)));
    }
  } catch (error) {
    // Validation failed (validate promise rejects on failure)
    console.log('Form validation failed:', error);
  }
};

</script>

<style scoped>
.form-item-help {
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
  margin-top: 4px;
}
/* Add any specific styles if needed */
</style>
