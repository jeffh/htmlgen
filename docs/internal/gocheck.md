# gocheck Reference (v0.2.0)

Property-based testing library for Go. Generate random test inputs and automatically shrink failing cases to minimal reproductions.

```
go get github.com/jeffh/gocheck
```

## Quick Start

```go
import (
    "testing"
    "github.com/jeffh/gocheck/check"
    "github.com/jeffh/gocheck/gen"
)

func TestAddCommutative(t *testing.T) {
    check.Spec(t, check.Config{}, gen.ForAll(
        gen.Tuple2(gen.Int, gen.Int),
        func(t gen.T2[int, int]) check.B {
            return t.A + t.B == t.B + t.A
        },
    ))
}
```

## Core API

### Running Tests

```go
// Spec runs a property test and reports to testing.T
check.Spec(t *testing.T, cfg Config, property Gen[Property[X]]) Property[X]

// Check runs a property and returns the result directly
check.Check(cfg Config, property Gen[Property[X]]) Property[X]
```

### Creating Properties

```go
// ForAll creates a property from a generator and test function
gen.ForAll(g Gen[X], testFn func(X) check.B) Gen[Property[X]]

// check.Prop is equivalent to gen.ForAll
check.Prop(g Gen[X], testFn func(X) check.B) Gen[Property[X]]

// PropT creates a property with testing.T subtest integration (v0.2.0)
// Each test iteration appears as a subtest in go test output
check.PropT(t *testing.T, g Gen[X], testFn func(X) check.B) Gen[Property[X]]
```

### Configuration

```go
check.Config{
    NumTests:       100,    // Number of test cases (default: 100)
    Seed:           0,      // Random seed (0 = random)
    MaxSizeHint:    0,      // Max size hint for sized generators (0 = no limit)
    NoShrink:       false,  // Disable shrinking
    PrintSeed:      false,  // Print seed before running
    ExpectsFailure: false,  // Expect the property to fail
    SkipLimit:      10,     // Max skips before failing
}
```

### Test Results (check.B)

Test functions return `check.B` which can be:
- `bool` - true = pass, false = fail
- `error` - nil = pass, non-nil = fail
- `check.Boolable` - custom type implementing `Bool() bool`

```go
// Helpers for test results
check.Fail(msg func() string) Boolable  // Create a failing result with message
check.Not(b B) B                         // Negate a result
check.Label(label string, res B) B       // Add label for distribution tracking
check.Postfix(b B, msg string) B         // Append message to result

// Human-readable representations (v0.2.0)
check.HumanRepresentation(r any) string
check.DerivedHumanRepresentation(format string, values ...any) string
```

## Generators

### Primitive Types

```go
// Booleans
gen.Bool                     // true/false, shrinks toward false
gen.BiasedBool(f, n uint64)  // false if random < f (out of n)

// Bytes
gen.Byte                     // 0-255, shrinks toward zero
gen.Bytes                    // Random byte slice
gen.BytesN(n)                // Exactly n bytes
gen.BytesRange(min, max)     // Between min and max bytes

// Integers (full range, shrink toward zero)
gen.Int, gen.Int8, gen.Int16, gen.Int32, gen.Int64
gen.Uint, gen.Uint8, gen.Uint16, gen.Uint32, gen.Uint64
gen.PosInt, gen.PosInt8, ...  // Positive only (includes zero)

// Non-zero integers (shrink toward 1)
gen.NonZeroInt, gen.NonZeroInt8, gen.NonZeroInt16, gen.NonZeroInt32, gen.NonZeroInt64

// Integers (sized - smaller values early in test run)
gen.SmallInt, gen.SmallInt8, gen.SmallInt16, gen.SmallInt32, gen.SmallInt64
gen.SmallUint, gen.SmallUint8, gen.SmallUint16, gen.SmallUint32, gen.SmallUint64
gen.SmallPosInt, gen.SmallPosInt8, ...

// Integer ranges
gen.SignedBetween[N](min, max N) Gen[N]        // Shrinks toward min
gen.SignedAround[N](min, mid, max N) Gen[N]    // Shrinks toward mid
gen.UnsignedBetween[N](min, max N) Gen[N]
gen.UnsignedAround[N](min, mid, max N) Gen[N]

// Floats (includes Inf, NaN rarely)
gen.Float32, gen.Float64
gen.Float32R(min, mid, max), gen.Float64R(min, mid, max)
gen.WithSpecialValues(floatGen)  // More likely to generate NaN, Inf
gen.NonZeroFloat32, gen.NonZeroFloat64

// Special numeric ranges
gen.Percentage      // 0.0 to 1.0
gen.Probability     // Alias for Percentage
gen.Degrees         // 0.0 to 360.0
gen.Radians         // 0.0 to 2π

// Complex
gen.Complex64, gen.Complex128

// Big integers
gen.BigInt, gen.PosBigInt

// Bitsets (v0.2.0)
gen.Bitset(flag1, flag2, ...)  // Combinations of flags via bitwise OR
```

