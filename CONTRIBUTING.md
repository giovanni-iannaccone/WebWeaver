# 🧩 Contributing
We welcome contributions! Please follow these steps:

1. Fork the repository.
2. Create a new branch ( using <a href="https://medium.com/@abhay.pixolo/naming-conventions-for-git-branches-a-cheatsheet-8549feca2534">this</a> convention).
3. Make your changes and commit them with descriptive messages.
4. Push your changes to your fork.
5. Create a pull request to the main repository.

# 👨‍💻 Development guidelines
- **Modularity**: functions should be small and focused on a single responsibility
- **Separation of Concerns**: maintain clear separation between business logic and utility functions
- **Error Handling**: handle errors gracefully. All errors should be logged using the utilities provided
- **Documentation**: add comments to explain complex logic, particularly when working with external libraries

# 🎨 Code style
- **Function naming**: use camelCase for function names
- **Error Handling**: check all errors explicitly, avoid ignoring errors
- **Formatting**: ensure your code is properly formatted using gofmt
- **Imports**: organize imports into three sections:
1. Standard library packages
2. Internal project packages
3. External libraries

Ensure all your tests pass before submitting a PR, update documentation if your changes affect the program's behaviour

Thank you for contributing to WebWeaver! Your efforts are deeply appreciated. If you have any questions, feel free to reach out in the project discussions or open an issue.