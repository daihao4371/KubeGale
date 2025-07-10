<template>
  <el-dialog
    v-model="dialogVisible"
    :title="isEdit ? '编辑 Deployment' : '创建 Deployment'"
    width="800px"
    :before-close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
      v-loading="loading"
    >
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="名称" prop="metadata.name">
            <el-input
              v-model="form.metadata.name"
              placeholder="请输入 Deployment 名称"
              :disabled="isEdit"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="副本数" prop="spec.replicas">
            <el-input-number
              v-model="form.spec.replicas"
              :min="0"
              :max="100"
              placeholder="副本数"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item label="镜像" prop="spec.template.spec.containers">
        <div
          v-for="(container, index) in form.spec.template.spec.containers"
          :key="index"
          class="container-item"
        >
          <el-row :gutter="10">
            <el-col :span="6">
              <el-input
                v-model="container.name"
                placeholder="容器名称"
                size="small"
              />
            </el-col>
            <el-col :span="12">
              <el-input
                v-model="container.image"
                placeholder="镜像地址"
                size="small"
              />
            </el-col>
            <el-col :span="4">
              <el-button
                type="danger"
                size="small"
                @click="removeContainer(index)"
                :disabled="form.spec.template.spec.containers.length === 1"
              >
                删除
              </el-button>
            </el-col>
          </el-row>
        </div>
        <el-button type="primary" size="small" @click="addContainer">
          添加容器
        </el-button>
      </el-form-item>

      <el-form-item label="端口映射">
        <div
          v-for="(container, containerIndex) in form.spec.template.spec.containers"
          :key="containerIndex"
          class="port-section"
        >
          <h4>{{ container.name || '容器' + (containerIndex + 1) }} 端口</h4>
          <div
            v-for="(port, portIndex) in container.ports"
            :key="portIndex"
            class="port-item"
          >
            <el-row :gutter="10">
              <el-col :span="6">
                <el-input
                  v-model="port.name"
                  placeholder="端口名称"
                  size="small"
                />
              </el-col>
              <el-col :span="6">
                <el-input-number
                  v-model="port.containerPort"
                  :min="1"
                  :max="65535"
                  placeholder="容器端口"
                  size="small"
                />
              </el-col>
              <el-col :span="6">
                <el-select v-model="port.protocol" size="small">
                  <el-option label="TCP" value="TCP" />
                  <el-option label="UDP" value="UDP" />
                </el-select>
              </el-col>
              <el-col :span="4">
                <el-button
                  type="danger"
                  size="small"
                  @click="removePort(containerIndex, portIndex)"
                >
                  删除
                </el-button>
              </el-col>
            </el-row>
          </div>
          <el-button
            type="primary"
            size="small"
            @click="addPort(containerIndex)"
          >
            添加端口
          </el-button>
        </div>
      </el-form-item>

      <el-form-item label="环境变量">
        <div
          v-for="(container, containerIndex) in form.spec.template.spec.containers"
          :key="containerIndex"
          class="env-section"
        >
          <h4>{{ container.name || '容器' + (containerIndex + 1) }} 环境变量</h4>
          <div
            v-for="(env, envIndex) in container.env"
            :key="envIndex"
            class="env-item"
          >
            <el-row :gutter="10">
              <el-col :span="8">
                <el-input
                  v-model="env.name"
                  placeholder="变量名"
                  size="small"
                />
              </el-col>
              <el-col :span="8">
                <el-input
                  v-model="env.value"
                  placeholder="变量值"
                  size="small"
                />
              </el-col>
              <el-col :span="4">
                <el-button
                  type="danger"
                  size="small"
                  @click="removeEnv(containerIndex, envIndex)"
                >
                  删除
                </el-button>
              </el-col>
            </el-row>
          </div>
          <el-button
            type="primary"
            size="small"
            @click="addEnv(containerIndex)"
          >
            添加环境变量
          </el-button>
        </div>
      </el-form-item>

      <el-form-item label="资源限制">
        <div
          v-for="(container, containerIndex) in form.spec.template.spec.containers"
          :key="containerIndex"
          class="resource-section"
        >
          <h4>{{ container.name || '容器' + (containerIndex + 1) }} 资源</h4>
          <el-row :gutter="20">
            <el-col :span="12">
              <h5>请求资源</h5>
              <el-row :gutter="10">
                <el-col :span="12">
                  <el-input
                    v-model="container.resources.requests.cpu"
                    placeholder="CPU (如: 100m)"
                    size="small"
                  />
                </el-col>
                <el-col :span="12">
                  <el-input
                    v-model="container.resources.requests.memory"
                    placeholder="内存 (如: 128Mi)"
                    size="small"
                  />
                </el-col>
              </el-row>
            </el-col>
            <el-col :span="12">
              <h5>限制资源</h5>
              <el-row :gutter="10">
                <el-col :span="12">
                  <el-input
                    v-model="container.resources.limits.cpu"
                    placeholder="CPU (如: 200m)"
                    size="small"
                  />
                </el-col>
                <el-col :span="12">
                  <el-input
                    v-model="container.resources.limits.memory"
                    placeholder="内存 (如: 256Mi)"
                    size="small"
                  />
                </el-col>
              </el-row>
            </el-col>
          </el-row>
        </div>
      </el-form-item>

      <el-form-item label="更新策略" prop="spec.strategy.type">
        <el-radio-group v-model="form.spec.strategy.type">
          <el-radio label="RollingUpdate">滚动更新</el-radio>
          <el-radio label="Recreate">重新创建</el-radio>
        </el-radio-group>
        <div v-if="form.spec.strategy.type === 'RollingUpdate'" class="strategy-config">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-input
                v-model="form.spec.strategy.rollingUpdate.maxUnavailable"
                placeholder="最大不可用数量 (如: 25%)"
                size="small"
              />
            </el-col>
            <el-col :span="12">
              <el-input
                v-model="form.spec.strategy.rollingUpdate.maxSurge"
                placeholder="最大超出数量 (如: 25%)"
                size="small"
              />
            </el-col>
          </el-row>
        </div>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="submitting">
        {{ isEdit ? '更新' : '创建' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, nextTick } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { deploymentApi } from '@/api/kubernetes/workload'
import type { Deployment, CreateWorkloadRequest, UpdateWorkloadRequest } from '@/types/kubernetes/workload'

// Props
interface Props {
  visible: boolean
  deployment?: Deployment | null
  clusterId: number
  namespace: string
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  deployment: null,
  clusterId: 0,
  namespace: ''
})

