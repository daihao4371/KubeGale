<template>
  <el-form :model="form" :rules="rules" ref="deploymentFormRef" label-width="120px" label-position="top">
    <el-tabs v-model="activeTab">
      <!-- Basic Information Tab -->
      <el-tab-pane label="基本信息" name="basic">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Deployment名称" prop="metadata.name">
              <el-input v-model="form.metadata.name" placeholder="小写字母、数字、'-' 或 '.'" :disabled="isEditMode"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="副本数" prop="spec.replicas">
              <el-input-number v-model="form.spec.replicas" :min="0" placeholder="实例数量"></el-input-number>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="标签 (Labels)" prop="metadata.labels">
           <VueJsoneditor v-model="form.metadata.labels" :options="jsonEditorOptions" height="150px" />
           <div class="form-item-help">应用于Deployment和Pod模板的标签，JSON格式。</div>
        </el-form-item>

        <el-form-item label="选择器 (Selector MatchLabels)" prop="spec.selector.matchLabels">
           <VueJsoneditor v-model="form.spec.selector.matchLabels" :options="jsonEditorOptions" height="150px" />
           <div class="form-item-help">用于匹配Pod的标签，JSON格式。通常与Pod模板标签一致。</div>
        </el-form-item>

      </el-tab-pane>

      <!-- Pod Template Tab -->
      <el-tab-pane label="Pod模板" name="podTemplate">
        <el-form-item label="Pod标签 (Template Labels)" prop="spec.template.metadata.labels">
            <VueJsoneditor v-model="form.spec.template.metadata.labels" :options="jsonEditorOptions" height="150px" />
            <div class="form-item-help">将应用于创建的Pod的标签，JSON格式。</div>
        </el-form-item>

        <!-- Containers Configuration -->
        <div v-for="(container, index) in form.spec.template.spec.containers" :key="index" class="container-section">
          <h4>容器 #{{ index + 1 }}
            <el-button type="danger" link icon="delete" @click="removeContainer(index)" v-if="form.spec.template.spec.containers.length > 1">移除容器</el-button>
          </h4>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item :label="`容器名称`" :prop="`spec.template.spec.containers[${index}].name`" :rules="containerRules.name">
                <el-input v-model="container.name" placeholder="小写字母、数字、'-'"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="`镜像名称`" :prop="`spec.template.spec.containers[${index}].image`" :rules="containerRules.image">
                <el-input v-model="container.image" placeholder="例如: nginx:latest"></el-input>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item :label="`镜像拉取策略`" :prop="`spec.template.spec.containers[${index}].imagePullPolicy`">
            <el-select v-model="container.imagePullPolicy" placeholder="选择拉取策略">
              <el-option label="IfNotPresent" value="IfNotPresent"></el-option>
              <el-option label="Always" value="Always"></el-option>
              <el-option label="Never" value="Never"></el-option>
            </el-select>
          </el-form-item>

          <!-- Ports, Env, VolumeMounts, Resources can be added here -->
           <el-form-item label="端口 (Ports)">
            <div v-for="(port, pIndex) in container.ports" :key="pIndex" class="sub-item-row">
              <el-input-number v-model="port.containerPort" placeholder="容器端口" :min="1" :max="65535" style="width: 120px; margin-right: 8px;"></el-input-number>
              <el-input v-model="port.name" placeholder="端口名称 (可选)" style="width: 150px; margin-right: 8px;"></el-input>
              <el-select v-model="port.protocol" placeholder="协议" style="width: 100px; margin-right: 8px;">
                <el-option label="TCP" value="TCP"></el-option>
                <el-option label="UDP" value="UDP"></el-option>
              </el-select>
              <el-button type="danger" link @click="removeContainerPort(container, pIndex)">移除</el-button>
            </div>
            <el-button type="primary" link @click="addContainerPort(container)">添加端口</el-button>
          </el-form-item>

        </div>
        <el-button type="primary" plain @click="addContainer">添加容器</el-button>

        <el-form-item label="重启策略" prop="spec.template.spec.restartPolicy">
            <el-select v-model="form.spec.template.spec.restartPolicy" placeholder="选择Pod重启策略">
              <el-option label="Always" value="Always"></el-option>
              <!-- Only Always is valid for Deployments -->
            </el-select>
        </el-form-item>
      </el-tab-pane>

      <!-- Advanced Settings (Strategy, etc.) -->
      <el-tab-pane label="高级设置" name="advanced">
         <el-form-item label="更新策略" prop="spec.strategy.type">
            <el-select v-model="form.spec.strategy.type" placeholder="选择更新策略">
              <el-option label="RollingUpdate" value="RollingUpdate"></el-option>
              <el-option label="Recreate" value="Recreate"></el-option>
            </el-select>
          </el-form-item>
          <div v-if="form.spec.strategy.type === 'RollingUpdate'">
            <el-form-item label="Max Unavailable" prop="spec.strategy.rollingUpdate.maxUnavailable">
              <el-input v-model="form.spec.strategy.rollingUpdate.maxUnavailable" placeholder="数字或百分比 (e.g., 1 or 25%)"></el-input>
            </el-form-item>
            <el-form-item label="Max Surge" prop="spec.strategy.rollingUpdate.maxSurge">
              <el-input v-model="form.spec.strategy.rollingUpdate.maxSurge" placeholder="数字或百分比 (e.g., 1 or 25%)"></el-input>
            </el-form-item>
          </div>
           <el-form-item label="最小就绪时间 (MinReadySeconds)" prop="spec.minReadySeconds">
              <el-input-number v-model="form.spec.minReadySeconds" :min="0" placeholder="Pod就绪等待秒数"></el-input-number>
            </el-form-item>
            <el-form-item label="历史版本保留数 (RevisionHistoryLimit)" prop="spec.revisionHistoryLimit">
              <el-input-number v-model="form.spec.revisionHistoryLimit" :min="0" placeholder="保留的旧ReplicaSet数量"></el-input-number>
            </el-form-item>
      </el-tab-pane>

    </el-tabs>

    <el-form-item style="margin-top: 20px;">
      <el-button type="primary" @click="handleSubmit" :loading="submitLoading">提交</el-button>
      <el-button @click="handleCancel">取消</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref, reactive, watch, onMounted } from 'vue';
