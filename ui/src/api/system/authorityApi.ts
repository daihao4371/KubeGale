import service from '@/api/request'

export const getApisByAuthority = (authorityId: number) => {
  return service({
    url: '/authorityApi/getApisByAuthority',
    method: 'post',
    data: { authorityId }
  })
} 