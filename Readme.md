# CLI Expense Tracker 💸

A simple command-line tool to track your expenses directly from the terminal.

---

## 🚀 Features

- Add new expenses with description, amount, and date
- List recent expenses
- Simple and lightweight — no external database required

---

## 📦 Installation

> Go 1.21+ is required

```bash
git clone https://github.com/Arsiievych/cli-expense-tracker.git
cd cli-expense-tracker
go build -o expense-tracker
```

## ⚙️ Usage
```bash
./expense-tracker [command] [flags]
```
## Available Commands

| Command  | Description               | Flags                           | 
|----------|---------------------------|---------------------------------|
| `add`    | Add a new expense         | --amount / -a <br/> --desc / -d |
| `list`   | Show all expenses         |                                 |
| `remove` | Remove item from the list | --id                            |
| `help`   | Show help for any command |                                 |


### Add a new expense
```bash 
./expense-tracker add --desc "Coffee" --amount 3.50
```

### Show list
```bash 
./expense-tracker list
```
![listCmdExample](./assets/listCmdExample.png)

### Remove from the list
```bash 
./expense-tracker remove --id exp-1749806964834981000
```
