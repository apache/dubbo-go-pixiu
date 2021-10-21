

#### readme


there are two spring cloud java project in serve dir which offer same restful api:
- user-service
- auth-service

you can use those to mock below situation
- two service with each one instance
- one service with two instances, just change spring.application.name in  auth-service/src/resource/application.properties from auth-service to user-service notice ! change the port if registry center is nacos
- start a new instance or kill the old to mock instance up or down