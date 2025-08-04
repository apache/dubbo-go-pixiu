/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import Cookies from 'js-cookie'

const TokenKey = 'ticket'

export function getToken() {
  return Cookies.get(TokenKey)
}

export function setToken(token) {
  return Cookies.set(TokenKey, token)
}

export function removeToken() {
  return Cookies.remove(TokenKey)
}

export function getLocalStorage (name) {
  
  var val = localStorage.getItem(name);
  return val ? JSON.parse(localStorage.getItem(name)) : ''
}

export function setLocalStorage (name,val){
  localStorage.setItem(name,JSON.stringify(val))
}

export function addLocalStorage (name,addVal){
  
  let oldVal = getLocalStorage(name) ? getLocalStorage(name) : [];
  let newVal = oldVal.concat(addVal);

  setLocalStorage(name,newVal);
}
export function deleteLocalStorage (name, deleteVal){
  let oldVal = getLocalStorage(name) ? getLocalStorage(name) : [];
  let _index = oldVal.findIndex((val) => val.time == deleteVal);

  if(_index > -1){
    oldVal.splice(_index,1);
  }
  setLocalStorage(name, oldVal);
}
export function clearLocalStorage(name){
  
  localStorage.removeItem(name);
}


//深拷贝
export const deepcopy = function (source) {
  if (!source) {
    return source;
  }
  let sourceCopy = source instanceof Array ? [] : {};
  for (let item in source) {
    sourceCopy[item] = typeof source[item] === 'object' ? deepcopy(source[item]) : source[item];
  }
  return sourceCopy;
};