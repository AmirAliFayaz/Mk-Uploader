# 🚀 Mk-Uploader

**Mk-Uploader** is a simple and fast file uploader written in **Go**. Its straightforward usage lets you upload files with just a single command in the terminal.

### 📜 Features
- **Fast** and lightweight
- **Simple CLI Usage**
- Developed with **Go**, perfect for performance-focused applications

---

## 🛠️ Installation & Setup

### Step 1: Download and Install Go
Before using Mk-Uploader, ensure **Go** is installed on your system. You can download it from the official website:

➡️ [Download Go](https://golang.org/dl/)

### Step 2: Clone and Compile Mk-Uploader
1. Clone the repository:
   ```bash
   git clone https://github.com/AmirAliFayaz/Mk-Uploader.git
   cd Mk-Uploader
### Compile the Go file:
```bash
go mod init upload.go
go mod tidy
go build -o upload upload.go
```
This will create an executable named **`upload`** in your project directory.

### 📦 Usage
To use Mk-Uploader, simply run the following command:
```bash
upload.exe <filename>
```
#### Example:
```bash
upload.exe myfile.txt
```
#### 💡 Tips
-   Ensure that the file you want to upload exists in the same directory or provide the full path.
-   **Mk-Uploader**

### 📝 License
This project is open-source and available under the **MIT License**.

##### Enjoy using **Mk-Uploader** for your file uploads! 🎉
