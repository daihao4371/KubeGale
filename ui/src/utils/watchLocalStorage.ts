export const dispatchClusterIDEventStrage = (cluster_id: string) => {
    // 创建一个StorageEvent事件
    const newStorageEvent = document.createEvent('StorageEvent')
    const storage = {
        setItem: function(cluster_id: string) {
            const oldValue = localStorage.getItem('cluster_id') || '';
            localStorage.setItem('cluster_id', cluster_id)
            // 初始化创建的事件
            newStorageEvent.initStorageEvent('setItemCluster', false, false, 'cluster_id', oldValue, cluster_id, window.location.href, window.localStorage)
            // 派发对象
            window.dispatchEvent(newStorageEvent)
        }
    }
    return storage.setItem(cluster_id)
}

export const dispatchNamespaceEventStrage = (namespace: string) => {
    // 创建一个StorageEvent事件
    const newStorageEvent = document.createEvent('StorageEvent')
    const storage = {
        setItem: function(namespace: string) {
            const oldValue = localStorage.getItem('namespace') || '';
            localStorage.setItem('namespace', namespace)
            // 初始化创建的事件
            newStorageEvent.initStorageEvent('setItemNamespace', false, false, 'namespace', oldValue, namespace, window.location.href, window.localStorage)
            // 派发对象
            window.dispatchEvent(newStorageEvent)
        }
    }
    return storage.setItem(namespace)
}
