import {defineConfig,loadEnv} from 'vite'; import react from '@vitejs/plugin-react';

const spaBypass=(req:{method?:string;url?:string})=>req.method==='GET'||req.method==='HEAD'?req.url:undefined;

export default defineConfig(({mode})=>{
  const {VITE_BACKEND_URL='http://127.0.0.1:8081'}=loadEnv(mode,'.','VITE_');
  return {plugins:[react()],server:{host:'0.0.0.0',port:8088,proxy:{'/login':{target:VITE_BACKEND_URL,changeOrigin:true,bypass:spaBypass},'/register':{target:VITE_BACKEND_URL,changeOrigin:true},'/user':VITE_BACKEND_URL,'/config':VITE_BACKEND_URL,'/swagger':VITE_BACKEND_URL}}};
});
