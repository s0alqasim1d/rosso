# ProtoBuf

https://stackoverflow.com/questions/41348512/protobuf-unmarshal-unknown

## Embedded message wire type

Currently Protocol Buffers is an inherently flawed protocol, in that it is not
self describing. This is evidenced by the fact that "type 2" wire type [1], can
be any of these values:

1. string
2. byte slice
3. embedded message

Its pretty straight forward to check for string or byte slice [2]. Also, not
all input is a valid message. That gives you this result:

valid string | valid message | result
-------------|---------------|-----------
no           | no            | byte slice
no           | yes           | message or byte slice
yes          | no            | string
yes          | yes           | message or string

Ideally the protocol would have an extra wire type, for example "type 6", to
designate embedded messages. But without that, its not possible in every case
to know the structure of the input, unless you have a schema. Other protocols
such as JSON are self describing, nothing is stopping ProtoBuf from being self
describing as well. Granted, you lose the field names, but even just having the
field numbers and values, with a defined structure as read from the input,
would be quite useful. Quoting from the spec:

> in other words, the last three bits of the number store the wire type

So that gives you this currently:

~~~
000 Varint
001 64-bit
010 Length-delimited
011 Start group
100 End group
101 32-bit
110
111
~~~

So you have at least two extra types that could be added. Even a single extra
type for embedded messages would essentially make ProtoBuf self describing.

1. https://developers.google.com/protocol-buffers/docs/encoding
2. https://github.com/89z/format/blob/v1.34.0/format.go#L17-L35

## Why not use generic add?

If you use generic add with method interface, it would allow recurive slices. We
could solve this by also adding type interface, but its probably simpler to just
not use generics for this.

## Why not implement a decoder?

I tried, but Unmarshal is way faster:

~~~
Benchmark_Decode-12       367295 ns/op         1481091 B/op       2588 allocs/op
Benchmark_Unmarshal-12     86033 ns/op           80245 B/op       1104 allocs/op
~~~

