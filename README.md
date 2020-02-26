# SOFAMesh

**This project is deprecated. It will contribute to istio directly instead of developing in a forked repo.**

---

SOFAMesh is a large-scale implementation scheme of Service Mesh based on [Istio](https://istio.io). On the basis of inheriting the powerful functions and rich features of Istio, in order to meet the performance requirements in large-scale deployments and to respond to the actual situation in the implementation, the following improvements are made:

- [MOSN](https://github.com/mosn/mosn) written in Golang instead of [Envoy](https://github.com/envoyproxy/envoy)
- Merge Mixer to data plane to resolve performance bottlenecks
- Enhance Pilot for more flexible service discovery mechanism
- Added support for [SOFA RPC](https://github.com/sofastack/sofa-rpc), Dubbo

The initial version was contributed by Ant Financial and Alibaba UC Business Unit.

---

**该项目已弃用。该项目将直接向 Istio 贡献，而不是在 fork 的仓库中开发。**

SOFAMesh 是基于 [Istio](https://istio.io) 改进和扩展而来的 Service Mesh 大规模落地实践方案。在继承 Istio 强大功能和丰富特性的基础上，为满足大规模部署下的性能要求以及应对落地实践中的实际情况，有如下改进：

- 采用 Golang 编写的 [MOSN](https://github.com/mosn/mosn) 取代 [Envoy](https://github.com/envoyproxy/envoy)
- 合并 Mixer 到数据平面以解决性能瓶颈
- 增强 Pilot 以实现更灵活的服务发现机制
- 增加对 [SOFA RPC](https://github.com/sofastack/sofa-rpc)、Dubbo 的支持

初始版本由蚂蚁金服和阿里大文娱UC事业部携手贡献。