import { ElMessage, FormInstance, FormRules } from 'element-plus';
// TODO: Create and import Deployment API functions
// import { createDeployment, updateDeployment } from '@/api/kubernetes/workload/deployment';
// TODO: Define proper Kubernetes Deployment type
import VueJsoneditor from 'vue3-ts-jsoneditor'; // Using a JSON editor for labels/selectors


type K8sDeployment = any; // Placeholder for actual Kubernetes Deployment type
type K8sContainer = any; // Placeholder
type K8sContainerPort = any; // Placeholder

interface Props {
  clusterId: string;
  namespaceId: string;
  deployment?: K8sDeployment | null; // For editing
}

const props = defineProps<Props>();
const emit = defineEmits(['submit', 'cancel']);

const deploymentFormRef = ref<FormInstance>();
const activeTab = ref('basic');
const isEditMode = ref(false);
const submitLoading = ref(false);

const jsonEditorOptions = ref({
  mode: 'code', // 'tree', 'form', 'code', 'text'
  mainMenuBar: false,
  // navigationBar: false,
  // statusBar: false,
});

const initialContainer = (): K8sContainer => ({
  name: '',
  image: '',
  imagePullPolicy: 'IfNotPresent',
  ports: [],
  // env: [],
  // volumeMounts: [],
  // resources: { limits: { cpu: '', memory: '' }, requests: { cpu: '', memory: '' } },
});

const initialFormState = (): K8sDeployment => ({
  apiVersion: 'apps/v1',
  kind: 'Deployment',
  metadata: {
    name: '',
    namespace: props.namespaceId, // Pre-fill from props
    labels: {},
  },
  spec: {
    replicas: 1,
    selector: {
      matchLabels: {},
    },
    template: {
      metadata: {
        labels: {},
      },
      spec: {
        containers: [initialContainer()],
        restartPolicy: 'Always',
        // volumes: [],
      },
    },
    strategy: {
      type: 'RollingUpdate',
      rollingUpdate: {
        maxUnavailable: '25%',
        maxSurge: '25%',
      },
    },
    minReadySeconds: 0,
    revisionHistoryLimit: 10,
  },
});

const form = reactive<K8sDeployment>(initialFormState());

watch(() => props.deployment, (newDeployment) => {
  if (newDeployment && newDeployment.metadata?.name) {
    isEditMode.value = true;
    // Deep merge or assign. Simple assignment might lose reactivity for nested objects.
    // A proper deep merge utility or careful assignment is needed for complex objects.
    Object.assign(form, JSON.parse(JSON.stringify(newDeployment)));
    // Ensure required nested structures exist if not present in incoming data
    if (!form.spec.strategy) form.spec.strategy = { type: 'RollingUpdate', rollingUpdate: { maxUnavailable: '25%', maxSurge: '25%' }};
    if (!form.spec.strategy.rollingUpdate && form.spec.strategy.type === 'RollingUpdate') {
        form.spec.strategy.rollingUpdate = { maxUnavailable: '25%', maxSurge: '25%' };
    }
    if (!form.spec.template?.spec?.containers?.length) {
        form.spec.template.spec.containers = [initialContainer()];
    }
     // Ensure metadata.labels and spec.selector.matchLabels are objects
    form.metadata.labels = form.metadata.labels || {};
    form.spec.selector.matchLabels = form.spec.selector.matchLabels || {};
    form.spec.template.metadata.labels = form.spec.template.metadata.labels || {};


  } else {
    isEditMode.value = false;
    Object.assign(form, initialFormState());
    form.metadata.namespace = props.namespaceId; // Ensure namespace is set for new deployments
  }
  // Sync selector with pod labels if they are empty and metadata labels exist
   if (!Object.keys(form.spec.selector.matchLabels).length && Object.keys(form.metadata.labels).length) {
    form.spec.selector.matchLabels = { ...form.metadata.labels };
  }
  if (!Object.keys(form.spec.template.metadata.labels).length && Object.keys(form.metadata.labels).length) {
    form.spec.template.metadata.labels = { ...form.metadata.labels };
  }
}, { immediate: true, deep: true });


