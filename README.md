# Golang

1. History
    * Created by Robert Griesemer, Rob Pike, and Ken Thompson
    * "Go was designed at Google in 2007 to improve programming productivity in an era of multicore, networked machines and large codebases." - wikipedia
    * Initially announced in 2009, and released in 2012, the current stable version is 1.13.4
    * Compiler can be found at golang.org, or any linux package repository

2. Paradigm
    * The language has functional features, but appears to be largely imperative/OO
3. Typing System
    * Strongly typed, declarations not required
4. Control Structures
    * for loops (also function as while loops)
    * if statements
        * can have pre-conditions
        * eg: ``` if x := true; !x {...} ``` makes x true before testing on it
        * no ternary operator :(
    * Switch statements
        * no need for ```break``` like in java, only one case is run
        * cases don't need to be constants
    * Defer
        * defers the execution of a function until the surrounding function returns
5. Semantics
6. Desirable Language Characteristics
    * Security - Go has garbage collection
