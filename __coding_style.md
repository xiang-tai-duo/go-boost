1. All functions must have only one return statement. The return value variable must be named result and declared using result := false (or the corresponding type initialization). The return value of error type must be named err and declared using err := error(nil). If there are other return values, they must also be declared in this format.
2. The if statement for error checking must be written in the style of if err := func(...); err == nil {} else {}.
3. Abbreviations for all variables are prohibited; full words must be used instead.
4. All variables must be declared using assignment syntax, for example: result := false.
5. Delete all comments except for the ones that begin with //goland
6. All global structs, global constants, and global variables must be declared right after the imports in the order of ```type```,```const```,```var```, and defined in the style of ```type(...)```,```const(...)```, ```var(...)```. Multiple declarations of ```type```,```const```,```var``` must not appear.
7. All struct and constant names must be in full uppercase with words separated by underscores, and their internal declarations must be sorted in ascending alphabetical order by the first letter.
8. For all functions, the main function must come first, followed by functions prefixed with New. If there are multiple New-prefixed functions, they must be sorted in ascending alphabetical order by the part after New. Other public functions must be sorted in ascending alphabetical order and placed after the New-prefixed functions, then private functions sorted in ascending alphabetical order after the public functions. Public methods of structs must be sorted in ascending alphabetical order and placed after private functions, and private methods of structs must be sorted in ascending alphabetical order and placed after the public methods of structs.
9. Write comments at the beginning of the file, following the format below.
    ```
    // Package network
    // File:        network.go
    // Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/network/network.go
    // Author:      Vibe Coding
    // Created:     2025/12/20 12:31:58
    // Description: NETWORK provides functions to get network IP addresses with the smallest metric.
    // --------------------------------------------------------------------------------
    ```
   Creation time must be fixed ```// Created: 2025/12/20 12:31:58```