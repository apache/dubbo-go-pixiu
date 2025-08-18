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
  <el-row :gutter="10" :span="24">
    <template v-if="!detail">
    <el-col :span="span">
      <!-- style="width:calc(100% - 6px);" -->
      <el-select v-model="province" clearable placeholder="请选择" class="w-206">
        <el-option label="全国" value="" v-if="showAll"></el-option>
        <el-option v-for="item in Provinces" :key="item.provCode" :label="item.provName" :value="item.provCode">
        </el-option>
      </el-select>
    </el-col>
    <el-col :span="span">
      <!-- style="width:calc(100% - 6px);" -->
      <el-select v-model="city" clearable placeholder="请选择" class="w-206">
        <el-option v-for="item in Cities" :key="item.cityCode" :label="item.cityName" :value="item.cityCode">
        </el-option>
      </el-select>
    </el-col>
    <el-col :span="span">
      <!-- style="width:calc(100% - 6px);" -->
      <el-select v-model="country" clearable placeholder="请选择" class="w-206">
        <el-option v-for="item in Countries" :key="item.regiCode" :label="item.regiName" :value="item.regiCode">
        </el-option>
      </el-select>
    </el-col>
    </template>
    <span v-else>{{currentProvince}} {{currentCity}}{{currentCountry}}</span>
    <slot></slot>
  </el-row>
</template>

<script>
export default {
  name: 'AddrSelect',
  props: {
    // v-model 绑定
    value: {
      type: String,
      required: true,
      default: ''
    },
    // 传输值的分隔符
    symbol: {
      type: String,
      default: '@'
    },
    // 是否展示“全国”选项，表单搜索情况下需要
    showAll: {
      type: Boolean,
      default: true
    },
    span: {
      type: Number,
      default: 8
    },
    detail: {
      type: Boolean
    }
  },
  data () {
    return {
      setValueing: false, // 是否正在一串异步设置选择值过程中，在此期间，不能更新相关的值，只能清空对应的表单
      province: '',
      city: '',
      country: '',
      Provinces: [],
      Cities: [],
      Countries: []
    }
  },
  computed: {
    addr () {
      return this.value ? this.value.replace(/null/g, '').split(this.symbol) : []
    },
    currentProvince () {
      return this.getAddress(this.Provinces, this.province, 'provCode', 'provName')
    },
    currentCity () {
      return this.getAddress(this.Cities, this.city, 'cityCode', 'cityName')
    },
    currentCountry () {
      return this.getAddress(this.Countries, this.country, 'regiCode', 'regiName')
    }
  },
  watch: {
    province (v) {
      this.Cities = []
      this.Countries = []
      if (this.setValueing) {
        return
      }
      this.city = ''
      this.country = ''
      v && this.getCities(v)
      this.emitValue()
    },
    city (v) {
      this.Countries = []
      if (this.setValueing) {
        return
      }
      this.country = ''
      v && this.getCountries(v)
      this.emitValue()
    },
    country () {
      this.emitValue()
    },
    value (v) {
      console.log(v, 'v')
      if (v && v.indexOf(this.symbol) === -1) {
        this.$message.warning('未识别的地址格式，请重新选择')
        return
      }
      if (!v) {
        this.clear()
      }
      this.setAddr()
    }
  },
  methods: {
    getAddress (array, val, key, target) {
      let result = ''
      for (let index = 0; index < array.length; index++) {
        const element = array[index][key]
        if (element === val) {
          result = array[index][target]
          break
        }
      }
      return result
    },
    emitValue () {
      this.$emit('input', [this.province, this.city, this.country].join(this.symbol))
      this.$emit('selected', {
        province: this.province,
        city: this.city,
        country: this.country
      })
    },
    setAddr () {
      if (this.setValueing) {
        return
      }
      this.setValueing = true
      this.getProvincesIfNeeded()
        .then(() => {
          this.province = this.addr[0] || ''
          if (this.province) {
            return Promise.resolve(true)
          } else {
            return Promise.reject(false)
          }
        })
        .then(this.getCitiesIfNeeded)
        .then(() => {
          this.city = this.addr[1] || ''
          if (this.city) {
            return Promise.resolve(true)
          } else {
            return Promise.reject(false)
          }
        })
        .then(this.getCountriesIfNeeded)
        .then(() => {
          this.country = this.addr[2] || ''
        })
        .catch(err => {
          console.log(err)
        })
        .finally(() => {
          this.setValueing = false
        })
    },
    clear () {
      this.province = ''
      this.city = ''
      this.country = ''
    },
    getProvincesIfNeeded () {
      return new Promise((resolve, reject) => {
        if (this.Provinces.length) {
          resolve(true)
        }
        this.getProvinces().then(res => {
          if (this.Provinces.length) {
            resolve(true)
          } else {
            reject(false)
          }
        }).catch(() => {
          reject(false)
        })
      })
    },
    getCitiesIfNeeded () {
      return new Promise((resolve, reject) => {
        if (this.Cities.length) {
          resolve(true)
        }
        this.getCities(this.province).then(res => {
          if (this.Cities.length) {
            resolve(true)
          } else {
            reject(false)
          }
        }).catch(() => {
          reject(false)
        })
      })
    },
    getCountriesIfNeeded () {
      return new Promise((resolve, reject) => {
        if (this.Countries.length) {
          resolve(true)
        }
        this.getCountries(this.city).then(res => {
          if (this.Countries.length) {
            resolve(true)
          } else {
            reject(false)
          }
        }).catch(() => {
          reject(false)
        })
      })
    },
    getProvinces () {
      return this.$http.get('/basic/province').then(res => {
        res = res.data
        if (res.code === 200) {
          this.Provinces = res.data || []
        } else {
          this.$message.error(res.msg)
        }
      }).catch(err => {
        this.$message.error('网络错误，请重试！')
      })
    },
    getCities (key) {
      return this.$http.get('/basic/city', {
        params: {
          key
        }
      }).then(res => {
        res = res.data
        if (res.code === 200) {
          this.Cities = res.data || []
        } else {
          this.Cities = []
          this.$message.error(res.msg)
        }
      }).catch(err => {
        this.Cities = []
        this.$message.error('网络错误，请重试！')
      })
    },
    getCountries (key) {
      return this.$http.get('/basic/region', {
        params: {
          key
        }
      }).then(res => {
        res = res.data
        if (res.code === 200) {
          this.Countries = res.data || []
        } else {
          this.Countries = []
          this.$message.error(res.msg)
        }
      }).catch(err => {
        this.Countries = []
        this.$message.error('网络错误，请重试！')
      })
    },
    findKeybyName (name, type) {
      if (type === 'province') {
        return this.Provinces.find(p => {
          return p.provName === name
        }).provId
      } else if (type === 'city') {
        return this.Cities.find(p => {
          return p.cityName === name
        }).cityId
      } else {
        return this.Countries.find(p => {
          return p.regiName === name
        }).regiId
      }
    }
  },
  created () {
    this.setAddr()
  },
  mounted () {

  }
}

</script>

<style scoped lang='scss'>
  
</style>
