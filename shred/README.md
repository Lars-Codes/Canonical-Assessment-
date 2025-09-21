# Canonical coding assessment (Go)

This shred script will overwrite a given file (eg. 'randomfile') a given number of times with 
random data and then delete the file afterwards. Test coverage is included in this repo in the 
file "shred_test.go."

---

## Table of Contents

- [Installation](#installation)  
- [Usage](#usage)  
- [Tests](#tests)  
- [License](#license)  


## Installation

Clone the repository and build the binary:

```bash
git clone <repository_url>
cd shred
go build -o shred

```

## Usage 
Call the shred function in your code:  
```
package main

func main() {
    filePath := "example.txt"
    overwriteCount := 3

    err := Shred(filePath, overwriteCount)
    if err != nil {
        log.Fatal("Shred failed: ", err)
    }
}

```

## Tests 
Run: 
```
go test -v 
```

## License 
MIT License. [License](https://github.com/Lars-Codes/Canonical-Assessment-/blob/master/LICENSE)
