# EarleyParser
Simple Earley Parser for a toy grammar, written in Go. 

I learned about Earley Parsers during Formal Models of Languages course at the University of Cambridge.
This is my implementation of the algorithm according to the example in the lectures - it builds and prints out the parse chart for a (hard-coded) sentence. 

The grammar:

N = { S, NP, VP, PP, N, V, P }

Σ = { they, can, fish, in, rivers, December }

S = S

P = { S -> NP VP,

NP -> N PP | N,

PP -> P NP,

VP -> VP PP | V VP | V NP | V,

N  -> can | fish | rivers | December,

P -> in,

V -> can | fish }
