# Depends-on for Kubernetes

![Kubernetes Logo](https://kubernetes.io/images/kubernetes-horizontal-color.png)

## 🚀 Introduction

Depends-on is a powerful tool that brings the familiar `depends_on` logic from Docker Compose to Kubernetes environments. It addresses a common challenge in Kubernetes deployments by providing a simple and effective way to manage initialization order between multiple services.

## 🌟 Key Features

- **Docker Compose Compatibility**: Implement `depends_on`-like functionality in Kubernetes
- **Flexible Service Orchestration**: Easily manage boot order and dependencies between services
- **Simplified Deployment**: Reduce complexity in service initialization and startup processes

## 🤔 Why Depends-on?

While Kubernetes intentionally doesn't natively support the `depends_on` logic (encouraging self-contained, resilient services), there are scenarios where managing service dependencies can significantly simplify deployment and reduce development time.

### Use Cases:

- Legacy applications that aren't designed for Kubernetes-native resilience
- Development and testing environments where quick setups are prioritized
- Complex systems with intricate startup requirements

## ⚠️ Important Considerations

Before using Depends-on, consider the following:

1. In production-grade, high-demand infrastructures, it's generally recommended that each service manages its own boot state natively.
2. Implementing dependency management at the application level can lead to more robust and scalable solutions in the long run.
3. Depends-on is particularly useful for:
   - Rapid prototyping
   - Migrating existing Docker Compose setups to Kubernetes
   - Scenarios where modifying application code is not feasible

## 🛠 Installation

Current version: 0.1

### RBAC

You need to apply the following RBAC to the right namespace

```
kubectl apply -f https://github.com/martient/depends-on/tree/{ VERSION }/config/rbac.yml -n { YOUR-NAMESPACE }
```

## 📘 Usage

```
initContainers:
    - name: depends-on-{the target}
    image: ghcr.io/martient/depends-on:main
    imagePullPolicy: Always
    args: ["--service=first", "--check_interval=10"]
```
## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for more details.

## 📜 License

This project is licensed under the [MIT License](LICENSE).

## 📞 Support

If you encounter any issues or have questions, please [open an issue](https://github.com/martient/depends-on/issues) on our GitHub repository.

---

⭐ If you find this project useful, please consider giving it a star on GitHub! ⭐