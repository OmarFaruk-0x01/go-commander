# CLI Task Manager with GO

![Interface](/demo/demo.gif)

[Demo Video](/demo/demo.webm)

The Commander is a flexible and powerful cli tool for running multiple commands simultaneously, with additional features like graceful task termination, keypress interactions, log management, and restart tasks. It provides a simple and intuitive UI to manage concurrent tasks in your Terminal.

## Why Use the Commander?

- **Efficient Concurrent Execution**: Run multiple commands concurrently, allowing for efficient resource utilization and improved performance.

- **Graceful Task Termination**: Stop running tasks gracefully by sending interrupt signals, enabling cleanup operations before termination.

- **Restart Running Tasks**: Restart running tasks gracefully by stopping and restarting them.

- **Keypress Interactions**: Listen for keypress events to perform actions like restarting or killing specific tasks, providing interactive control over the concurrent execution.

- **Log Management**: Capture and display logs for each running task, enabling real-time monitoring.

## Features

- Run multiple commands simultaneously with Goroutines.
- Gracefully stop tasks by sending interrupt signals.
- Restart tasks by stopping and starting new ones with a different command.
- Listen for keypress events to perform actions on running tasks.
- Capture and display logs for each task.
- Interactive prompt for restarting or killing tasks.
- Colorful indication for each tasks.

## Installation

To use the Commander in your system, you need to have Go installed. Then, you can run the following command:

```shell
# Linux
sudo curl -L https://github.com/OmarFaruk-0x01/go-commander/releases/download/v1.0.1/commander-linux-amd64 --output /usr/bin/commander

# Mac
sudo curl -L https://github.com/OmarFaruk-0x01/go-commander/releases/download/v1.0.1/commander-darwin-amd64 --output /usr/local/bin/commander
```

## Usage

Here's a basic example demonstrating the usage of the Commander:

```shell
commander -cmd "<command 1>" -cmd "<command 2>" -cmd "<command 3>" ...
```

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request. Make sure to follow the project's code style and guidelines.

## License

This project is licensed under the MIT License.
