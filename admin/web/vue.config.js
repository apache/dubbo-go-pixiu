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
const path = require('path');

const BundleAnalyzerPlugin = require('webpack-bundle-analyzer').BundleAnalyzerPlugin

// const ThreeExamples = require('import-three-examples')
function resolve(dir) {
    return path.join(__dirname,'.', dir);
}

module.exports = {
    baseUrl : '/',
    publicPath : '/',
    productionSourceMap: true,
    chainWebpack: config => {
        config.resolve.alias
        .set('@', resolve('src'))
        .set('@assets',resolve('src/assets'))
        .set('styles',resolve('src/styles'))
        .set('@utils',resolve('src/utils'))
        .set('@views',resolve('src/views'))
    },
    devServer:{
        // 设置代理
        // before: require('./src/mock'),
        host: '0.0.0.0',
        port: 8080,
        hot: true,
        https: false,
        open: false,
        disableHostCheck: true,
        proxy: {
            "/config": {
                target: process.env.VUE_APP_BACKEND_URL,
                ws: true, // enable websocket proxy
                changOrigin: true, // enable proxy module
                // proxy url replace rules
                pathRewrite:{
                    '^/config':''
                },
            },
            "/login": {
                target: process.env.VUE_APP_BACKEND_URL,
                ws: true, // enable websocket proxy
                changOrigin: true, // enable proxy module
                // proxy url replace rules
                pathRewrite:{
                    '^/login':''
                },
            },
        }

    },
    configureWebpack: config => {
        if (process.env.NODE_ENV === 'production') {
            return {
                plugins: [new BundleAnalyzerPlugin()]
            }
        }
    },
     // other vue-cli plugin things...
     pluginOptions: {
        // ...
        // ...ThreeExamples
    },
}
/**
 * page generator
 * @param {String} name page name
 * @param {String} title page title
 */
 function PageReset (name, title) {
    // main.js
    this.entry = `src/entry/${name}.js`
    // template source path
    this.template = 'public/index.html'
    // output filename
    this.filename = `${name}.html`
    this.title = title
    this.chunks = ['chunk-vendors', 'chunk-common', name]
  }
