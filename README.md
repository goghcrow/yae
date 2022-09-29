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

## BuildIn

```
+ :: num -> num -> num 
- :: num -> num -> num
* :: num -> num -> num
/ :: num -> num -> num
% :: num -> num -> num

+  :: str -> str -> str

>  :: num -> num -> bool
>= :: num -> num -> bool
<  :: num -> num -> bool
<= :: num -> num -> bool
== :: num -> num -> bool
!= :: num -> num -> bool

>  :: time -> time -> bool
>= :: time -> time -> bool
<  :: time -> time -> bool
<= :: time -> time -> bool
== :: time -> time -> bool
!= :: time -> time -> bool

todo...

```