---
name: Bug Report
about: Create a report to help us improve
title: '[BUG] '
labels: bug
assignees: ''

---

## 🐛 Bug Report

### Description
A clear and concise description of what the bug is.

### Environment
- **Go Version**: (e.g. go1.24.2)
- **OS**: (e.g. macOS 14.5, Ubuntu 22.04, Windows 11)
- **Architecture**: (e.g. arm64, amd64)
- **Vibe Coding Logger Version**: (e.g. v0.1.0)

### Steps to Reproduce
1. Create a logger with `logger.Default()`
2. Enable system info with `log.EnableSystemInfo(true)`
3. Log a message with `log.Info("test")`
4. See error

### Expected Behavior
A clear and concise description of what you expected to happen.

### Actual Behavior
A clear and concise description of what actually happened.

### Minimal Reproduction Code
```go
package main

import "github.com/your-username/vibe-coding-logger/pkg/logger"

func main() {
    log := logger.Default()
    // Code that reproduces the bug
}
```

### Error Messages
```
Paste any error messages here
```

### Log Output
```
Paste relevant log output here
```

### Additional Context
Add any other context about the problem here.

### System Information
If available, please include the output of:
```go
systemInfo := logger.GetSystemInfo()
fmt.Printf("%+v\n", systemInfo)
```

### Possible Solution
If you have any ideas on how to fix this, please describe them here.

### Checklist
- [ ] I have searched for existing issues
- [ ] I have provided a minimal reproduction case
- [ ] I have included environment information
- [ ] I have included relevant error messages