### Strings and Runes

```go
// Runes
gen.Rune              // Any valid rune
gen.RuneASCII         // ASCII only (was RuneAscii in v0.1)
gen.RuneDigit         // '0'-'9'
gen.RuneAlpha         // a-z, A-Z
gen.RuneAlphaNumeric  // a-z, A-Z, 0-9
gen.RuneURLSafe       // Alphanumeric + -._~ (was RuneUrlSafe in v0.1)
gen.RunePrintable     // Printable runes
gen.RuneRange(min, max rune)
gen.RuneOf(set string)    // Runes from string

// Strings (shrink toward empty)
gen.StringAny           // Any runes (including non-printable)
gen.StringASCII         // ASCII only (was StringAscii in v0.1)
gen.StringNumbers       // Digits only
gen.StringAlpha         // Letters only
gen.StringAlphaNumeric  // Letters and digits
gen.StringURLSafe       // URL-safe characters (was StringUrlSafe in v0.1)
gen.StringPrintable     // Printable characters
gen.StringOf(runeGen)   // String from custom rune generator
gen.HexString           // Lowercase hex characters

// Identifiers and formats
gen.Identifier          // Valid Go identifiers, shrinks toward "a"
gen.UUID                // UUID v4 strings
gen.SemVer              // Semantic versions (1.2.3, 1.0.0-alpha)
gen.SemVerRange         // Version ranges (^1.2.3, >=1.0.0, *)

// String with options
gen.String(StringOptions{
    Slice: SliceOptions[rune]{
        Item:    gen.RuneAlpha,
        MinSize: 1,
        MaxSize: 50,
    },
})
```

### Time

```go
gen.Time                          // Shrinks toward Unix epoch
gen.TimeAfter(t time.Time)        // After t, shrinks toward t
gen.TimeBefore(t time.Time)       // Before t, shrinks toward t
gen.TimeBetween(t1, t2 time.Time) // Between t1 and t2, shrinks toward t1
gen.SmallDuration, gen.Duration
gen.PosSmallDuration, gen.PosDuration
gen.DayOfWeek     // time.Weekday
gen.MonthOfYear   // time.Month
gen.TimeZone      // IANA timezone identifiers
```

### Network

```go
// IP addresses
gen.IPv4              // net.IP, shrinks toward 0.0.0.0
gen.IPv6              // net.IP, shrinks toward ::
gen.IP                // Either IPv4 or IPv6
gen.HardwareAddr      // MAC address

// CIDR notation
gen.CIDR              // IPv4 or IPv6 CIDR
gen.IPv4CIDR          // e.g., "192.168.1.0/24"
gen.IPv6CIDR          // e.g., "2001:db8::/32"
gen.Netmask           // IPv4 netmask

// Ports
gen.Port              // 1-65535
gen.PortPrivileged    // 1-1023
gen.PortUnprivileged  // 1024-65535

// Addresses
gen.TCPAddr           // *net.TCPAddr
gen.UDPAddr           // *net.UDPAddr
gen.UnixAddr          // *net.UnixAddr

// DNS
gen.DNSRecordGen(recordType)  // A, AAAA, MX, TXT, CNAME, NS, PTR, SOA, SRV, CAA
```

