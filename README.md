# fillstruct

fillstruct - fills a struct literal with random values

---

For example, given the following types,
```golang
type C struct {
	C1 int64
	C2 *int64
	C3 string
	C4 *string
}
type B struct {
	B1 int64
	B2 string
	B3 float64
	B4 *string
	B5 C
	B6 *C
}
type A struct {
	A1 int64
	A2 string
	A3 float64
	A4 *string
	A5 B
	A6 *B
}

func testA() {
	a := A{
		A1: 1222222,
	}
	_ = a
}

```
the following struct literal
```golang
    a := A{
        A1: 1222222,
    }
```
becomes:
```golang

	a := A{
		A1: 1222222,
		A2: "HT5VwfJhIn",
		A3: 0.1544788879744687,
		A4: thrift.StringPtr("ECnHy1mGBd"),
		A5: B{
			B1: 6031280647731874003,
			B2: "rhk25yNM8e",
			B3: 0.9526626845624445,
			B4: thrift.StringPtr("J9Q8aIYoNp"),
			B5: C{
				C1: 7280439440819657315,
				C2: thrift.Int64Ptr(int64(2786761508419750547)),
				C3: "3xbBltMb9K",
				C4: thrift.StringPtr("2BDzV4A7Kv"),
			},
			B6: &C{
				C1: 2707018739610228200,
				C2: thrift.Int64Ptr(int64(7065193991060922179)),
				C3: "nGLoLmKFwq",
				C4: thrift.StringPtr("uQI56s9ZSf"),
			},
		},
		A6: &B{
			B1: 7604776742081145612,
			B2: "vDIkrP3hQs",
			B3: 0.8395196802725778,
			B4: thrift.StringPtr("5jbJcx4A4m"),
			B5: C{
				C1: 6281496988027306234,
				C2: thrift.Int64Ptr(int64(4015257045867961535)),
				C3: "u2utLX8xFM",
				C4: thrift.StringPtr("X5jwZIUL0A"),
			},
			B6: &C{
				C1: 876092331258217238,
				C2: thrift.Int64Ptr(int64(504834319128121057)),
				C3: "jObWv8JnpJ",
				C4: thrift.StringPtr("Zg3Gs9uEGI"),
			},
		},
	}
	_ = a


```
after applying fillstruct.

```
type C struct {
	C1 int64
	C2 *int64
	C3 string
	C4 *string
}
type B struct {
	B1 int64
	B2 string
	B3 float64
	B4 *string
	B5 C
	B6 *C
}
type A struct {
	A1 int64
	A2 string
	A3 float64
	A4 *string
	A5 B
	A6 *B
}

func testA() {
	a := A{
		A1: 1222222,
		A2: "HT5VwfJhIn",
		A3: 0.1544788879744687,
		A4: thrift.StringPtr("ECnHy1mGBd"),
		A5: B{
			B1: 6031280647731874003,
			B2: "rhk25yNM8e",
			B3: 0.9526626845624445,
			B4: thrift.StringPtr("J9Q8aIYoNp"),
			B5: C{
				C1: 7280439440819657315,
				C2: thrift.Int64Ptr(int64(2786761508419750547)),
				C3: "3xbBltMb9K",
				C4: thrift.StringPtr("2BDzV4A7Kv"),
			},
			B6: &C{
				C1: 2707018739610228200,
				C2: thrift.Int64Ptr(int64(7065193991060922179)),
				C3: "nGLoLmKFwq",
				C4: thrift.StringPtr("uQI56s9ZSf"),
			},
		},
		A6: &B{
			B1: 7604776742081145612,
			B2: "vDIkrP3hQs",
			B3: 0.8395196802725778,
			B4: thrift.StringPtr("5jbJcx4A4m"),
			B5: C{
				C1: 6281496988027306234,
				C2: thrift.Int64Ptr(int64(4015257045867961535)),
				C3: "u2utLX8xFM",
				C4: thrift.StringPtr("X5jwZIUL0A"),
			},
			B6: &C{
				C1: 876092331258217238,
				C2: thrift.Int64Ptr(int64(504834319128121057)),
				C3: "jObWv8JnpJ",
				C4: thrift.StringPtr("Zg3Gs9uEGI"),
			},
		},
	}
	_ = a
}

```

## Installation

```
cd fillstruct-rand
go install .
which fillstruct-rand
```

## Usage

```
% fillstruct [-modified] -file=<filename, required>  -line=<line number, required> -w=<filename to be written, optional> 
```

Flags:

	-file:     required. filename
	-line:     required. line number of the struct literal
	-w:        optional. when set this, the generated code will write to the file which names w

