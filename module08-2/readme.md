1. 创建 ingress nginx controller
   ```bash
   kubectl create -f nginx-ingress-deployment.yaml
   ```

2. 创建证书
   ```bash
   openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=cncamp.com/O=cncamp" -addext "subjectAltName = DNS:cncamp.com"
   ```
   
3. 创建 secret 对象
   ```
   kubectl create secret tls cncamp-tls --cert=./tls.crt --key=./tls.key
   ```

4. 定义 service
   
   ```yaml
    apiVersion: v1
    kind: Service
    metadata:
      name: http-service
    spec:
      type: ClusterIP
      ports:
        - port: 80
          targetPort: 8080
          protocol: TCP
          name: http
      selector:
        app: httpserver
   ```

5. 定义 ingress
   
   ```yaml
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    metadata:
      name: gateway
      annotations:
        kubernetes.io/ingress.class: "nginx"
    spec:
      tls:
        - hosts:
            - cncamp.com
          secretName: cncamp-tls
      rules:
        - host: cncamp.com
          http:
            paths:
              - path: "/"
                pathType: Prefix
                backend:
                  service:
                    name: http-service
                    port:
                      number: 80
   ```

6. 发布 httpserver
   ```
   kubectl apply -f httpserver.yaml
   ```
