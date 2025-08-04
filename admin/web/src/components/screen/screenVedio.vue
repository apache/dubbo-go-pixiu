<!--
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
 -->
<template>
    <div class="cvedio">
        <div class="cvedio-left ve">
            <div class="txt">主驾侧摄像头</div>
            <div class="player"></div>
        </div>
        <div class="cvedio-right ve">
            <div class="txt">侧驾侧摄像头</div>
            <div class="player"></div>
        </div>
    </div>
</template>
<script>
export default {
    name:'screenVedio',
    data(){
        return{

        }
    }
}
</script>
<style scoped lang="scss">
    .cvedio{
        display: flex;
        padding: 20px 32px;
        height: 340px;
        // line-height: 340px;
        .ve{
            border:1px solid rgb(161, 161, 88);
            flex:1;
            color: #fff;
            display: flex;
            font-size: 18px;
            align-items: center;
            padding: 10px;
            .txt{
                width: 18px;
                margin: 0 12px;
            }
            .player{
                flex: 1;
            }
        }   
        &-left{
            
        }
        &-right{
            
        }
    }
</style>