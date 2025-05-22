<div align="center">
    <h1>Go To Typescript</h1>
    <h1>Convert Go type to Typescript interface</h1>
</div>

### Usage

Install the Go module in your project:
```
go get github.com/preampinbut/gots@v1.0.0
```

Next run
```
go run github.com/preampinbut/gots --output <output.ts> <input.go>
```

### Example

cmd
```
go run github.com/preampinbut/gots --output output.ts input.go
```

input.go
```go
package mytype

type Name string
type AddressB Address

type Person struct {
	ID       string   `json:"id"`
	Name     Name     `json:"name"`
	Age      int      `json:"age,omitempty"`
	AddressA Address  `json:"address"`
	AddressB AddressB `json:"address2"`
	Friends  []Person `json:"friends"`
}
```

imported.go
```go
package mytype

type UnusedType struct {
	T string `json:"t"`
}

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
}
```

output.ts
```typescript
interface Address {
  street: string;
  city: string;
}

interface Person {
  id: string;
  name: string;
  age?: number;
  address: Address;
  address2: Address;
  friends: Person[];
}
```

### Unhandled Type

```go
type UnhandledType struct {
	Data interface{}      `json:"data"`
	Func func(int) string `json:"func"`
	Chan chan string      `json:"chan"`
	Map  map[string]int   `json:"map"`
}
```

will result in
```typescript
interface UnhandledType {
  data: any;
  func: any;
  chan: any;
  map: any;
}
```