### Web and URLs

```go
gen.Email             // Valid email addresses
gen.Hostname          // Valid hostnames
gen.URL               // url.URL with various components
gen.BasicURL          // Simpler URLs (http/https)

// TLDs
gen.TLDs              // Any TLD
gen.OriginalTLDs      // com, org, net, edu, gov, mil, int, arpa
gen.CountryTLDs       // Country-code TLDs

// URL with options
gen.URLOf(&URLGenOptions{
    Scheme: gen.Enum("http", "https"),
    Domain: gen.StringAlphaNumeric,
    Tld:    gen.TLDs,
    Path:   gen.Path(),
    Query:  gen.Slice(gen.Tuple2(gen.StringURLSafe, gen.StringASCII)),
})

// Path components
gen.Path()            // URL paths
gen.PathComponent     // Single path segment (URL-safe string)
```

### Geography

```go
gen.Latitude          // -90.0 to 90.0, shrinks toward 0 (equator)
gen.Longitude         // -180.0 to 180.0, shrinks toward 0 (prime meridian)
gen.LatLngGen         // LatLng struct with both coordinates
gen.CountryCode       // ISO 3166-1 alpha-2 codes
```

### Colors

```go
gen.HexColor          // "#RRGGBB" format, shrinks toward #000000
gen.RGB               // RGBColor struct (R, G, B 0-255)
gen.RGBA              // color.RGBA from image/color
gen.ColorName         // CSS color names ("red", "blue", etc.)
```

### Finance

```go
gen.Currency          // ISO 4217 codes ("USD", "EUR", etc.)
gen.CreditCard        // Valid card numbers (pass Luhn checksum)
gen.BTC               // Bitcoin amounts (0 to 21M with 8 decimals)
gen.IBAN              // International Bank Account Numbers
```

### Phone Numbers

```go
gen.PhoneNumberE164           // International E.164 format (+12025551234)
gen.PhoneNumberUS             // US formats ((555) 123-4567, etc.)
gen.PhoneNumberInternational  // Various international formats
```

### Files

```go
gen.FileMode          // fs.FileMode with type and permission bits
gen.FileModeUnix      // Unix permission bits (0-0777)
gen.DirMode           // Directory permissions (0700-0777)
gen.FileErrors(path)  // Common filesystem errors
```

### Collections

```go
// Slices
gen.Slice(itemGen)  // 0-100 elements
gen.SliceOf(SliceOptions[X]{
    Item:          itemGen,
    MinSize:       0,
    MaxSize:       100,
    PartitionSize: 0,  // Hint for shrinking large slices
})

// Maps
gen.Map(keyGen, valueGen)
gen.MapOf(MapOptions[K, V]{
    Keys:          keyGen,
    Values:        valueGen,
    MinSize:       0,
    MaxSize:       100,
    PartitionSize: 0,
})

// Channels (v0.2.0)
gen.Chan(ctx, itemGen)                    // Buffered channel with items
gen.ChanOf(ctx, ChanOptions[T]{...})      // Full control over channel
gen.ChanReceiveOnly(ctx, itemGen)         // <-chan T
gen.ChanSized(ctx, itemGen, sizeGen)      // Variable buffer size
gen.EmptyChan[T](capGen)                  // Empty channel with capacity
```

### Tuples

```go
gen.Tuple2(genA, genB)                    // T2[A, B]
gen.Tuple3(genA, genB, genC)              // T3[A, B, C]
gen.Tuple4(genA, genB, genC, genD)        // T4[A, B, C, D]
gen.Tuple5(genA, genB, genC, genD, genE)  // T5[A, B, C, D, E]
gen.Tuple6(...)                           // T6[...]
gen.TupleAny(gens ...Gen[any])            // []any

// Tuple fields
t := gen.T2[int, string]{A: 1, B: "hello"}
t.A  // first value
t.B  // second value
```

