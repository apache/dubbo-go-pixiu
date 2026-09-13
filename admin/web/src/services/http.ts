import {API_BASE,API_STATUS} from '../app/constants'
import {clearSession,readSession,redirectToLogin} from './session'

export class ApiError extends Error{code:string; raw:string; constructor(message:string,code='',raw=''){super(message);this.code=code;this.raw=raw}}

function sessionHeaders(){
  const session=readSession();
  if(!session)return {};
  return {token:session.token,username:session.username};
}

function isAuthFailure(status:number,bodyData:unknown){
  const body=bodyData&&typeof bodyData==='object'?bodyData as {code?:unknown;data?:unknown}:{};
  if(status===401||status===403)return true;
  if(body.code!==API_STATUS.NOT_FOUND)return false;
  const message=typeof body.data==='string'?body.data.toLowerCase():'';
  return /token|login|access|authoriz|认证|登录|权限/.test(message);
}

export function isNotFoundError(error:unknown){return error instanceof ApiError&&error.code===API_STATUS.NOT_FOUND}

export function parseArrayResponse<T>(data:unknown):T[]{
  let parsed=data;
  if(typeof data==='string'){
    try{parsed=JSON.parse(data) as unknown}catch(error){throw new ApiError('接口返回的数据不是有效 JSON','BAD_RESPONSE',error instanceof Error?error.message:String(error))}
  }
  if(!Array.isArray(parsed))throw new ApiError('接口返回的数据格式无效','BAD_RESPONSE',String(data));
  return parsed as T[];
}

export async function request<T>(path:string,init:RequestInit={}){
  const body=init.body;
  const headers=new Headers(init.headers);
  Object.entries(sessionHeaders()).forEach(([key,value])=>headers.set(key,value));
  if(body && !(body instanceof URLSearchParams) && !(body instanceof FormData) && !headers.has('Content-Type'))headers.set('Content-Type','application/json');
  const r=await fetch(API_BASE+path,{...init,headers});
  const raw=await r.text();
  if(r.status===401||r.status===403){clearSession();redirectToLogin();throw new ApiError('登录状态已失效，请重新登录','AUTH_REQUIRED',raw)}
  if(r.status===204)return undefined as T;
  let bodyData:unknown;
  try{bodyData=JSON.parse(raw)}catch{
    const boundary=raw.indexOf('}{');
    if(boundary>0){try{bodyData=JSON.parse(raw.slice(0,boundary+1))}catch{bodyData=null}}
  }
  if(bodyData===null||typeof bodyData!=='object')throw new ApiError(`接口返回了非 JSON 内容（HTTP ${r.status}）`,'BAD_RESPONSE',raw);
  const envelope=bodyData as {code?:unknown;data?:unknown;message?:unknown};
  if(!r.ok){
    if(isAuthFailure(r.status,envelope)){clearSession();redirectToLogin();throw new ApiError('登录状态已失效，请重新登录','AUTH_REQUIRED',String(envelope.data||raw))}
    const message=typeof envelope.data==='string'?envelope.data:typeof envelope.message==='string'?envelope.message:`请求失败（HTTP ${r.status}）`;
    throw new ApiError(message,String(envelope.code||`HTTP_${r.status}`),raw);
  }
  if(envelope.code&&envelope.code!==API_STATUS.SUCCESS){
    if(isAuthFailure(r.status,envelope)){clearSession();redirectToLogin();throw new ApiError('登录状态已失效，请重新登录','AUTH_REQUIRED',String(envelope.data||''))}
    const message=typeof envelope.data==='string'?envelope.data:envelope.code===API_STATUS.NOT_FOUND?'数据不存在':envelope.code===API_STATUS.CONCURRENT?'操作冲突，请刷新后重试':'请求失败';
    throw new ApiError(message,String(envelope.code),String(envelope.data||''));
  }
  if(!('data' in envelope))throw new ApiError('接口响应缺少 data 字段','BAD_RESPONSE',raw);
  return envelope.data as T
}
