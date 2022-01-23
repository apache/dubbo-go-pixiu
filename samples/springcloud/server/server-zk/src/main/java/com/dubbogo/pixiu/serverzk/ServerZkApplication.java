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