## Combinators

### Value Selection

```go
// Constant value (no shrinking)
gen.Value(x)

// Choose from values (shrinks toward first)
gen.Enum(val1, val2, val3, ...)
gen.EnumSlice([]T{val1, val2, ...})

// Choose from generators
gen.OneOf(gen1, gen2, ...)      // Shrinks value, keeps same generator
gen.Union(gen1, gen2, ...)      // Shrinks toward first generator
gen.Frequency(pairs ...)        // Weighted selection, shrinks toward first
gen.Pick(pairs ...)             // Weighted selection, keeps same generator

// Weighted pairs for Frequency/Pick
gen.P(weight uint32, generator)

// Combine generators (v0.2.0)
gen1.Combine(gen2)              // Merge two generators
```

### Transformations

```go
// Transform generated values
gen.Then(gen, func(X) Y) Gen[Y]

// Filter values (may panic if too many rejects)
gen.Filter(gen, func(X) bool) Gen[X]

// Ensure non-zero integers
gen.NonZero(intGen)
gen.NonZeroFn(gen, transformFn)  // Custom non-zero transform

// Type casting
gen.Cast[From, To](gen)
gen.ToAny(gen)  // Gen[X] -> Gen[any]

// Make integers sized (smaller early in test run)
gen.Sized(intGen)
```

## Custom Generators

### Simple Function Generator

```go
type myGen struct{}

func (g myGen) Gen(ptr *MyType, r *check.Run) error {
    // Use r.Rand for random values
    ptr.Field = r.Rand.Int64R(0, 100)
    return nil
}
```

### Using gen.Func

```go
myGen := gen.Func[MyType](func(ptr *MyType, r *check.Run) error {
    ptr.Value = r.Rand.Int64R(0, 100)
    return nil
})
```

### Stateful Generator with Clone

```go
gen.Factory(func() gen.Func[[]int] {
    buf := make([]int, 10)  // State initialized per clone
    return func(ptr *[]int, r *check.Run) error {
        for i := range buf {
            buf[i] = int(r.Rand.Int64R(0, 100))
        }
        *ptr = buf
        return nil
    }
})
```

### Human-Readable Representation (v0.2.0)

Implement `Representable` for custom types to improve error messages:

```go
type MyType struct {
    Name string
    Age  int
}

func (m MyType) HumanRepresentation() string {
    return fmt.Sprintf("MyType{Name: %q, Age: %d}", m.Name, m.Age)
}
```

### Random API (check.Run.Rand)

```go
r.Rand.Bool()                       // Fair coin flip
r.Rand.BiasedBool(f, n)             // false if < f (out of n)
r.Rand.Uint32(), r.Rand.Uint64()    // Full range
r.Rand.Uint32N(n), r.Rand.Uint64N(n) // [0, n]
r.Rand.Uint64R(min, max)            // [min, max]
r.Rand.Uint64RM(min, mid, max)      // [min, max], shrinks toward mid
r.Rand.Int64R(min, max)             // [min, max]
r.Rand.Int64RM(min, mid, max)       // [min, max], shrinks toward mid
```

### Generator Errors

```go
check.ErrSkipGen   // Skip this generation (counts toward skip limit)
check.ErrRetryGen  // Retry with different seed (doesn't count)
check.ErrSameGen   // Produced same value as before
```

## Examples

### Testing JSON Round-Trip

```go
func TestJSONRoundTrip(t *testing.T) {
    check.Spec(t, check.Config{}, gen.ForAll(gen.StringAny, func(s string) check.B {
        data, err := json.Marshal(s)
        if err != nil {
            return err
        }
        var result string
        if err := json.Unmarshal(data, &result); err != nil {
            return err
        }
        return s == result
    }))
}
```

