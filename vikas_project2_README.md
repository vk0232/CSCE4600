# Project Files CSCE 4600



## [Project 2: Shell Builtins](https://github.com/vk0232/CSCE4600/tree/main/Project2)

A twist on a classic "build your own shell". The *very* basic shell is already written, but you will choose five (5) shell builtins (or shell-adjacent) commands to rewrite into Go, and integrate into the Go shell.

Main.go

made changes to add the builtins - to add the mkdir, rmdir, pwd, time and mv shell builtin commands
func handleInput(w io.Writer, input string, exit chan<- struct{}) error {
        // Remove trailing spaces.
        input = strings.TrimSpace(input)
        //fmt.Printf("Entered handleInput\n")

        // Split the input separate the command name and the command arguments.
        args := strings.Split(input, " ")
        name, args := args[0], args[1:]

        // Check for built-in commands.
        // New builtin commands should be added here. Eventually this should be refactored to its own func.
        switch name {
        case "cd":
                return builtins.ChangeDirectory(args...)
        case "env":
                return builtins.EnvironmentVariables(w, args...)
        case "mkdir":
                return builtins.MakeDirectory(args...)
        case "rmdir":
                return builtins.RemoveDirectory(args...)
        case "pwd":
                 builtins.GetworkDirectory()
        case "time":
                 builtins.Gettime()
        case "mv":
                return builtins.Renamefile(args...)
        case "exit":
                exit <- struct{}{}
                return nil
        }

        return executeCommand(name, args...)
}

Compiled the code 

vijaykarthikeyaraja@Vijays-MBP CSCE4600 % go build
vijaykarthikeyaraja@Vijays-MBP CSCE4600 % ls -ltr
total 4608
-rw-r--r--@ 1 vijaykarthikeyaraja  staff     2564 Dec  3 17:58 go.sum
drwxr-xr-x@ 7 vijaykarthikeyaraja  staff      224 Dec  3 17:58 Project1
-rw-r--r--@ 1 vijaykarthikeyaraja  staff    35149 Dec  3 17:58 LICENSE
-rw-r--r--@ 1 vijaykarthikeyaraja  staff     1579 Dec  3 18:03 README.md
-rw-r--r--@ 1 vijaykarthikeyaraja  staff      375 Dec  5 12:56 go.mod
-rw-r--r--@ 1 vijaykarthikeyaraja  staff     1029 Dec  5 12:57 main_test.go
drwxr-xr-x@ 7 vijaykarthikeyaraja  staff      224 Dec  8 16:12 builtins
drwxr-xr-x@ 7 vijaykarthikeyaraja  staff      224 Dec  8 17:12 Project2
-rw-r--r--  1 vijaykarthikeyaraja  staff        0 Dec  8 18:28 raja
-rw-r--r--@ 1 vijaykarthikeyaraja  staff     2514 Dec  9 16:15 main.go
-rwxr-xr-x  1 vijaykarthikeyaraja  staff  2298146 Dec  9 16:21 CSCE4600
vijaykarthikeyaraja@Vijays-MBP CSCE4600 % 



