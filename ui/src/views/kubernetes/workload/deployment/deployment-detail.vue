<template>
  <div class="deployment-detail-page" v-loading="loading">
    <el-page-header @back="goBack" :content="pageTitle" style="margin-bottom: 20px;"></el-page-header>

    <div v-if="deployment">
      <el-card class="box-card">
        <template #header>
          <div class="card-header">
            <span>基本信息</span>
            <div>
              <!-- Action buttons like Edit YAML, Scale, Rollback can go here -->
               <el-button type="primary" @click="handleEditYaml" icon="edit">编辑YAML</el-button>
            </div>
          </div>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="名称">{{ deployment.metadata?.name }}</el-descriptions-item>
          <el-descriptions-item label="命名空间">{{ deployment.metadata?.namespace }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDateTime(deployment.metadata?.creationTimestamp) }}</el-descriptions-item>
          <el-descriptions-item label="UID">{{ deployment.metadata?.uid }}</el-descriptions-item>
          <el-descriptions-item label="期望副本数">{{ deployment.spec?.replicas }}</el-descriptions-item>
          <el-descriptions-item label="更新策略">
            <el-tag>{{ deployment.spec?.strategy?.type }}</el-tag>
          </el-descriptions-item>
           <el-descriptions-item v-if="deployment.spec?.strategy?.type === 'RollingUpdate'" label="Max Unavailable">
            {{ deployment.spec?.strategy?.rollingUpdate?.maxUnavailable }}
          </el-descriptions-item>
          <el-descriptions-item v-if="deployment.spec?.strategy?.type === 'RollingUpdate'" label="Max Surge">
            {{ deployment.spec?.strategy?.rollingUpdate?.maxSurge }}
          </el-descriptions-item>
          <el-descriptions-item label="选择器 (Selector)">
            <div v-if="deployment.spec?.selector?.matchLabels">
                <el-tag v-for="(value, key) in deployment.spec.selector.matchLabels" :key="key" style="margin-right: 5px;">
                    {{ key }}: {{ value }}
                </el-tag>
            </div>
            <span v-else>-</span>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <el-tabs v-model="activeTab" style="margin-top: 20px;">
        <el-tab-pane label="状态" name="status">
          <el-card>
            <template #header><span>运行状态</span></template>
            <el-descriptions :column="2" border>
                <el-descriptions-item label="当前副本数">{{ deployment.status?.replicas || 0 }}</el-descriptions-item>
                <el-descriptions-item label="可用副本数">{{ deployment.status?.availableReplicas || 0 }}</el-descriptions-item>
                <el-descriptions-item label="就绪副本数">{{ deployment.status?.readyReplicas || 0 }}</el-descriptions-item>
                <el-descriptions-item label="更新后副本数">{{ deployment.status?.updatedReplicas || 0 }}</el-descriptions-item>
                <el-descriptions-item label="碰撞次数">{{ deployment.status?.observedGeneration }}</el-descriptions-item>
            </el-descriptions>
            <div v-if="deployment.status?.conditions?.length" style="margin-top:20px;">
                <h4>Conditions</h4>
                <el-table :data="deployment.status.conditions" stripe border size="small">
                    <el-table-column prop="type" label="Type"></el-table-column>
                    <el-table-column prop="status" label="Status"></el-table-column>
                    <el-table-column prop="lastUpdateTime" label="Last Update" :formatter="row => formatDateTime(row.lastUpdateTime)"></el-table-column>
                    <el-table-column prop="lastTransitionTime" label="Last Transition" :formatter="row => formatDateTime(row.lastTransitionTime)"></el-table-column>
                    <el-table-column prop="reason" label="Reason"></el-table-column>
                    <el-table-column prop="message" label="Message" min-width="200">
                         <template #default="scope">
                            <pre style="white-space: pre-wrap; word-break: break-all;">{{scope.row.message}}</pre>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
          </el-card>
        </el-tab-pane>

        <el-tab-pane label="关联Pods" name="pods">
          <el-card>
            <template #header><span>管理的Pods</span></template>
            <!-- Placeholder for Pods list component -->
            <p>Pods列表将在此处显示 (例如, 引用 PodList 组件，并传递label selector)</p>
            <!-- <PodsTable :clusterId="clusterId" :namespace="namespace" :selector="deployment.spec.selector.matchLabels" /> -->
          </el-card>
        </el-tab-pane>

        <el-tab-pane label="历史版本 (ReplicaSets)" name="replicasets">
          <el-card>
            <template #header><span>历史版本 (ReplicaSets)</span></template>
            <!-- Placeholder for ReplicaSets list -->
            <p>ReplicaSets列表将在此处显示</p>
          </el-card>
        </el-tab-pane>

        <el-tab-pane label="事件 (Events)" name="events">
          <el-card>
            <template #header><span>相关事件</span></template>
            <!-- Placeholder for Events list -->
            <p>事件列表将在此处显示</p>
          </el-card>
        </el-tab-pane>
      </el-tabs>

    </div>
    <el-empty v-else-if="!loading && !deployment" description="未找到Deployment信息"></el-empty>

    <!-- YAML Editor Dialog -->
    <el-dialog title="编辑YAML" v-model="yamlEditorVisible" width="70%" :close-on-click-modal="false">
        <div v-if="yamlData">
             <VueJsoneditor v-model="editableYamlData" :options="jsonEditorOptions" height="500px" />
        </div>
        <template #footer>
            <el-button @click="yamlEditorVisible = false">取消</el-button>
            <el-button type="primary" @click="submitYamlUpdate" :loading="yamlSubmitLoading">应用</el-button>
        </template>
    </el-dialog>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
