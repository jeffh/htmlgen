package h

// Body is an element body that streams child content through b. A nil Body
// writes an empty element.
type Body = func(*B)
