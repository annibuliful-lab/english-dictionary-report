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
curl -L -o edr https://github.com/annibuliful-lab/english-dictionary-report/releases/latest/download/linux-amd64
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
Invoke-WebRequest -Uri "https://github.com/annibuliful-lab/english-dictionary-report/releases/latest/download/windows-amd64.exe" -OutFile "edr.exe"

# Run the program with arguments
.\edr.exe v1 .\words.txt .\output
```

### Screenshots

<img width="584" alt="Screenshot 2568-05-25 at 20 06 27" src="https://github.com/user-attachments/assets/ecde204e-afe8-4123-8c6c-87c35c7a57e9" />
<img width="642" alt="Screenshot 2568-05-25 at 20 07 09" src="https://github.com/user-attachments/assets/1654e21d-41ab-48e9-9e24-34285831af30" />
<img width="561" alt="Screenshot 2568-05-25 at 20 06 51" src="https://github.com/user-attachments/assets/3a448205-a297-4197-ac00-6e9b6f22801e" />