// TODO: Create and import Deployment API functions
// import { getDeploymentByName, updateDeploymentWithYaml } from '@/api/kubernetes/workload/deployment';
// TODO: Define Deployment type
import VueJsoneditor from 'vue3-ts-jsoneditor'; // Using a JSON editor for labels/selectors


type K8sDeployment = any; // Placeholder for actual Kubernetes Deployment type

const route = useRoute();
const router = useRouter();

const deployment = ref<K8sDeployment | null>(null);
const loading = ref(true);
const activeTab = ref('status');

const yamlEditorVisible = ref(false);
const yamlData = ref<object | null>(null); // For raw YAML editing
const editableYamlData = ref<object>({});
const yamlSubmitLoading = ref(false);
const jsonEditorOptions = ref({
  mode: 'code',
  mainMenuBar: false,
});


const clusterId = computed(() => route.params.clusterId as string);
const namespace = computed(() => route.params.namespace as string);
const deploymentName = computed(() => route.params.name as string);

const pageTitle = computed(() => deployment.value ? `Deployment: ${deployment.value.metadata?.name}` : 'Deployment详情');

const formatDateTime = (dateTime: string | undefined) => {
  if (!dateTime) return '';
  return new Date(dateTime).toLocaleString();
};

const fetchDeploymentDetails = async () => {
  loading.value = true;
  try {
    // const response = await getDeploymentByName(clusterId.value, namespace.value, deploymentName.value);
    // MOCK DATA
    await new Promise(resolve => setTimeout(resolve, 300));
    const mockDeployment = {
      apiVersion: "apps/v1",
      kind: "Deployment",
      metadata: {
        name: deploymentName.value,
        namespace: namespace.value,
        uid: "mock-uid-" + Math.random().toString(36).substring(7),
        creationTimestamp: new Date().toISOString(),
        labels: { app: deploymentName.value }
      },
      spec: {
        replicas: 2,
        selector: { matchLabels: { app: deploymentName.value }},
        template: {
          metadata: { labels: { app: deploymentName.value }},
          spec: {
            containers: [{ name: 'nginx', image: 'nginx:latest', ports: [{containerPort: 80}]}]
          }
        },
        strategy: { type: 'RollingUpdate', rollingUpdate: {maxUnavailable: '25%', maxSurge: '25%'}}
      },
      status: {
        replicas: 2,
        readyReplicas: 2,
        availableReplicas: 2,
        updatedReplicas: 2,
        observedGeneration: 1,
        conditions: [
            {type: "Available", status: "True", lastUpdateTime: new Date().toISOString(), lastTransitionTime: new Date().toISOString(), reason: "MinimumReplicasAvailable", message: "Deployment has minimum availability."},
            {type: "Progressing", status: "True", lastUpdateTime: new Date().toISOString(), lastTransitionTime: new Date().toISOString(), reason: "NewReplicaSetAvailable", message: `ReplicaSet "${deploymentName.value}-xyz" has successfully progressed.`}
        ]
      }
    };
    const response = { code: 0, data: mockDeployment };
    // END MOCK
    if (response.code === 0) {
      deployment.value = response.data;
      yamlData.value = JSON.parse(JSON.stringify(response.data)); // For YAML editor
      editableYamlData.value = JSON.parse(JSON.stringify(response.data));
    } else {
      // ElMessage.error(response.message || '获取Deployment详情失败');
      deployment.value = null;
    }
  } catch (error) {
    console.error('Error fetching deployment details:', error);
    ElMessage.error('获取Deployment详情接口错误');
    deployment.value = null;
  } finally {
    loading.value = false;
  }
};

const goBack = () => {
  // router.push({ name: 'DeploymentList' }); // Assuming this is the route name for the list
  // Or go back in history if appropriate
  router.back();
};

const handleEditYaml = () => {
    if(yamlData.value) {
        editableYamlData.value = JSON.parse(JSON.stringify(yamlData.value)); // Reset to original for editing
        yamlEditorVisible.value = true;
    } else {
        ElMessage.warning("无YAML数据可编辑");
    }
};

const submitYamlUpdate = async () => {
    yamlSubmitLoading.value = true;
    try {
        // const res = await updateDeploymentWithYaml(clusterId.value, namespace.value, deploymentName.value, editableYamlData.value);
        // MOCK
        await new Promise(resolve => setTimeout(resolve, 500));
        const res = { code: 0, message: "通过YAML更新成功"};
        // END MOCK

        if (res.code === 0) {
            ElMessage.success(res.message || "更新成功");
            yamlEditorVisible.value = false;
            fetchDeploymentDetails(); // Refresh details
        } else {
            ElMessage.error(res.message || "更新失败");
        }
    } catch (error: any) {
        ElMessage.error("更新失败: " + error.message);
    } finally {
        yamlSubmitLoading.value = false;
    }
};

onMounted(() => {
  if (clusterId.value && namespace.value && deploymentName.value) {
    fetchDeploymentDetails();
  } else {
    ElMessage.error('缺少必要的路由参数 (clusterId, namespace, name)');
    loading.value = false;
  }
});
</script>

<style scoped>
.deployment-detail-page {
  padding: 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.el-card {
  margin-bottom: 20px;
}
</style>
