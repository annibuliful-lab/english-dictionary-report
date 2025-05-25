# English Dictionary Report

A high-performance, cross-platform CLI tool for processing English wordlists, built with Go.

## 🚀 Download Binaries

Download the latest binaries from the [Releases page](https://github.com/annibuliful-lab/english-dictionary-report/releases/latest).

| Platform | Architecture | Download Link                                                                                                                                          |
| -------- | ------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Windows  | amd64        | [english-dictionary-report-windows-amd64.exe](https://github.com/annibuliful-lab/english-dictionary-report/releases/latest/download/windows-amd64.exe) |
| Linux    | amd64        | [english-dictionary-report-linux-amd64](https://github.com/annibuliful-lab/english-dictionary-report/releases/latest/download/linux-amd64)             |
| macOS    | amd64        | [english-dictionary-report-darwin-amd64](https://github.com/annibuliful-lab/english-dictionary-report/releases/latest/download/darwin-amd64)           |

> **Note:** Replace `english-dictionary-report` with the actual binary name if it differs.

---

## 🛠 Usage

### ✅ Linux / macOS

```bash
# Download the binary
curl -L -o edr https://github.com/annibuliful-lab/english-dictionary-report/releases/latest/download/english-dictionary-report-linux-amd64
# For macOS, use this instead:
# curl -L -o edr https://github.com/annibuliful-lab/english-dictionary-report/releases/latest/download/english-dictionary-report-darwin-amd64

# Make it executable
chmod +x edr

# Run the program with arguments
./edr v1 ./words.txt ./output
```

### ✅ On Windows

Open PowerShell and run the following commands:

```powershell
# Download the executable
Invoke-WebRequest -Uri "https://github.com/annibuliful-lab/english-dictionary-report/releases/latest/download/english-dictionary-report-windows-amd64.exe" -OutFile "edr.exe"

# Run the program with arguments
.\edr.exe v1 .\words.txt .\output
```
