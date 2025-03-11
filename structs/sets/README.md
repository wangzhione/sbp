# struct go map set

[glang-set](https://github.com/deckarep/glang-set) deckarep 项目的精简版本 set simple. Thank you ✍

## Install

Use `go get` to install this package.

```shell
go get github.com/wangzhione/sbp@latest
```

## Usage

```go
// Syntax example, doesn't compile.
mySet := sets.NewSet[T]() // where T is some concrete comparable type.

// Therefore this code creates an int set
mySet := sets.NewSet[int]()

// Or perhaps you want a string set
mySet := sets.NewSet[string]()

type myStruct struct {
  name string
  age uint8
}

// Alternatively a set of structs
mySet := sets.NewSet[myStruct]()

// Lastly a set that can hold anything using the any or empty interface keyword: interface{}. This is effectively removes type safety.
mySet := sets.NewSet[any]()
```

### Comprehensive Example

```go
package main

import (
  "fmt"

  "github.com/wangzhione/sbp/structs/sets"
)

func main() {
  // Create a string-based set of required classes.
  required := sets.NewSet[string]()
  required.Add("cooking")
  required.Add("english")
  required.Add("math")
  required.Add("biology")

  // Create a string-based set of science classes.
  sciences := sets.NewSet[string]()
  sciences.Add("biology")
  sciences.Add("chemistry")
  
  // Create a string-based set of electives.
  electives := sets.NewSet[string]()
  electives.Add("welding")
  electives.Add("music")
  electives.Add("automotive")

  // Create a string-based set of bonus programming classes.
  bonus := sets.NewSet[string]()
  bonus.Add("beginner go")
  bonus.Add("python for dummies")
}
```

