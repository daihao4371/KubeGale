<template>
  <el-form :model="form" :rules="rules" ref="clusterFormRef" label-width="140px">
    <el-form-item label="集群名称" prop="name">
      <el-input v-model="form.name" placeholder="请输入集群名称"></el-input>
    </el-form-item>

    <el-form-item label="集群类型" prop="kube_type">
      <el-radio-group v-model="form.kube_type">
        <el-radio :label="1">KubeConfig</el-radio>
        <el-radio :label="2">Agent (暂未支持)</el-radio>
      </el-radio-group>
    </el-form-item>

    <el-form-item v-if="form.kube_type === 1" label="KubeConfig" prop="kube_config">
      <el-input
        type="textarea"
        :rows="10"
        v-model="form.kube_config"
        placeholder="请粘贴 KubeConfig 内容 (YAML格式)"
      ></el-input>
    </el-form-item>

    <el-form-item label="API Server地址" prop="api_address">
      <el-input v-model="form.api_address" placeholder="例如: https://192.168.0.1:6443"></el-input>
      <div class="form-item-help">当使用 KubeConfig 时，若不填写，将尝试从 KubeConfig 中自动解析。</div>
    </el-form-item>

    <el-form-item label="Prometheus URL" prop="prometheus_url">
      <el-input v-model="form.prometheus_url" placeholder="例如: http://prometheus.example.com"></el-input>
    </el-form-item>

    <el-form-item label="Prometheus认证类型" prop="prometheus_auth_type">
      <el-radio-group v-model="form.prometheus_auth_type">
        <el-radio :label="0">None</el-radio>
        <el-radio :label="1">Basic Auth</el-radio>
      </el-radio-group>
    </el-form-item>

    <template v-if="form.prometheus_auth_type === 1">
      <el-form-item label="Prometheus 用户名" prop="prometheus_user">
        <el-input v-model="form.prometheus_user" placeholder="请输入Prometheus认证用户名"></el-input>
      </el-form-item>
      <el-form-item label="Prometheus 密码" prop="prometheus_pwd">
        <el-input type="password" v-model="form.prometheus_pwd" placeholder="请输入Prometheus认证密码" show-password></el-input>
      </el-form-item>
    </template>

    <el-form-item>
      <el-button type="primary" @click="handleSubmit">提交</el-button>
      <el-button @click="handleCancel">取消</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref, watch, reactive, onMounted } from 'vue';
import { ElMessage, FormInstance, FormRules } from 'element-plus';
import { CreateCluster, UpdateCluster } from '@/api/kubernetes/cluster/k8sCluster';
import type { ClusterData } from '@/api/kubernetes/cluster/k8sCluster'; // Assuming this will be more detailed

interface Props {
  cluster?: ClusterData | null; // For editing
}

const props = defineProps<Props>();
const emit = defineEmits(['submit', 'cancel']);

const clusterFormRef = ref<FormInstance>();

const initialFormState: ClusterData = {
  name: '',
  kube_type: 1, // Default to KubeConfig
  kube_config: '',
  api_address: '',
  prometheus_url: '',
  prometheus_auth_type: 0, // Default to None
  prometheus_user: '',
  prometheus_pwd: '',
};

const form = reactive<ClusterData>({ ...initialFormState });
const isEditMode = ref(false);

onMounted(() => {
  if (props.cluster && props.cluster.id) {
    isEditMode.value = true;
    // Deep copy props.cluster to form
    Object.assign(form, JSON.parse(JSON.stringify(props.cluster)));
  } else {
    isEditMode.value = false;
    Object.assign(form, { ...initialFormState }); // Reset for create
  }
});

watch(() => props.cluster, (newCluster) => {
  if (newCluster && newCluster.id) {
    isEditMode.value = true;
    Object.assign(form, JSON.parse(JSON.stringify(newCluster)));
     if (form.prometheus_auth_type === undefined) { // Handle case where it might be null from backend
      form.prometheus_auth_type = 0;
    }
  } else {
    isEditMode.value = false;
    Object.assign(form, { ...initialFormState });
  }
}, { immediate: true, deep: true });


const rules = reactive<FormRules<ClusterData>>({
  name: [{ required: true, message: '请输入集群名称', trigger: 'blur' }],
  kube_type: [{ required: true, message: '请选择集群类型', trigger: 'change' }],
  kube_config: [
    {
      validator: (rule, value, callback) => {
        if (form.kube_type === 1 && !value) {
          callback(new Error('KubeConfig类型集群必须提供KubeConfig内容'));
        } else {
          callback();
        }
      },
      trigger: 'blur',
    },
  ],
  api_address: [
    // Optional for KubeConfig type if it can be parsed, but good to have a pattern if provided
    { pattern: /^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$/, message: '请输入有效的URL', trigger: 'blur' }
  ],
  prometheus_url: [
    { pattern: /^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$/, message: '请输入有效的Prometheus URL', trigger: 'blur' }
  ],
  prometheus_user: [
      {
      validator: (rule, value, callback) => {
        if (form.prometheus_auth_type === 1 && !value) {
          callback(new Error('Prometheus用户名为必填项'));
        } else {
          callback();
        }
      },
      trigger: 'blur',
    },
  ],
  prometheus_pwd: [
    {
      validator: (rule, value, callback) => {
        if (form.prometheus_auth_type === 1 && !value) {
          callback(new Error('Prometheus密码为必填项'));
        } else {
          callback();
        }
      },
      trigger: 'blur',
    },
  ]
});

const handleSubmit = async () => {
  if (!clusterFormRef.value) return;
  await clusterFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        let response;
        const payload = { ...form };
        // Ensure numeric fields are numbers if backend expects them
        payload.kube_type = Number(payload.kube_type);
        payload.prometheus_auth_type = Number(payload.prometheus_auth_type);


        if (isEditMode.value && payload.id) {
          response = await UpdateCluster(payload);
        } else {
          // For create, ensure no ID is sent if backend auto-generates it
          const { id, ...createPayload } = payload;
          response = await CreateCluster(createPayload as ClusterData); // Type assertion
        }

        if (response.code === 0) {
          ElMessage.success(isEditMode.value ? '集群更新成功' : '集群创建成功');
          emit('submit');
        } else {
          ElMessage.error(response.message || (isEditMode.value ? '更新失败' : '创建失败'));
        }
      } catch (error) {
        console.error('Error submitting cluster form:', error);
        ElMessage.error((isEditMode.value ? '更新失败' : '创建失败'));
      }
    } else {
      ElMessage.error('请检查表单输入');
      return false;
    }
  });
};

const handleCancel = () => {
  if (clusterFormRef.value) {
    clusterFormRef.value.resetFields();
  }
  emit('cancel');
};

</script>

<style scoped>
.form-item-help {
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
}
</style>
