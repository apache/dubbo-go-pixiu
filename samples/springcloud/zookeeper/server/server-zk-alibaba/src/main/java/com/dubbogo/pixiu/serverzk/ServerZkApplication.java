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
package com.dubbogo.pixiu.serverzk;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cloud.client.discovery.DiscoveryClient;
import org.springframework.cloud.client.serviceregistry.Registration;
import org.springframework.web.bind.annotation.*;

import java.util.Objects;

@RestController
@SpringBootApplication
public class ServerZkApplication {

    @Value("${spring.application.name:pixiu-springcloud-server}")
    private String appName;

    @Autowired
    private DiscoveryClient discovery;

    @Autowired(required = false)
    private Registration registration;

    public static void main(String[] args) {
        SpringApplication.run(ServerZkApplication.class, args);
    }

    @GetMapping({"/hello","/hi"})
    public String hello() {
        return "Hello Pixiu World! from " + this.registration;
    }

    @GetMapping({"/service/{sid}", "/service"})
    public Object service(@PathVariable(value = "sid", required = false) @RequestParam(value = "sid", required = false) String sid) {
        return this.discovery.getInstances(Objects.isNull(sid)||sid.length()==0?appName:sid);
    }

    @GetMapping("/services")
    public Object services() {
        return discovery.getServices();
    }
}