// Emits
const emit = defineEmits<{
  'update:visible': [value: boolean]
  'success': []
}>()

// 响应式数据
const formRef = ref<FormInstance>()
const loading = ref(false)
const submitting = ref(false)

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const isEdit = computed(() => !!props.deployment)

// 表单数据
const form = reactive({
  metadata: {
    name: '',
    namespace: props.namespace,
    labels: {} as Record<string, string>,
    annotations: {} as Record<string, string>
  },
  spec: {
    replicas: 1,
    selector: {
      matchLabels: {} as Record<string, string>
    },
    template: {
      metadata: {
        labels: {} as Record<string, string>
      },
      spec: {
        containers: [
          {
            name: '',
            image: '',
            ports: [] as Array<{
              name: string
              containerPort: number
              protocol: string
            }>,
            env: [] as Array<{
              name: string
              value: string
            }>,
            resources: {
              requests: {
                cpu: '',
                memory: ''
              },
              limits: {
                cpu: '',
                memory: ''
              }
            }
          }
        ]
      }
    },
    strategy: {
      type: 'RollingUpdate' as 'RollingUpdate' | 'Recreate',
      rollingUpdate: {
        maxUnavailable: '25%',
        maxSurge: '25%'
      }
    }
  }
})

// 表单验证规则
const rules: FormRules = {
  'metadata.name': [
    { required: true, message: '请输入 Deployment 名称', trigger: 'blur' },
    { pattern: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/, message: '名称只能包含小写字母、数字和连字符', trigger: 'blur' }
  ],
  'spec.replicas': [
    { required: true, message: '请输入副本数', trigger: 'blur' },
    { type: 'number', min: 0, max: 100, message: '副本数必须在 0-100 之间', trigger: 'blur' }
  ],
  'spec.template.spec.containers': [
    { required: true, message: '至少需要一个容器', trigger: 'blur' }
  ]
}

// 方法
const initForm = () => {
  if (props.deployment) {
    // 编辑模式，填充现有数据
    Object.assign(form, JSON.parse(JSON.stringify(props.deployment)))
  } else {
    // 创建模式，重置表单
    form.metadata.name = ''
    form.metadata.namespace = props.namespace
    form.spec.replicas = 1
    form.spec.template.spec.containers = [
      {
        name: '',
        image: '',
        ports: [],
        env: [],
        resources: {
          requests: { cpu: '', memory: '' },
          limits: { cpu: '', memory: '' }
        }
      }
    ]
  }
}

const addContainer = () => {
  form.spec.template.spec.containers.push({
    name: '',
    image: '',
    ports: [],
    env: [],
    resources: {
      requests: { cpu: '', memory: '' },
      limits: { cpu: '', memory: '' }
    }
  })
}

const removeContainer = (index: number) => {
  if (form.spec.template.spec.containers.length > 1) {
    form.spec.template.spec.containers.splice(index, 1)
  }
}

const addPort = (containerIndex: number) => {
  form.spec.template.spec.containers[containerIndex].ports.push({
    name: '',
    containerPort: 80,
    protocol: 'TCP'
  })
}

const removePort = (containerIndex: number, portIndex: number) => {
  form.spec.template.spec.containers[containerIndex].ports.splice(portIndex, 1)
}

const addEnv = (containerIndex: number) => {
  form.spec.template.spec.containers[containerIndex].env.push({
    name: '',
    value: ''
  })
}

const removeEnv = (containerIndex: number, envIndex: number) => {
  form.spec.template.spec.containers[containerIndex].env.splice(envIndex, 1)
}

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    const requestData = {
      cluster_id: props.clusterId,
      namespace: props.namespace,
      content: form
    }

    if (isEdit.value) {
      await deploymentApi.updateDeployment({
        ...requestData,
        name: form.metadata.name
      })
      ElMessage.success('更新成功')
    } else {
      await deploymentApi.createDeployment(requestData)
      ElMessage.success('创建成功')
    }

    emit('success')
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
  } finally {
    submitting.value = false
  }
}

const handleClose = () => {
  dialogVisible.value = false
}

// 监听对话框显示状态
watch(() => props.visible, (visible) => {
  if (visible) {
    nextTick(() => {
      initForm()
    })
  }
})

// 监听命名空间变化
watch(() => props.namespace, (namespace) => {
  form.metadata.namespace = namespace
})
</script>

<style scoped>
.container-item {
  margin-bottom: 10px;
  padding: 10px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
}

.port-section,
.env-section,
.resource-section {
  margin-bottom: 20px;
  padding: 15px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
}

.port-section h4,
.env-section h4,
.resource-section h4 {
  margin: 0 0 10px 0;
  color: var(--el-text-color-primary);
  font-size: 14px;
}

.port-item,
.env-item {
  margin-bottom: 10px;
  padding: 10px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.strategy-config {
  margin-top: 10px;
  padding: 10px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

h5 {
  margin: 0 0 8px 0;
  color: var(--el-text-color-regular);
  font-size: 12px;
}
</style> 