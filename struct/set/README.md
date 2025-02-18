# struct go map set

[glang-set](https://github.com/deckarep/glang-set) deckarep 项目的精简版本. Thank you ✍

## Install

Use `go get` to install this package.

```shell
go get github.com/wangzhione/sbp@latest
```

## Usage

```go
// Syntax example, doesn't compile.
mySet := set.NewSet[T]() // where T is some concrete comparable type.

// Therefore this code creates an int set
mySet := set.NewSet[int]()

// Or perhaps you want a string set
mySet := set.NewSet[string]()

type myStruct struct {
  name string
  age uint8
}

// Alternatively a set of structs
mySet := set.NewSet[myStruct]()

// Lastly a set that can hold anything using the any or empty interface keyword: interface{}. This is effectively removes type safety.
mySet := set.NewSet[any]()
```

### Comprehensive Example

```go
package main

import (
  "fmt"
  "github.com/wangzhione/sbp/struct/set"
)

func main() {
  // Create a string-based set of required classes.
  required := set.NewSet[string]()
  required.Add("cooking")
  required.Add("english")
  required.Add("math")
  required.Add("biology")

  // Create a string-based set of science classes.
  sciences := set.NewSet[string]()
  sciences.Add("biology")
  sciences.Add("chemistry")
  
  // Create a string-based set of electives.
  electives := set.NewSet[string]()
  electives.Add("welding")
  electives.Add("music")
  electives.Add("automotive")

  // Create a string-based set of bonus programming classes.
  bonus := set.NewSet[string]()
  bonus.Add("beginner go")
  bonus.Add("python for dummies")
}
```

