/*
Package inspect provides means to do type serialization, deserialization and documentation.

A type that wants to be serializable should implement Inspectable interface
and a serializer or deserializer should implement InspectorImpl

Supported composite types are objects, arrays and maps (arrays of key-value pairs) with the key being of type string
*/
package inspect
