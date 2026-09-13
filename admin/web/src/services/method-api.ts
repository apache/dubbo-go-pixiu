import {request} from './http';
export const methodApi={
 list:async(resourceId:string)=>{try{const data=await request<unknown>(`/config/api/resource/method/list?resourceId=${encodeURIComponent(resourceId)}`);try{return typeof data==='string'?JSON.parse(data):data}catch{return []}}catch(e:any){if(e.code==='10002'||String(e.raw||e.message).includes('k/v pair not found'))return [];throw e}},
 detail:(resourceId:string,id:string)=>request<any>(`/config/api/resource/method/detail?resourceId=${encodeURIComponent(resourceId)}&methodId=${encodeURIComponent(id)}`),
 create:(resourceId:string,content:string)=>request(`/config/api/resource/method?resourceId=${encodeURIComponent(resourceId)}`,{method:'POST',body:new URLSearchParams({content})}),
 update:(resourceId:string,id:string,content:string)=>request(`/config/api/resource/method?resourceId=${encodeURIComponent(resourceId)}&methodId=${encodeURIComponent(id)}`,{method:'PUT',body:new URLSearchParams({content})}),
 remove:(resourceId:string,id:string)=>request(`/config/api/resource/method?resourceId=${encodeURIComponent(resourceId)}&methodId=${encodeURIComponent(id)}`,{method:'DELETE'})
}
