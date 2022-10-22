package db

// DataType define the data structure type.
type DataType = int8

// File different data types, support String, List, Hash, Set, Sorted Set right now.
const (
	String DataType = iota
	List
	Hash
	Set
	ZSet
)
