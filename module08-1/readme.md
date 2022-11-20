1. 创建配置
   
   ```bash
   kubectl create configmap http-config --from-env-file=http-config.properties
   ```
2. 创建deployment
   
   ```bash
   kubectl create -f deployment.yaml
   ```