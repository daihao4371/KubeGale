<template>
  <div class="menu-item-wrapper">
    <div
      class="menu-item"
      :class="{
        active: activeMenu === item.id,
        'parent-active': hasActiveChild(item)
      }"
      :style="{ paddingLeft: `${16 * (level - 1)}px` }"
      @click="handleClick"
    >
      <el-icon class="menu-item-icon">
        <component :is="item.icon"></component>
      </el-icon>
      <span class="menu-title">{{ item.title }}</span>
      <el-icon
        v-if="item.children && item.children.length"
        class="menu-arrow"
        @click.stop="toggleExpand"
      >
        <arrow-down v-if="isExpanded" />
        <arrow-right v-else />
      </el-icon>
    </div>
    <div
      v-if="item.children && item.children.length && isExpanded"
      class="submenu"
    >
      <MenuItem
        v-for="child in item.children"
        :key="child.id"
        :item="child"
        :activeMenu="activeMenu"
        :expandedMenus="expandedMenus"
        :level="level + 1"
        @selectMenu="$emit('selectMenu', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, defineProps, defineEmits, watch } from 'vue'
import { ArrowDown, ArrowRight } from '@element-plus/icons-vue'
const props = defineProps({
  item: {
    type: Object,
    required: true
  },
  activeMenu: String,
  expandedMenus: Array,
  level: {
    type: Number,
    default: 1
  }
})
const emit = defineEmits(['selectMenu'])
const hasActiveChild = (item: any) => {
  if (!item.children) return false
  return item.children.some((child: any) => child.id === props.activeMenu || hasActiveChild(child))
}
const handleClick = () => {
  if (!props.item.children || !props.item.children.length) {
    emit('selectMenu', props.item.id)
  }
}
const isExpanded = ref(false)
const toggleExpand = () => {
  isExpanded.value = !isExpanded.value
}
// 自动展开包含当前激活菜单的父级
watch(
  () => props.activeMenu,
  (val) => {
    if (hasActiveChild(props.item)) {
      isExpanded.value = true
    }
  },
  { immediate: true }
)
</script>

<style scoped>
.menu-item-wrapper {
  position: relative;
}
</style> 