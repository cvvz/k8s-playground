client-go 源码路径 `vendor/k8s.io/client-go`

代码仓库 https://github.com/kubernetes/client-go

### Client客户端
`rest` 目录提供`RESTClient`客户端，**`RESTClient`是最底层的客户端**

`kubernetes` 目录提供ClientSet客户端。clientSet在RESTClient的基础上封装了资源的管理方法。clientSet是多个客户端的集合（set）。

此外还有DynamicClient客户端和DiscoveryClient客户端


### 版本转换
版本偏差支持策略：https://kubernetes.io/zh/docs/setup/release/version-skew-policy/


#### RESTful
**最常见的一种设计错误，就是URI包含动词。
因为"资源"表示一种实体，所以应该是名词，URI不应该有动词，
动词应该放在HTTP协议中。**

资源注册表


同一种资源在kubernetes中可能有多个版本，每次kubernetes都会帮你存多份。
例如，向kubernetes中存v1版本的sts，它会帮你将资源转换为v1beta1和b1beta2版本，并且保存。
这样，你用低版本的客户端（只支持v1beta1或v1beta2）也能够对sts对象进行操作。

这样，当你用不同版本的客户端去访问kubernetes api-server时也不会有问题。
但是如果服务端已经不支持这个版本了的时候，访问就会有问题。


kube-apiserver会根据你给的manifest中的GVK信息，将资源在etcd中对应的
路径下存储，客户端可以指定任意版本获取资源信息。
可以通过`kubectl api-resources`和`kubectl api-versions`两条命令知道服务端（kubernetes）
的`组/版本/资源`信息。


### Informer机制

https://github.com/kubernetes/sample-controller/blob/master/docs/controller-client-go.md