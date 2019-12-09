# Goborg

I created a markov-chain based chatbot that connects to the discord chat 
service. As go is a reasonably popular language, I was able to find pre-existing
code to connect to discord. Making use of that project, I was able to focus on
making the bot, instead of making the connection. Additionally, it should be
easy to connect it to different services if I choose to do so in the future.

What I wrote myself was the Markov chain implementation. This takes strings of
input (either messages from discord, or from console), and breaks them into pairs
of consecutive words, each word forming a pair with both the words directly
before and after it. These pairs are then added to the chain, updating the
weights on any previously existing edges, and adding new ones as necessary.
Occasionally, the bot will generate a reply to a message it receives, by
traversing the chain from a random word in the message.

The chain itself is just an associative array with a mutex. While the data
structure is present in just about every language, go's maps did make it about
as easy as I would expect it to be. More unique to go is it's concurrency, which
was equally easy to make use of. By just adding a single word to a function call,
```go```, I was able to have it fork into a new thread. This was super useful
when I wanted to have the bot periodically save it's brain in the background.
Originally, I had also intended to use it to read in multiple pairs of words to
the chain at once. Unfortunately, since the bot is based on a _single_ markov
chain in memory, the necessary mutexs blocked near to the point of being a single
thread process. That said, it still allows for multiple reads to occur at once, 
and pushing everything into goroutines took care of queueing up all edge changes
for me. Even if the changes to the chain blocked each other, the rest of the 
program was allowed to continue working.

# Golang

1. History
    * Created by Robert Griesemer, Rob Pike, and Ken Thompson
    * "Go was designed at Google in 2007 to improve programming productivity in an era of multicore, networked machines and large codebases." - wikipedia
    * Initially announced in 2009, and released in 2012, the current stable version is 1.13.4
    * Compiler can be found at golang.org, or any linux package repository

2. Paradigm
    * The language has functional features
        * Closures, function passing/returning
    * Still largely imperitave
3. Typing System
    * Strongly typed, declarations not required
    * Programmer is free to create new types
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
    * Statically scoped
    * Garbage is handled automatically by the runtime
6. Desirable Language Characteristics
    * Security - Go has garbage collection
