import {request} from './http';
export const baseApi={get:()=>request<string>('/config/api/base'),save:(content:string)=>request('/config/api/base/',{method:'POST',body:new URLSearchParams({content})})};
