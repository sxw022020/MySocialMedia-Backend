## Syntax:
1. `map[string]interface{}`
    - represents a map with string keys and values of any type.
    - `map` indicates that we're defining a `map` data structure.
    - `string` specifies the type of the `keys` in the map, in this case, they must be strings.
    - `interface{}` is used as the `value` type. In Go, interface{} is the empty interface, which represents any type. It can hold values of any underlying type.

## Basics:
1. The name of a `method` to be imported in one package, should have Capital letter in the beginning

## Configuration Thinking:
1. Goland is larger than vscode, we need a larger VM with more memory to install Goland in VM
2. 1 key-pair could be used for multiple VMs, no need to create different key-pairs for different VMs
3. To run a go program and use elasticsearch, we need to download them in the VM
4. Apply customized firewall to your VM