### Testing with Multiple Inputs

```go
func TestMapOperations(t *testing.T) {
    check.Spec(t, check.Config{}, gen.ForAll(
        gen.Tuple2(gen.StringAlpha, gen.Int),
        func(t gen.T2[string, int]) check.B {
            m := make(map[string]int)
            m[t.A] = t.B
            return m[t.A] == t.B
        },
    ))
}
```

### Using PropT for Subtests (v0.2.0)

```go
func TestWithSubtests(t *testing.T) {
    check.Spec(t, check.Config{}, check.PropT(t, gen.Int, func(x int) check.B {
        return x*2 == x+x
    }))
}
```

### Custom Generator Composition

```go
type Person struct {
    Name string
    Age  int
}

var personGen = gen.Then(
    gen.Tuple2(gen.StringAlpha, gen.UnsignedBetween[int](0, 120)),
    func(t gen.T2[string, int]) Person {
        return Person{Name: t.A, Age: t.B}
    },
)

func TestPersonValidation(t *testing.T) {
    check.Spec(t, check.Config{}, gen.ForAll(personGen, func(p Person) check.B {
        return len(p.Name) >= 0 && p.Age >= 0 && p.Age <= 120
    }))
}
```

### Labeling for Distribution

```go
func TestWithLabels(t *testing.T) {
    check.Spec(t, check.Config{}, gen.ForAll(gen.Int, func(n int) check.B {
        var label string
        switch {
        case n < 0:
            label = "negative"
        case n == 0:
            label = "zero"
        default:
            label = "positive"
        }
        return check.Label(label, true)
    }))
}
```

### Testing with Bitflags (v0.2.0)

```go
type Permission uint8

const (
    None    Permission = 0
    Read    Permission = 1 << iota
    Write
    Execute
)

func TestPermissions(t *testing.T) {
    permGen := gen.Bitset(None, Read, Write, Execute)
    check.Spec(t, check.Config{}, gen.ForAll(permGen, func(p Permission) check.B {
        // Test that permissions can be combined and checked
        hasRead := p&Read != 0
        hasWrite := p&Write != 0
        return hasRead || hasWrite || p == None || p&Execute != 0
    }))
}
```

## Command Line Flags

```
-gocheck-num-tests=N   Number of test cases (default: 100)
-gocheck-seed=N        Random seed for reproducibility
```

## Tips

1. **Shrinking**: Values shrink toward "simpler" values (zero, empty, first element). Design generators so simpler inputs are meaningful test cases.

2. **Reproducibility**: Save failing seeds. Add them as explicit test cases:
   ```go
   check.Spec(t, check.Config{Seed: 12345}, property)
   ```

3. **Performance**: Use `gen.SmallInt` etc. for faster tests with smaller values.

4. **Debugging**: Set `PrintSeed: true` in Config to see seeds for all runs.

5. **Skip vs Retry**: Use `ErrSkipGen` when a generated value is invalid but counts toward the limit. Use `ErrRetryGen` for transient issues.

6. **Subtests (v0.2.0)**: Use `PropT` instead of `Prop` to get individual test iterations as subtests in `go test` output.

7. **Human-Readable Output (v0.2.0)**: Implement `Representable` interface on custom types for better failure messages.

## Migration from v0.1.0

Breaking changes in v0.2.0:
- `RuneAscii` renamed to `RuneASCII`
- `StringAscii` renamed to `StringASCII`
- `RuneUrlSafe` renamed to `RuneURLSafe`
- `StringUrlSafe` renamed to `StringURLSafe`
- `IP()` function changed to `IP` variable
- `StringPrintable()` function changed to `StringPrintable` variable
- `RunePrintable()` function changed to `RunePrintable` variable

New additions:
- `PropT` for subtest integration
- Many new domain-specific generators (colors, finance, phone numbers, etc.)
- Channel generators
- `Bitset` for flag combinations
- `Combine` method on generators
- `HumanRepresentation` for better error messages
