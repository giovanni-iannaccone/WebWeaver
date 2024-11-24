![Logo](https://github.com/user-attachments/assets/195bc208-0297-4128-83ed-d2a48b25cd83)

# 🕸 WebWeaver

WebWeaver is a load balancer written in Go, designed to handle load balancing efficiently and at scale. With a modular architecture and simple configuration, WebWeaver is ideal for modern deployments and high-availability environments.

## 📚 Features

- **Load Balancing**: Distributes requests across backend servers using various strategies, including Round-Robin, Random and IP Hash.
- **Dynamic Configuration**: Supports real-time configuration changes without needing a restart.
- **Monitoring and Logging**: Provides detailed statistics and request logging for in-depth monitoring.
- **Automatic Failover**: Manages automatic failover of unavailable backend servers.
- **Security**: Includes configuration options for protection against common threats and connection management.


## 👨‍💻 Installation

### 📜  Prerequisites

- **Go 1.25+**: WebWeaver is written in Go and requires a compatible version of Go for compilation.

### 🧪 Installation Steps
1. **Clone the Repository**

    ```sh
    git clone https://github.com/giovanni-iannaccone/WebWeaver
    cd WebWeaver
    ```

2. **Run the Project**

    Ensure you have Go installed, then run:

    ```sh
    go run ./cmd/main.go
    ```

    or compile it with 
    ```sh
    go build ./cmd
    ```

## ⚙ Configuration

WebWeaver configuration is managed through a JSON file. Here's an example configuration:

```json
{
    "algorithm": "rnd",
    "host": "localhost:8080",
    "servers": [
        "localhost:80",
        "localhost:8081"
    ],
    
    "healthCheck": 10,
    "logs": "%USERPROFILE%\\Desktop\\logs.txt",

    "prohibited": [
        "/.env",
        "/secret/"
    ]
}
```

- algorithm: rr for Round Robin, rnd for random choice, iph for ip hash
- host: the main server address
- servers: write here your servers addresses and ports
- healthCheck:  seconds of the healthCheck timeout, put less than or 0 if you don't want the server to do any
- logs: file where to save logs, put nothing between quotes if you don't want to save logs 
- prohibited: file you don't want the server to show


## 🎮 Usage

1. Write your configurations in the configs/configs.json file or give it as an argument with ```--config``` or  ```-c``` flag
2. Run the main file with go


## 🧩 Contributing
We welcome contributions! Please follow these steps:

1. Fork the repository.
2. Create a new branch ( using <a href="https://medium.com/@abhay.pixolo/naming-conventions-for-git-branches-a-cheatsheet-8549feca2534">this</a> convention).
3. Make your changes and commit them with descriptive messages.
4. Push your changes to your fork.
5. Create a pull request to the main repository.


## 🔭 Learn
Golang: https://go.dev/doc/ <br>
Load Balancing: https://www.cloudflare.com/learning/performance/what-is-load-balancing/


## ⚖ License
This project is licensed under the GPL-3.0 License. See the LICENSE file for details.


## ⚔ Contact
- For any inquiries or support, please contact <a href="mailto:iannacconegiovanni444@gmail.com"> iannacconegiovanni444@gmail.com </a>.
- Visit my site for more informations about me and my work <a href="https://giovanni-iannaccone.gith
ub.io" target=”_blank” rel="noopener noreferrer"> https://giovanni-iannaccone.github.io </a>
