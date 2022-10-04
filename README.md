# yae

Yet Another Golang Expression Engine

## Motivation

Type safe expression for terminated computing

## Overview

`lex | parse | desugar | typecheck | compile | eval`

- regex based lexer
- tdop parser
- type checking and inferring by abstract interpretation and unify
- compile to closure

## Syntax

- lexicon: lexer/lexicon.go
- reserved: lexer/reserved.go
- operator precedence and associativity: token/type.go
- syntax: parser/factory.go


## Type And Value

`statically strongly type`

type-mapping

| yae      | golang (Val)    |
|----------|-----------------|
| bool     | bool            |
| str      | string          |
| num      | float64         |
| time     | time.Time       |
| list[T]  | []*val          |
| map[K V] | map[string]*val |
| object   | map[string]*val |

`type key = bool | str | num | time`

`type val = bool | str | num | time | list[any] | map[key any] | object`

## Overload And Rank-1 Parametric polymorphism

## Execution

No local scope, only global scope denoted by execution environment.

Name binding and resolution are implemented through the global 
environment(divided into compile-time environment and runtime environment).

## Functions

All callable is function, `obj.method(args ...)` will be desugar to `method(obj, args ...)`.

So, if you want to implement `3.repeat("a") == "aaa"`, you can register func `repeat :: num -> str -> str`

## BuildIn-Funs

```
+ :: num -> num
+ :: num -> num -> num
+ :: str -> str -> str
- :: num -> num
- :: num -> num -> num
- :: time -> time -> num
* :: num -> num -> num
/ :: num -> num -> num
% :: num -> num -> num
^ :: num -> num -> num

== :: bool -> bool -> bool
== :: num -> num -> bool
== :: str -> str -> bool
== :: time -> time -> bool
== :: forall a. (list[a] -> list[a] -> bool)
== :: forall k v. (map[k,v] -> map[k,v] -> bool)

!= :: bool -> bool -> bool
!= :: num -> num -> bool
!= :: str -> str -> bool
!= :: time -> time -> bool
!= :: forall a. (list[a] -> list[a] -> bool)
!= :: forall k v .(map[k,v] -> map[k,v] -> bool)

< :: num -> num -> bool
< :: time -> time -> bool
<= :: num -> num -> bool
<= :: time -> time -> bool
> :: num -> num -> bool
> :: time -> time -> bool
>= :: num -> num -> bool
>= :: time -> time -> bool

abs :: num -> num
round :: num -> num
ceil :: num -> num
floor :: num -> num

max :: num -> num -> num
max :: list[num] -> num
min :: num -> num -> num
min :: list[num] -> num

len :: forall a. (list[a] -> num)
len :: forall k v. (map[k, v] -> num)
len :: str -> num

if :: forall a. (bool -> α -> α -> α)
and :: bool -> bool -> bool
or :: bool -> bool -> bool
not :: bool -> bool

match :: str -> str -> bool
string :: forall a. (a -> str)

isset :: forall k v. (map[k, v] -> k -> bool)

strtotime :: str -> time

print :: forall a. (a -> a)
```