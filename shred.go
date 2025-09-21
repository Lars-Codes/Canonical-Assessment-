package main 

import (
    "os"
    "log"
    "crypto/rand"
)

func main(){

    filePath := "test.txt"
    overwriteCount := 3 

    shred(filePath, overwriteCount)
}

/*
    Purpose: Generate a given number of random bytes using secure library crypto/rand. 
    Parameters: 
        - numBytes: The number of bytes to return 
    
    Return types: 
        - []byte: Returns a random byte array of length numBytes
        - err: Returns error if applicable 

*/
func GenerateRandomBytes(numBytes int64)([]byte, error){
    randBytes := make([]byte, numBytes) // Random byte array 

    _, err := rand.Read(randBytes) // Generate random bytes 
    if err != nil {
        return nil, err
    }

    return randBytes, nil 
}


/*
    Purpose: "Shred" a file a certain number of times by overwriting it with random bytes 
    Parameters: 
        - filePath: Path to file to shred 
        - overwriteCount: Number of times to overwrite the file with random bytes 

*/
func shred(filePath string, overwriteCount int){
    // Check if file exists 
    _, err := os.Stat(filePath)
    if os.IsNotExist(err){
        log.Fatalf("Filepath %s does not exist", filePath)
    }

    // Get size of file 
    fileStats, err := os.Stat(filePath)
    if err != nil{
        log.Fatal("Error getting stats from file: ", err, " Exiting application.")
    }
    fileSize := fileStats.Size()
    log.Println("Size of file: ", fileSize)

    // Store read/write permissions 
    permissions := os.FileMode(0644)

    // Overwrite file specified # of times 
    for i:=0; i<overwriteCount; i++{
        randomBytes, err := GenerateRandomBytes(fileSize)
        if err != nil{
            log.Fatal("Error generating random bytes: ", err, " Exiting application.")
        }
        err = os.WriteFile(filePath, randomBytes, permissions)
        if err != nil {
            log.Fatal("Error writing random bytes to file: ", err, " Exiting application.")
        }
        log.Println("Shredding file. Pass #", i+1, " complete.")
    }

    log.Println("Deleting file...")
    os.Remove(filePath)
    if err != nil {
        log.Fatal("Error deleting file: ", err, " Exiting application.")
    }

}