const rules = reactive<FormRules>({
  'metadata.name': [
    { required: true, message: '请输入Deployment名称', trigger: 'blur' },
    { pattern: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$/, message: '名称不符合Kubernetes命名规范', trigger: 'blur' }
  ],
  'spec.replicas': [{ required: true, type: 'number', message: '请输入副本数', trigger: 'blur' }],
  // Add more rules for other fields as needed
});

const containerRules = {
  name: [{ required: true, message: '请输入容器名称', trigger: 'blur' }, { pattern: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/, message: '名称不符合Kubernetes命名规范', trigger: 'blur' }],
  image: [{ required: true, message: '请输入镜像名称', trigger: 'blur' }],
};

const addContainer = () => {
  form.spec.template.spec.containers.push(initialContainer());
};

const removeContainer = (index: number) => {
  if (form.spec.template.spec.containers.length > 1) {
    form.spec.template.spec.containers.splice(index, 1);
  } else {
    ElMessage.warning('至少需要一个容器');
  }
};

const addContainerPort = (container: K8sContainer) => {
  if (!container.ports) container.ports = [];
  container.ports.push({ containerPort: 80, protocol: 'TCP', name: '' });
};

const removeContainerPort = (container: K8sContainer, portIndex: number) => {
  container.ports.splice(portIndex, 1);
};


const handleSubmit = async () => {
  if (!deploymentFormRef.value) return;
  await deploymentFormRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true;
      try {
        // Ensure labels are not null/undefined before sending
        const payload = JSON.parse(JSON.stringify(form)); // Deep clone
        if (!payload.metadata.labels) payload.metadata.labels = {};
        if (!payload.spec.selector.matchLabels) payload.spec.selector.matchLabels = {};
        if (!payload.spec.template.metadata.labels) payload.spec.template.metadata.labels = {};

        // Convert replicas to number if it's a string from input
        payload.spec.replicas = Number(payload.spec.replicas);


        // Remove empty/null ports before submitting
        payload.spec.template.spec.containers.forEach((c: K8sContainer) => {
            if (c.ports) {
                c.ports = c.ports.filter((p: K8sContainerPort) => p.containerPort); // Only keep ports with a port number
                if (c.ports.length === 0) delete c.ports; // Remove empty ports array
            }
        });


        // MOCK API Call
        console.log('Submitting Deployment:', { clusterId: props.clusterId, namespace: props.namespaceId, deployment: payload });
        await new Promise(resolve => setTimeout(resolve, 1000));
        const mockResponse = { code: 0, message: isEditMode.value ? 'Deployment更新成功' : 'Deployment创建成功' };
        // END MOCK

        // let response;
        // if (isEditMode.value) {
        //   response = await updateDeployment(props.clusterId, props.namespaceId, payload.metadata.name, payload);
        // } else {
        //   response = await createDeployment(props.clusterId, props.namespaceId, payload);
        // }

        if (mockResponse.code === 0) {
          ElMessage.success(mockResponse.message);
          emit('submit');
        } else {
          ElMessage.error(mockResponse.message || '操作失败');
        }
      } catch (error: any) {
        console.error('Error submitting deployment form:', error);
        ElMessage.error(`操作失败: ${error.message || '请检查网络或联系管理员'}`);
      } finally {
        submitLoading.value = false;
      }
    } else {
      ElMessage.error('表单校验失败，请检查输入项');
      // Find the first invalid tab and switch to it
      for (const field in deploymentFormRef.value.fields) {
        const fieldInstance = deploymentFormRef.value.fields[field];
        if (fieldInstance.validateState === 'error') {
            if (field.startsWith('spec.template.spec.containers') || field.startsWith('spec.template.metadata')) {
                 activeTab.value = 'podTemplate';
            } else if (field.startsWith('spec.strategy') || field.startsWith('spec.minReadySeconds') || field.startsWith('spec.revisionHistoryLimit')) {
                 activeTab.value = 'advanced';
            } else {
                 activeTab.value = 'basic';
            }
            break;
        }
      }
      return false;
    }
  });
};

const handleCancel = () => {
  emit('cancel');
};

onMounted(() => {
    // Set initial namespace from props
    if(!isEditMode.value) {
        form.metadata.namespace = props.namespaceId;
    }
});

</script>

<style scoped>
.container-section {
  border: 1px solid #eee;
  padding: 15px;
  margin-bottom: 15px;
  border-radius: 4px;
}
.container-section h4 {
  margin-top: 0;
  margin-bottom: 15px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.form-item-help {
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
  margin-top: 4px;
}
.sub-item-row {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}
.sub-item-row .el-form-item {
  margin-bottom: 0; /* Remove default margin for nested form items */
}
</style>
