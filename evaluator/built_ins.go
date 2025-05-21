package evaluator

import (
	"fmt"
	"interpreter/object"
	"slices"
)

var built_ins map[string]*object.BuiltIn

func init() {
	built_ins = map[string]*object.BuiltIn{
		"len": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `len`: got %d", len(args))
				}
				switch arg := args[0].(type) {
				case *object.String:
					return &object.Integer{Value: int64(len(arg.Value))}
				case *object.Array:
					return &object.Integer{Value: int64(len(arg.Elements))}
				default:
					return newError("argument to `len` not supported, got %s", args[0].Type())
				}
			},
		},
		"push": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments to `push`: got %d", len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `push` must be an array, but got %s", args[0].Type())
				}
				array := args[0].(*object.Array)
				length := len(array.Elements)

				newElements := make([]object.Object, length+1)
				copy(newElements, array.Elements)
				newElements[length] = args[1]

				return &object.Array{Elements: newElements}
			},
		},
		"pop": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `pop`: got %d", len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `pop` must be an array, but got %s", args[0].Type())
				}
				array := args[0].(*object.Array)
				lastElement := array.Elements[len(array.Elements)-1]
				array.Elements = array.Elements[0 : len(array.Elements)-1]

				return lastElement
			},
		},
		"concat": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments to `pop`: got %d", len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument 1 to `concat` must be an array, but got %s", args[0].Type())
				}
				if args[1].Type() != object.ARRAY_OBJ {
					return newError("argument 2 to `concat` must be an array, but got %s", args[0].Type())
				}
				array1 := args[0].(*object.Array)
				array2 := args[1].(*object.Array)
				length := len(array1.Elements) + len(array2.Elements)

				newElements := make([]object.Object, length)
				copy(newElements, append(array1.Elements, array2.Elements...))

				return &object.Array{Elements: newElements}
			},
		},
		"insert": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 3 {
					return newError("wrong number of arguments to `insert`: got %d", len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument 1 to `insert` must be an array, but got %s", args[0].Type())
				}
				if args[1].Type() != object.INTEGER_OBJ {
					return newError("argument 2 to `concat` must be an integer, but got %s", args[0].Type())
				}
				array := args[0].(*object.Array)
				index := args[1].(*object.Integer)
				value := args[2]

				return &object.Array{Elements: slices.Insert(array.Elements, int(index.Value), value)}

			},
		},
		"reverse": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `insert`: got %d", len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `reverse` must be an array, but got %s", args[0].Type())
				}
				array := args[0].(*object.Array)
				slices.Reverse(array.Elements)

				return &object.Array{Elements: array.Elements}

			},
		},
		"sort": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `insert`: got %d", len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `sort` must be an array, but got %s", args[0].Type())
				}
				array := args[0].(*object.Array)
				for _, obj := range array.Elements {
					if obj.Type() != object.INTEGER_OBJ {
						return newError("`sort` currently only supports arrays of integers")
					}
				}
				slices.SortFunc(array.Elements, func(a, b object.Object) int {
					return int(a.(*object.Integer).Value) - int(b.(*object.Integer).Value)
				})

				return &object.Array{Elements: array.Elements}

			},
		},
		"set": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments to `insert`: got %d", len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `set` must be an array, but got %s", args[0].Type())
				}
				seen := make(map[string]object.Object)
				array := args[0].(*object.Array)
				for _, obj := range array.Elements {
					seen[obj.Inspect()] = obj
				}
				keys := make([]object.Object, 0, len(seen))
				for k := range seen {
					keys = append(keys, seen[k])
				}

				return &object.Array{Elements: keys}

			},
		},
		"transform": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments to `insert`: got %d", len(args))
				}
				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument 1 to `transform` must be an array, but got %s", args[0].Type())
				}
				if args[1].Type() != object.FUNCTION_OBJ && args[1].Type() != object.BUILT_IN_OBJ {
					return newError("argument 2 to `transform` must be a function, but got %s", args[0].Type())
				}
				array := args[0].(*object.Array)
				fn := args[1]

				newElements := make([]object.Object, 0, len(array.Elements))

				for _, obj := range array.Elements {
					val := applyFunction(fn, []object.Object{obj})
					newElements = append(newElements, val)
				}

				return &object.Array{Elements: newElements}
			},
		},
		"print": {
			Fn: func(args ...object.Object) object.Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}
				return NULL
			},
		},
	}
}
