import {API_BASE,API_STATUS} from '../app/constants'
import {clearSession,readSession,redirectToLogin} from './session'

export class ApiError extends Error{code:string; raw:string; constructor(message:string,code='',raw=''){super(message);this.code=code;this.raw=raw}}

function sessionHeaders(){
  const session=readSession();
  if(!session)return {};
  return {token:session.token,username:session.username};
}

function isAuthFailure(status:number,bodyData:any){
  if(status===401||status===403)return true;
  if(bodyData?.code!==API_STATUS.NOT_FOUND)return false;
  const message=typeof bodyData.data==='string'?bodyData.data.toLowerCase():'';
  return /token|login|access|authoriz|认证|登录|权限/.test(message);
}

export async function request<T>(path:string,init:RequestInit={}){
  const body=init.body;
  const headers=new Headers(init.headers);
  Object.entries(sessionHeaders()).forEach(([key,value])=>headers.set(key,value));
  if(body && !(body instanceof URLSearchParams) && !(body instanceof FormData) && !headers.has('Content-Type'))headers.set('Content-Type','application/json');
  const r=await fetch(API_BASE+path,{...init,headers});
  const raw=await r.text();
  if(r.status===401||r.status===403){clearSession();redirectToLogin();throw new ApiError('登录状态已失效，请重新登录','AUTH_REQUIRED',raw)}
  let bodyData:any;
  try{bodyData=JSON.parse(raw)}catch{
    const boundary=raw.indexOf('}{');
    if(boundary>0){try{bodyData=JSON.parse(raw.slice(0,boundary+1))}catch{bodyData=null}}
  }
  if(!bodyData)throw new ApiError(`接口返回了非 JSON 内容（HTTP ${r.status}）`,'BAD_RESPONSE',raw);
  if(bodyData.code&&bodyData.code!==API_STATUS.SUCCESS){
    if(isAuthFailure(r.status,bodyData)){clearSession();redirectToLogin();throw new ApiError('登录状态已失效，请重新登录','AUTH_REQUIRED',String(bodyData.data||''))}
    throw new ApiError(bodyData.data|| (bodyData.code===API_STATUS.NOT_FOUND?'数据不存在':bodyData.code===API_STATUS.CONCURRENT?'操作冲突，请刷新后重试':'请求失败'),bodyData.code,String(bodyData.data||''));
  }
  return bodyData.data as T
}
