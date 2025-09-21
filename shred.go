package main 

import (
    "os"
    "log"
    "crypto/rand"
)

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

    Return value: 
        - error: Returns error if error 
*/
func Shred(filePath string, overwriteCount int) error{
    // Check if file exists 
    _, err := os.Stat(filePath)
    if os.IsNotExist(err){
        log.Printf("Filepath %s does not exist", filePath)
        return err 
    }

    // Get size of file 
    fileStats, err := os.Stat(filePath)
    if err != nil{
        log.Print("Error getting stats from file: ", err)
        return err 
    }
    fileSize := fileStats.Size()
    log.Println("Size of file: ", fileSize)

    // Store read/write permissions 
    permissions := os.FileMode(0644)

    // Assign permissions to file 
    err = os.Chmod(filePath, 0644)
    if err != nil {
        log.Print("Error setting file permissions: ", err)
        return err
    }
    // Overwrite file specified # of times 
    for i:=0; i<overwriteCount; i++{
        randomBytes, err := GenerateRandomBytes(fileSize)
        if err != nil{
            log.Print("Error generating random bytes: ", err)
            return err 
        }
        err = os.WriteFile(filePath, randomBytes, permissions)
        if err != nil {
            log.Print("Error writing random bytes to file: ", err)
            return err 
        }
        log.Println("Shredding file. Pass #", i+1, " complete.")
    }

    log.Println("Deleting file...")
    os.Remove(filePath)
    if err != nil {
        log.Print("Error deleting file: ", err)
        return err 
    }
    return nil 
}