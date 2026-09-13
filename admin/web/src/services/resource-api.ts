import {request} from './http';
export const resourceApi={
 list:async()=>{try{const data=await request<unknown>('/config/api/resource/list');try{return typeof data==='string'?JSON.parse(data):data}catch{return []}}catch(e:any){if(e.code==='10002'||String(e.raw||e.message).includes('k/v pair not found'))return [];throw e}},
 detail:(id:string)=>request<any>(`/config/api/resource/detail?resourceId=${encodeURIComponent(id)}`),
 create:(content:string)=>request('/config/api/resource',{method:'POST',body:new URLSearchParams({content})}),
 update:(id:string,content:string)=>request(`/config/api/resource?resourceId=${encodeURIComponent(id)}`,{method:'PUT',body:new URLSearchParams({content})}),
 remove:(id:string)=>request(`/config/api/resource?resourceId=${encodeURIComponent(id)}`,{method:'DELETE'})
}
